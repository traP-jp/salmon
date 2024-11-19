package model

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type CommandName string

const (
	JudgeVote CommandName = "judge-vote"
)

type ScheduledTask struct {
	Id          string
	Command     CommandName
	Arg         string
	ScheduledAt time.Time
	CreatedAt   time.Time
	ExecutedAt  sql.NullTime
}

func (c Client) CreateScheduledTask(command CommandName, arg string, scheduledAt time.Time) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	task := ScheduledTask{
		Id:          id.String(),
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

func (c Client) GetActiveScheduledTasks() ([]ScheduledTask, error) {
	var tasks []ScheduledTask
	err := c.db.Where("scheduled_at <= ? AND executed_at IS NULL", time.Now()).Find(&tasks).Error
	return tasks, err
}

func (c Client) UpdateScheduledTask(task ScheduledTask) error {
	return c.db.Save(&task).Error
}
