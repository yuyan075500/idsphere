package model

import (
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"time"
)

type ScheduledTask struct {
	gorm.Model
	Name                  string        `json:"name"`                                // 任务名称
	Type                  uint          `json:"type"`                                // 执行方法（1：调用URL，2：内置方法）
	CronExpr              string        `json:"cron_expr" gorm:"default:null"`       // 周期任务的表达式
	BuiltInMethod         string        `json:"built_in_method" gorm:"default:null"` // 内置方法
	Enabled               bool          `json:"enabled"`                             // 是否启用
	NotifyType            *int          `json:"notify_type"`                         // 通知方式
	Receiver              *string       `json:"receiver"`                            // 通知接收人（可以是：收件人、WebhookURL）
	LastRunAt             *time.Time    `json:"last_run_at"`                         // 上次运行时间
	LastRunResult         string        `json:"last_run_result"`                     // 上次执行结果
	NextRunAt             *time.Time    `json:"next_run_at"`                         // 下次运行时间
	EntryID               *cron.EntryID `json:"entry_id" gorm:"unique"`              // 任务运行ID
	ExecutionCount        uint          `json:"execution_count" gorm:"default:0"`    // 执行次数
	ScheduledTaskExecLogs []ScheduledTaskExecLog
}

type ScheduledTaskExecLog struct {
	ID              uint       `json:"id" json:"primaryKey;autoIncrement"`
	ScheduledTaskID uint       `json:"scheduled_task_id"`
	RunAt           *time.Time `json:"run_at"`
	FinishAt        *time.Time `json:"finish_at"`
	Result          string     `json:"result"`
}
