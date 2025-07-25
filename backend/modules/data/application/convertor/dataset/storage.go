// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package dataset

import (
	"github.com/coze-dev/cozeloop/backend/kitex_gen/coze/loop/data/domain/dataset"
	"github.com/coze-dev/cozeloop/backend/modules/data/domain/entity"
)

func StorageProviderDTO2DO(s dataset.StorageProvider) entity.Provider {
	switch s {
	case dataset.StorageProvider_TOS:
		return entity.ProviderTOS
	case dataset.StorageProvider_VETOS:
		return entity.ProviderVETOS
	case dataset.StorageProvider_HDFS:
		return entity.ProviderHDFS
	case dataset.StorageProvider_ImageX:
		return entity.ProviderImageX
	case dataset.StorageProvider_S3:
		return entity.ProviderS3
	case dataset.StorageProvider_LocalFS:
		return entity.ProviderLocalFS
	case dataset.StorageProvider_Abase:
		return entity.ProviderAbase
	default:
		return entity.ProviderUnknown
	}
}

func ProviderDO2DTO(sp entity.Provider) dataset.StorageProvider {
	switch sp {
	case entity.ProviderTOS:
		return dataset.StorageProvider_TOS
	case entity.ProviderVETOS:
		return dataset.StorageProvider_VETOS
	case entity.ProviderHDFS:
		return dataset.StorageProvider_HDFS
	case entity.ProviderImageX:
		return dataset.StorageProvider_ImageX
	case entity.ProviderS3:
		return dataset.StorageProvider_S3
	case entity.ProviderLocalFS:
		return dataset.StorageProvider_LocalFS
	case entity.ProviderAbase:
		return dataset.StorageProvider_Abase
	case entity.ProviderRDS:
		return dataset.StorageProvider_RDS
	default:
		return dataset.ObjectStorage_Provider_DEFAULT
	}
}
