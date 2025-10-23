package domain

import (
	"testing"
	"time"
)

func TestTaskDueDateValidation_Fixed(t *testing.T) {
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
			wantErr: true,
		},
		{
			name: "task with due date exactly now should fail validation",
			task: Task{
				Title:       "Task with Current Due Date",
				Description: "This should fail validation",
				Status:      StatusTodo,
				Priority:    PriorityMedium,
				DueDate:     func() *time.Time { t := time.Now(); return &t }(),
			},
			wantErr: true,
		},
		{
			name: "task with due date in the past by hours should fail validation",
			task: Task{
				Title:       "Task with Past Due Date",
				Description: "This should fail validation",
				Status:      StatusTodo,
				Priority:    PriorityMedium,
				DueDate:     func() *time.Time { t := time.Now().Add(-2 * time.Hour); return &t }(),
			},
			wantErr: true,
		},
		{
			name: "task with no due date should pass validation",
			task: Task{
				Title:       "Task without Due Date",
				Description: "This should pass validation",
				Status:      StatusTodo,
				Priority:    PriorityMedium,
				DueDate:     nil,
			},
			wantErr: false,
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

func TestCreateTaskRequestDueDateValidation_Fixed(t *testing.T) {
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
			wantErr: true,
		},
		{
			name: "request with due date exactly now should fail validation",
			req: CreateTaskRequest{
				Title:       "Task with Current Due Date",
				Description: "This should fail validation",
				Priority:    PriorityMedium,
				DueDate:     func() *time.Time { t := time.Now(); return &t }(),
			},
			wantErr: true,
		},
		{
			name: "request with no due date should pass validation",
			req: CreateTaskRequest{
				Title:       "Task without Due Date",
				Description: "This should pass validation",
				Priority:    PriorityMedium,
				DueDate:     nil,
			},
			wantErr: false,
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

func TestUpdateTaskRequestDueDateValidation_Fixed(t *testing.T) {
	tests := []struct {
		name    string
		req     UpdateTaskRequest
		wantErr bool
	}{
		{
			name: "valid request with future due date",
			req: UpdateTaskRequest{
				Title:       stringPtrForDueDate("Valid Task"),
				Description: stringPtrForDueDate("Valid description"),
				Priority:    priorityPtr(PriorityMedium),
				DueDate:     func() *time.Time { t := time.Now().AddDate(0, 0, 1); return &t }(),
			},
			wantErr: false,
		},
		{
			name: "request with past due date should fail validation",
			req: UpdateTaskRequest{
				Title:       stringPtrForDueDate("Task with Past Due Date"),
				Description: stringPtrForDueDate("This should fail validation"),
				Priority:    priorityPtr(PriorityMedium),
				DueDate:     func() *time.Time { t := time.Now().AddDate(0, 0, -1); return &t }(),
			},
			wantErr: true,
		},
		{
			name: "request with due date exactly now should fail validation",
			req: UpdateTaskRequest{
				Title:       stringPtrForDueDate("Task with Current Due Date"),
				Description: stringPtrForDueDate("This should fail validation"),
				Priority:    priorityPtr(PriorityMedium),
				DueDate:     func() *time.Time { t := time.Now(); return &t }(),
			},
			wantErr: true,
		},
		{
			name: "request with no due date should pass validation",
			req: UpdateTaskRequest{
				Title:       stringPtrForDueDate("Task without Due Date"),
				Description: stringPtrForDueDate("This should pass validation"),
				Priority:    priorityPtr(PriorityMedium),
				DueDate:     nil,
			},
			wantErr: false,
		},
		{
			name: "request with only title update should pass validation",
			req: UpdateTaskRequest{
				Title: stringPtrForDueDate("Updated Title"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateTaskRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Helper function for string pointers
func stringPtrForDueDate(s string) *string {
	return &s
}

// Helper function for priority pointers
func priorityPtr(p TaskPriority) *TaskPriority {
	return &p
}
