package database

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type ScheduledTask struct {
	Id          string
	Command     string
	Arg         string
	ScheduledAt time.Time
	CreatedAt   time.Time
	ExecutedAt  sql.NullTime
}

func (c Client) CreateScheduledTask(command string, arg string, scheduledAt time.Time) error {
	task := ScheduledTask{
		Id:          uuid.NewString(),
		Command:     command,
		Arg:         arg,
		ScheduledAt: scheduledAt,
		CreatedAt:   time.Now(),
	}
	return c.db.Create(&task).Error
}

func (c Client) GetScheduledTask(id string) (ScheduledTask, error) {
	var task ScheduledTask
	err := c.db.Where("id = ?", id).First(&task).Error
	return task, err
}

func (c Client) GetScheduledTasks() ([]ScheduledTask, error) {
	var tasks []ScheduledTask
	err := c.db.Find(&tasks).Error
	return tasks, err
}

func (c Client) UpdateScheduledTask(task ScheduledTask) error {
	return c.db.Save(&task).Error
}
