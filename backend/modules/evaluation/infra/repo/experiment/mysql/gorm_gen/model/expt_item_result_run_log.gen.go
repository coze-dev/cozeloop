// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameExptItemResultRunLog = "expt_item_result_run_log"

// ExptItemResultRunLog expt_item_result_run_log
type ExptItemResultRunLog struct {
	ID          int64          `gorm:"column:id;type:bigint(20) unsigned;primaryKey;comment:id" json:"id"`                                                                                                                                          // id
	SpaceID     int64          `gorm:"column:space_id;type:bigint(20) unsigned;not null;uniqueIndex:uk_expt_run_item_turn,priority:1;index:idx_expt_item_turn,priority:1;index:idx_expt_run_result_state,priority:1;comment:空间 id" json:"space_id"` // 空间 id
	ExptID      int64          `gorm:"column:expt_id;type:bigint(20) unsigned;not null;uniqueIndex:uk_expt_run_item_turn,priority:2;index:idx_expt_item_turn,priority:2;index:idx_expt_run_result_state,priority:2;comment:实验 id" json:"expt_id"`   // 实验 id
	ExptRunID   int64          `gorm:"column:expt_run_id;type:bigint(20) unsigned;not null;uniqueIndex:uk_expt_run_item_turn,priority:3;index:idx_expt_run_result_state,priority:3;comment:实验运行 id" json:"expt_run_id"`                             // 实验运行 id
	ItemID      int64          `gorm:"column:item_id;type:bigint(20) unsigned;not null;uniqueIndex:uk_expt_run_item_turn,priority:4;index:idx_expt_item_turn,priority:3;comment:item_id" json:"item_id"`                                            // item_id
	Status      int32          `gorm:"column:status;type:int(11) unsigned;not null;comment:状态" json:"status"`                                                                                                                                       // 状态
	ErrMsg      *[]byte        `gorm:"column:err_msg;type:blob binary;comment:错误信息" json:"err_msg"`                                                                                                                                                 // 错误信息
	CreatedAt   time.Time      `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`                                                                                                          // 创建时间
	UpdatedAt   time.Time      `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`                                                                                                          // 更新时间
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp;comment:删除时间" json:"deleted_at"`                                                                                                                                             // 删除时间
	LogID       string         `gorm:"column:log_id;type:varchar(128);not null;comment:日志 id" json:"log_id"`                                                                                                                                        // 日志 id
	ResultState *int32         `gorm:"column:result_state;type:int(11);index:idx_expt_run_result_state,priority:4;comment:回写结果表状态" json:"result_state"`                                                                                             // 回写结果表状态
}

// TableName ExptItemResultRunLog's table name
func (*ExptItemResultRunLog) TableName() string {
	return TableNameExptItemResultRunLog
}
