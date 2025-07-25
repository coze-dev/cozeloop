// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/coze-dev/cozeloop/backend/infra/db"
	"github.com/coze-dev/cozeloop/backend/infra/idgen"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/domain/user/entity"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/domain/user/repo"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/infra/repo/mysql"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/infra/repo/mysql/convertor"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/infra/repo/mysql/gorm_gen/model"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/infra/repo/mysql/gorm_gen/query"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/pkg/errno"
	"github.com/coze-dev/cozeloop/backend/pkg/errorx"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/conv"
	"github.com/coze-dev/cozeloop/backend/pkg/lang/slices"
)

type UserRepoImpl struct {
	db             db.Provider
	idgen          idgen.IIDGenerator
	userDao        mysql.IUserDAO
	spaceDao       mysql.ISpaceDAO
	spaceMemberDao mysql.ISpaceUserDAO
}

func NewUserRepo(
	db db.Provider,
	idgen idgen.IIDGenerator,
	userDao mysql.IUserDAO,
	spaceDao mysql.ISpaceDAO,
	spaceMemberDao mysql.ISpaceUserDAO,
) repo.IUserRepo {
	return &UserRepoImpl{
		db:             db,
		idgen:          idgen,
		userDao:        userDao,
		spaceDao:       spaceDao,
		spaceMemberDao: spaceMemberDao,
	}
}

func (u UserRepoImpl) CreateUser(ctx context.Context, user *entity.User) (userID int64, err error) {
	if user == nil || user.Email == "" || user.HashPassword == "" {
		return 0, errorx.New("UserRepoImpl.CreateUser invalid param")
	}

	if user.UserID <= 0 {
		userID, err = u.idgen.GenID(ctx)
		if err != nil {
			return 0, errorx.Wrapf(err, "UserRepoImpl.CreateUser gen id error")
		}
		user.UserID = userID
	}

	// 初始化唯一用户名
	user.UniqueName = "user" + conv.ToString(user.UserID)
	if user.NickName == "" {
		user.NickName = user.UniqueName
	}

	err = u.db.Transaction(ctx, func(tx *gorm.DB) error {
		opt := db.WithTransaction(tx)
		// 创建用户
		userPO := convertor.UserDO2PO(user)
		err = u.userDao.Create(ctx, userPO, opt)
		if err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return errorx.WrapByCode(err, errno.CommonResourceDuplicatedCode, errorx.WithExtraMsg("userDao.CreateUser duplicate error"))
			}

			return errorx.WrapByCode(err, errno.CommonMySqlErrorCode, errorx.WithExtraMsg("userDao.CreateUser error"))
		}

		spaceID := user.SpaceID
		if spaceID <= 0 {
			// 创建空间
			spaceID, err = u.idgen.GenID(ctx)
			if err != nil {
				return errorx.Wrapf(err, "UserRepoImpl.CreateUser gen id error")
			}
			err = u.spaceDao.Create(ctx, &model.Space{
				ID:          spaceID,
				OwnerID:     user.UserID,
				Name:        "个人空间",
				Description: user.NickName + "的个人空间",
				SpaceType:   int32(entity.SpaceTypePersonal),
				IconURI:     "",
			}, opt)
			if err != nil {
				return err
			}
		}

		// 加入空间
		spaceMemberPO := &model.SpaceUser{
			SpaceID:  spaceID,
			UserID:   user.UserID,
			RoleType: int32(entity.SpaceUserTypeOwner),
		}
		err = u.spaceMemberDao.Create(ctx, spaceMemberPO, opt)
		if err != nil {
			return errorx.WrapByCode(err, errno.CommonMySqlErrorCode, errorx.WithExtraMsg("userDao.CreateSpaceUser error"))
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (u UserRepoImpl) GetUserByID(ctx context.Context, userID int64) (*entity.User, error) {
	var err error

	if userID <= 0 {
		return nil, errorx.New("UserRepoImpl.GetUserByID invalid param")
	}

	userPO, err := u.userDao.GetByID(ctx, userID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.CommonMySqlErrorCode, errorx.WithExtraMsg("userDao.GetUserByID error"))
	}
	user := convertor.UserPO2DO(userPO)
	return user, nil
}

func (u UserRepoImpl) MGetUserByIDs(ctx context.Context, userIDs []int64) ([]*entity.User, error) {
	var err error
	if len(userIDs) == 0 {
		return nil, errorx.New("UserRepoImpl.MGetUserByIDs invalid param")
	}
	userPOList, err := u.userDao.MGetByIDs(ctx, userIDs)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.CommonMySqlErrorCode, errorx.WithExtraMsg("userDao.MGetUserByIDs error"))
	}
	userList := make([]*entity.User, 0, len(userPOList))
	for _, userPO := range userPOList {
		user := convertor.UserPO2DO(userPO)
		userList = append(userList, user)
	}
	return userList, nil
}

