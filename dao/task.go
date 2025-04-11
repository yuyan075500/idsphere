package dao

import (
	"gorm.io/gorm"
	"ops-api/global"
	"ops-api/model"
)

var Task task

type task struct{}

// TaskList 任务列表
type TaskList struct {
	Items []*model.ScheduledTask `json:"items"`
	Total int64                  `json:"total"`
}

// TaskLogList 任务执行日志列表
type TaskLogList struct {
	Items []*model.ScheduledTaskExecLog `json:"items"`
	Total int64                         `json:"total"`
}

// TaskUpdate 更新构体
type TaskUpdate struct {
	ID            uint    `json:"id" binding:"required"`
	Name          string  `json:"name"`
	Type          uint    `json:"type"`
	CronExpr      string  `json:"cron_expr"`
	Method        uint    `json:"method"`
	BuiltInMethod string  `json:"built_in_method"`
	Enabled       bool    `json:"enabled"`
	NotifyType    *int    `json:"notify_type"`
	Receiver      *string `json:"receiver"`
}

// AddTask 新增定时任务
func (t *task) AddTask(data *model.ScheduledTask) (task *model.ScheduledTask, err error) {
	if err := global.MySQLClient.Create(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// DeleteTask 删除定时任务
func (t *task) DeleteTask(tx *gorm.DB, id int) (err error) {
	return tx.Where("id = ?", id).Unscoped().Delete(&model.ScheduledTask{}).Error
}

// UpdateTask 修改定时任务
func (t *task) UpdateTask(task *model.ScheduledTask, data *TaskUpdate) (*model.ScheduledTask, error) {
	if err := global.MySQLClient.Model(task).Select("*").Where("id = ?", data.ID).Updates(data).Error; err != nil {
		return nil, err
	}
	return task, nil
}

// GetTaskList 获取定时任务列表
func (t *task) GetTaskList(name string, page, limit int) (data *TaskList, err error) {

	// 定义数据的起始位置
	startSet := (page - 1) * limit

	// 定义返回的内容
	var (
		items []*model.ScheduledTask
		total int64
	)

	// 获取菜单列表
	tx := global.MySQLClient.Model(&model.ScheduledTask{}).
		Where("name like ?", "%"+name+"%").
		Count(&total).
		Limit(limit).
		Offset(startSet).
		Find(&items)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &TaskList{
		Items: items,
		Total: total,
	}, nil
}

// GetTaskLogList 获取定时任务执行日志列表
func (t *task) GetTaskLogList(id uint, page, limit int) (data *TaskLogList, err error) {

	// 定义数据的起始位置
	startSet := (page - 1) * limit

	// 定义返回的内容
	var (
		items []*model.ScheduledTaskExecLog
		total int64
	)

	// 获取菜单列表
	tx := global.MySQLClient.Model(&model.ScheduledTaskExecLog{}).
		Where("scheduled_task_id = ?", id).
		Count(&total).
		Limit(limit).
		Offset(startSet).
		Order("id desc").
		Find(&items)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &TaskLogList{
		Items: items,
		Total: total,
	}, nil
}

// DeleteTaskLogList 删除定时任务执行日志
func (t *task) DeleteTaskLogList(tx *gorm.DB, id int) error {
	return tx.Where("scheduled_task_id = ?", id).Delete(&model.ScheduledTaskExecLog{}).Error
}
