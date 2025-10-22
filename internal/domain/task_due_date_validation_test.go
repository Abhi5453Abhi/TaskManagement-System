package domain

import (
	"testing"
	"time"
)

func TestTaskDueDateValidation_Buggy(t *testing.T) {
	tests := []struct {
		name    string
		task    Task
		wantErr bool
	}{
		{
			name: "valid task with future due date",
			task: Task{
				Title:       "Valid Task",
				Description: "Valid description",
				Status:      StatusTodo,
				Priority:    PriorityMedium,
				DueDate:     func() *time.Time { t := time.Now().AddDate(0, 0, 1); return &t }(),
			},
			wantErr: false,
		},
		{
			name: "task with past due date should fail validation",
			task: Task{
				Title:       "Task with Past Due Date",
				Description: "This should fail validation",
				Status:      StatusTodo,
				Priority:    PriorityMedium,
				DueDate:     func() *time.Time { t := time.Now().AddDate(0, 0, -1); return &t }(),
			},
			wantErr: true, // This should fail but won't due to the bug
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.task.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Task.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateTaskRequestDueDateValidation_Buggy(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateTaskRequest
		wantErr bool
	}{
		{
			name: "valid request with future due date",
			req: CreateTaskRequest{
				Title:       "Valid Task",
				Description: "Valid description",
				Priority:    PriorityMedium,
				DueDate:     func() *time.Time { t := time.Now().AddDate(0, 0, 1); return &t }(),
			},
			wantErr: false,
		},
		{
			name: "request with past due date should fail validation",
			req: CreateTaskRequest{
				Title:       "Task with Past Due Date",
				Description: "This should fail validation",
				Priority:    PriorityMedium,
				DueDate:     func() *time.Time { t := time.Now().AddDate(0, 0, -1); return &t }(),
			},
			wantErr: true, // This should fail but won't due to the bug
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTaskRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