func (u UserRepoImpl) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var err error

	if email == "" {
		return nil, errorx.New("UserRepoImpl.GetUserByEmail invalid param")
	}
	userPO, err := u.userDao.FindByEmail(ctx, email)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.CommonMySqlErrorCode, errorx.WithExtraMsg("userDao.GetUserByEmail error"))
	}
	user := convertor.UserPO2DO(userPO)
	return user, nil
}

func (u UserRepoImpl) UpdateSessionKey(ctx context.Context, userID int64, sessionKey string) error {
	var err error
	if userID <= 0 || sessionKey == "" {
		return errorx.New("UserRepoImpl.UpdateSessionKey invalid param")
	}

	q := query.Use(u.db.NewSession(ctx))
	err = u.userDao.Update(ctx, userID, map[string]interface{}{
		q.User.SessionKey.ColumnName().String(): sessionKey,
	})
	if err != nil {
		return errorx.WrapByCode(err, errno.CommonMySqlErrorCode, errorx.WithExtraMsg("userDao.UpdateUserAttr error"))
	}
	return nil
}

func (u UserRepoImpl) ClearSessionKey(ctx context.Context, userID int64) error {
	q := query.Use(u.db.NewSession(ctx))
	err := u.userDao.Update(ctx, userID, map[string]interface{}{
		q.User.SessionKey.ColumnName().String(): "",
	})
	if err != nil {
		return errorx.WrapByCode(err, errno.CommonMySqlErrorCode, errorx.WithExtraMsg("userDao.UpdateUserAttr error"))
	}
	return nil
}

func (u UserRepoImpl) UpdatePassword(ctx context.Context, userID int64, password string) error {
	q := query.Use(u.db.NewSession(ctx))
	err := u.userDao.Update(ctx, userID, map[string]interface{}{
		q.User.Password.ColumnName().String(): password,
	})
	if err != nil {
		return errorx.WrapByCode(err, errno.CommonMySqlErrorCode, errorx.WithExtraMsg("userDao.UpdatePassword error"))
	}
	return nil
}

func (u UserRepoImpl) CheckUniqueNameExist(ctx context.Context, uniqueName string) (bool, error) {
	_, err := u.userDao.FindByUniqueName(ctx, uniqueName)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, errorx.WrapByCode(err, errno.CommonMySqlErrorCode, errorx.WithExtraMsg("userDao.CheckUniqueNameExist error"))
	}
	return true, nil
}

func (u UserRepoImpl) CheckEmailExist(ctx context.Context, email string) (bool, error) {
	_, err := u.userDao.FindByEmail(ctx, email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, errorx.WrapByCode(err, errno.CommonMySqlErrorCode, errorx.WithExtraMsg("userDao.CheckEmailExist error"))
	}
	return true, nil
}

func (u UserRepoImpl) UpdateProfile(ctx context.Context, userID int64, param *repo.UpdateProfileParam) (*entity.User, error) {
	if param == nil || userID <= 0 {
		return nil, errorx.New("UserRepoImpl.UpdateProfile invalid param")
	}

	q := query.Use(u.db.NewSession(ctx))
	updates := map[string]interface{}{}

	if param.UniqueName != nil {
		updates[q.User.UniqueName.ColumnName().String()] = *param.UniqueName
	}

	if param.NickName != nil {
		updates[q.User.Name.ColumnName().String()] = *param.NickName
	}

	if param.Description != nil {
		updates[q.User.Description.ColumnName().String()] = *param.Description
	}

	if param.IconURI != nil {
		updates[q.User.IconURI.ColumnName().String()] = *param.IconURI
	}

	// 如果没有更新任何字段，则返回
	if len(updates) == 0 {
		return nil, errorx.New("noting need update")
	}

	err := u.userDao.Update(ctx, userID, updates)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errorx.WrapByCode(err, errno.CommonResourceDuplicatedCode, errorx.WithExtraMsg("userDao.UpdateUserAttr duplicate error"))
		}
		return nil, errorx.WrapByCode(err, errno.CommonMySqlErrorCode, errorx.WithExtraMsg("userDao.UpdateUserAttr error"))
	}

	// 更新后，通过主库获取最新用户信息
	userPO, err := u.userDao.GetByID(ctx, userID, db.WithMaster())
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.CommonMySqlErrorCode, errorx.WithExtraMsg("userDao.GetUserByID error"))
	}

	return convertor.UserPO2DO(userPO), nil
}

