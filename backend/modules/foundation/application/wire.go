// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: Apache-2.0

//go:build wireinject
// +build wireinject

package application

import (
	"github.com/google/wire"

	"github.com/coze-dev/cozeloop/backend/infra/db"
	"github.com/coze-dev/cozeloop/backend/infra/fileserver"
	"github.com/coze-dev/cozeloop/backend/infra/idgen"
	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/foundation/auth"
	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/foundation/auth/authservice"
	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/foundation/authn"
	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/foundation/file"
	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/foundation/openapi"
	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/foundation/space"
	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/foundation/user"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/domain/user/service"
	auth2 "github.com/coze-dev/cozeloop/backend/modules/foundation/infra/auth"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/infra/repo"
	"github.com/coze-dev/cozeloop/backend/modules/foundation/infra/repo/mysql"
)

var (
	userDomainSet = wire.NewSet(
		service.NewUserService,
		repo.NewUserRepo,
		mysql.NewUserDAOImpl,
		mysql.NewSpaceDAOImpl,
		mysql.NewSpaceUserDAOImpl,
	)

	userSet = wire.NewSet(
		NewUserApplication,
		userDomainSet,
	)

	spaceSet = wire.NewSet(
		NewSpaceApplication,
		userDomainSet,
	)

	authSet = wire.NewSet(
		NewAuthApplication,
		userDomainSet,
	)

	authNSet = wire.NewSet(
		NewAuthNApplication,
		repo.NewAuthNRepo,
		mysql.NewAuthNDAOImpl,
	)

	fileSet = wire.NewSet(
		NewFileApplication,
		auth2.NewAuthProvider,
	)

	openAPISet = wire.NewSet(
		NewFoundationOpenAPIApplication,
		auth2.NewAuthProvider,
	)
)

func InitAuthApplication(idgen idgen.IIDGenerator,
	db db.Provider,
) (auth.AuthService, error) {
	wire.Build(authSet)
	return nil, nil
}

func InitAuthNApplication(
	idgen idgen.IIDGenerator,
	db db.Provider,
) (authn.AuthNService, error) {
	wire.Build(authNSet)
	return nil, nil
}

func InitSpaceApplication(
	idgen idgen.IIDGenerator,
	db db.Provider,
) (space.SpaceService, error) {
	wire.Build(spaceSet)
	return nil, nil
}

func InitUserApplication(
	idgen idgen.IIDGenerator,
	db db.Provider,
) (user.UserService, error) {
	wire.Build(userSet)
	return nil, nil
}

func InitFileApplication(
	objectStorage fileserver.BatchObjectStorage,
	authClient authservice.Client,
) (file.FileService, error) {
	wire.Build(fileSet)
	return nil, nil
}

func InitFoundationOpenAPIApplication(
	objectStorage fileserver.BatchObjectStorage,
	authClient authservice.Client,
) (openapi.FoundationOpenAPIService, error) {
	wire.Build(openAPISet)
	return nil, nil
}
