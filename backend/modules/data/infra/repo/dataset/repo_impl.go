// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package dataset

import (
	"github.com/coze-dev/cozeloop/backend/infra/db"
	"github.com/coze-dev/cozeloop/backend/infra/idgen"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/dataset/repo"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/entity"
	"github.com/coze-dev/cozeloop/backend/modules/data/infra/repo/dataset/item_dao"
	"github.com/coze-dev/cozeloop/backend/modules/data/infra/repo/dataset/mysql"
	"github.com/coze-dev/cozeloop/backend/modules/data/infra/repo/dataset/redis"
)

type DatasetRepo struct {
	idGen           idgen.IIDGenerator
	txDB            db.Provider
	datasetDAO      mysql.IDatasetDAO
	schemaDAO       mysql.ISchemaDAO
	datasetRedisDAO redis.DatasetDAO
	versionDao      mysql.IVersionDAO
	versionRedisDao redis.VersionDAO
	optDAO          redis.OperationDAO
	itemDAO         mysql.IItemDAO
	itemSnapshotDAO mysql.IItemSnapshotDAO
	ioJobDAO        mysql.IIOJobDAO
	itemProviderDAO map[entity.Provider]item_dao.ItemDAO
}

var _ repo.IDatasetAPI = (*DatasetRepo)(nil)

func NewDatasetRepo(idgen idgen.IIDGenerator, db db.Provider, datasetDAO mysql.IDatasetDAO, schemaDAO mysql.ISchemaDAO, datasetRedisDAO redis.DatasetDAO, versionDAO mysql.IVersionDAO, versionRedisDAO redis.VersionDAO, optDAO redis.OperationDAO, itemDAO mysql.IItemDAO, itemSnapshotDAO mysql.IItemSnapshotDAO, ioJobDAO mysql.IIOJobDAO, itemProviderDAO map[entity.Provider]item_dao.ItemDAO) repo.IDatasetAPI {
	return &DatasetRepo{
		idGen:           idgen,
		txDB:            db,
		datasetDAO:      datasetDAO,
		schemaDAO:       schemaDAO,
		datasetRedisDAO: datasetRedisDAO,
		versionDao:      versionDAO,
		versionRedisDao: versionRedisDAO,
		optDAO:          optDAO,
		itemDAO:         itemDAO,
		itemSnapshotDAO: itemSnapshotDAO,
		ioJobDAO:        ioJobDAO,
		itemProviderDAO: itemProviderDAO,
	}
}