func (u UserRepoImpl) UpdateAvatar(ctx context.Context, userID int64, iconURI string) error {
	if userID <= 0 || iconURI == "" {
		return errorx.New("UserRepoImpl.UpdateAvatar invalid param")
	}

	q := query.Use(u.db.NewSession(ctx))
	err := u.userDao.Update(ctx, userID, map[string]interface{}{
		q.User.IconURI.ColumnName().String(): iconURI,
	})
	if err != nil {
		return errorx.WrapByCode(err, errno.CommonMySqlErrorCode, errorx.WithExtraMsg("userDao.UpdateUserAttr error"))
	}
	return nil
}

func (u UserRepoImpl) ListUserSpace(ctx context.Context, userID int64, pageSize, pageNum int32) ([]*entity.Space, int32, error) {
	userSpaceList, total, err := u.spaceMemberDao.List(ctx, userID, pageSize, pageNum)
	if err != nil {
		return nil, 0, errorx.WrapByCode(err, errno.CommonMySqlErrorCode, errorx.WithExtraMsg("spaceMemberDao.GetUserSpaceList error"))
	}
	spaceIDs := slices.Transform(userSpaceList, func(us *model.SpaceUser, idx int) int64 {
		return us.SpaceID
	})

	spaceList, err := u.spaceDao.MGetByIDs(ctx, spaceIDs)
	if err != nil {
		return nil, 0, errorx.WrapByCode(err, errno.CommonMySqlErrorCode, errorx.WithExtraMsg("spaceDao.MGetSpaceByIDs error"))
	}
	spaceDOList := make([]*entity.Space, 0, len(spaceList))
	for _, space := range spaceList {
		spaceDOList = append(spaceDOList, convertor.SpacePO2DO(space))
	}

	return spaceDOList, total, nil
}

func (u UserRepoImpl) CreateSpace(ctx context.Context, space *entity.Space) (spaceID int64, err error) {
	if space == nil {
		return 0, errorx.New("UserRepoImpl.CreateSpace invalid param: space nil")
	}

	spaceID, err = u.idgen.GenID(ctx)
	if err != nil {
		return 0, err
	}

	spacePO := convertor.SpaceDO2PO(space)
	spacePO.ID = spaceID
	spacePO.CreatedBy = space.OwnerID

	err = u.spaceDao.Create(ctx, spacePO)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return 0, errorx.WrapByCode(err, errno.CommonResourceDuplicatedCode)
		}
		return 0, errorx.WrapByCode(err, errno.CommonMySqlErrorCode, errorx.WithExtraMsg("spaceDao.CreateSpace error"))
	}
	spaceID = spacePO.ID

	return spaceID, nil
}

func (u UserRepoImpl) GetSpaceByID(ctx context.Context, spaceID int64) (*entity.Space, error) {
	var err error

	if spaceID <= 0 {
		return nil, errorx.New("UserRepoImpl.GetSpaceByID invalid param")
	}

	spacePO, err := u.spaceDao.GetByID(ctx, spaceID)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.CommonMySqlErrorCode, errorx.WithExtraMsg("spaceDao.GetSpaceByID error"))
	}
	spaceDO := convertor.SpacePO2DO(spacePO)
	return spaceDO, nil
}

func (u UserRepoImpl) MGetSpaceByIDs(ctx context.Context, spaceIDs []int64) (spaceDOList []*entity.Space, err error) {
	if len(spaceIDs) <= 0 {
		return nil, errorx.New("UserRepoImpl.MGetSpaceByIDs invalid param")
	}

	spacePOList, err := u.spaceDao.MGetByIDs(ctx, spaceIDs)
	if err != nil {
		return nil, errorx.WrapByCode(err, errno.CommonMySqlErrorCode, errorx.WithExtraMsg("spaceDao.MGetSpaceByIDs error"))
	}

	for _, spacePO := range spacePOList {
		spaceDOList = append(spaceDOList, convertor.SpacePO2DO(spacePO))
	}

	return spaceDOList, nil
}

func (u UserRepoImpl) CheckUserSpaceExist(ctx context.Context, userID, spaceID int64) (bool, error) {
	space, err := u.GetSpaceByID(ctx, spaceID)
	if err != nil {
		return false, err
	}

	// OpenSource Version only check owner's space
	if space.OwnerID == userID {
		return true, nil
	}
	return false, nil
}
