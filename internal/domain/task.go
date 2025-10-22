package domain

import (
	"errors"
	"time"
)

type TaskStatus string

const (
	StatusTodo  TaskStatus = "todo"
	StatusDoing TaskStatus = "doing"
	StatusDone  TaskStatus = "done"
)

type TaskPriority string

const (
	PriorityLow      TaskPriority = "low"
	PriorityMedium   TaskPriority = "medium"
	PriorityHigh     TaskPriority = "high"
	PriorityCritical TaskPriority = "critical"
)

type Task struct {
	ID          int64        `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Status      TaskStatus   `json:"status"`
	Priority    TaskPriority `json:"priority"`
	DueDate     *time.Time   `json:"due_date"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type CreateTaskRequest struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Priority    TaskPriority `json:"priority"`
	DueDate     *time.Time   `json:"due_date"`
}

type UpdateTaskRequest struct {
	Title       *string       `json:"title,omitempty"`
	Description *string       `json:"description,omitempty"`
	Status      *TaskStatus   `json:"status,omitempty"`
	Priority    *TaskPriority `json:"priority,omitempty"`
	DueDate     *time.Time    `json:"due_date,omitempty"`
}

func (t *Task) Validate() error {
	if t.Title == "" {
		return errors.New("title is required")
	}
	if len(t.Title) > 200 {
		return errors.New("title must be less than 200 characters")
	}
	if len(t.Description) > 1000 {
		return errors.New("description must be less than 1000 characters")
	}
	if !isValidStatus(t.Status) {
		return errors.New("invalid status")
	}
	if !isValidPriority(t.Priority) {
		return errors.New("invalid priority")
	}
	if t.DueDate != nil && t.DueDate.Before(time.Now().Truncate(24*time.Hour)) {
		return errors.New("due date cannot be in the past")
	}
	return nil
}

func (r *CreateTaskRequest) Validate() error {
	if r.Title == "" {
		return errors.New("title is required")
	}
	if len(r.Title) > 200 {
		return errors.New("title must be less than 200 characters")
	}
	if len(r.Description) > 1000 {
		return errors.New("description must be less than 1000 characters")
	}
	if !isValidPriority(r.Priority) {
		return errors.New("invalid priority")
	}
	if r.DueDate != nil && r.DueDate.Before(time.Now().Truncate(24*time.Hour)) {
		return errors.New("due date cannot be in the past")
	}
	return nil
}

func (r *UpdateTaskRequest) Validate() error {
	if r.Title != nil {
		if *r.Title == "" {
			return errors.New("title cannot be empty")
		}
		if len(*r.Title) > 200 {
			return errors.New("title must be less than 200 characters")
		}
	}
	if r.Description != nil && len(*r.Description) > 1000 {
		return errors.New("description must be less than 1000 characters")
	}
	if r.Status != nil && !isValidStatus(*r.Status) {
		return errors.New("invalid status")
	}
	if r.Priority != nil && !isValidPriority(*r.Priority) {
		return errors.New("invalid priority")
	}
	return nil
}

func isValidStatus(status TaskStatus) bool {
	switch status {
	case StatusTodo, StatusDoing, StatusDone:
		return true
	default:
		return false
	}
}

func isValidPriority(priority TaskPriority) bool {
	switch priority {
	case PriorityLow, PriorityMedium, PriorityHigh, PriorityCritical:
		return true
	default:
		return false
	}
}

func GetPriorityOrder(priority TaskPriority) int {
	switch priority {
	case PriorityCritical:
		return 4
	case PriorityHigh:
		return 3
	case PriorityMedium:
		return 2
	case PriorityLow:
		return 1
	default:
		return 0
	}
}
