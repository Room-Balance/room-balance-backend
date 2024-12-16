package models

import (
	"database/sql/driver"
	"fmt"
)

// TaskType Enum
type TaskType string

const (
	TaskTypeOneTime    TaskType = "one-time"
	TaskTypeRepetitive TaskType = "repetitive"
)

func (t TaskType) Value() (driver.Value, error) {
	return string(t), nil
}

func (t *TaskType) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid type for TaskType: %T", value)
	}
	*t = TaskType(str)
	return nil
}

// TaskStatus Enum
type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in-progress"
	TaskStatusDone       TaskStatus = "done"
)

func (s TaskStatus) Value() (driver.Value, error) {
	return string(s), nil
}

func (s *TaskStatus) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid type for TaskStatus: %T", value)
	}
	*s = TaskStatus(str)
	return nil
}

// TaskFrequency Enum
type TaskFrequency string

const (
	TaskFrequencyDaily   TaskFrequency = "daily"
	TaskFrequencyWeekly  TaskFrequency = "weekly"
	TaskFrequencyMonthly TaskFrequency = "monthly"
)

func (f TaskFrequency) Value() (driver.Value, error) {
	return string(f), nil
}

func (f *TaskFrequency) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid type for TaskFrequency: %T", value)
	}
	*f = TaskFrequency(str)
	return nil
}
