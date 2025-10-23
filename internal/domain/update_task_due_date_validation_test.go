package domain

import (
	"testing"
	"time"
)

func TestUpdateTaskRequestDueDateValidation(t *testing.T) {
	tests := []struct {
		name    string
		req     UpdateTaskRequest
		wantErr bool
	}{
		{
			name: "valid request with future due date",
			req: UpdateTaskRequest{
				Title:       stringPtrForUpdate("Valid Task"),
				Description: stringPtrForUpdate("Valid description"),
				Priority:    priorityPtrForUpdate(PriorityMedium),
				DueDate:     func() *time.Time { t := time.Now().AddDate(0, 0, 1); return &t }(),
			},
			wantErr: false,
		},
		{
			name: "request with past due date should fail validation",
			req: UpdateTaskRequest{
				Title:       stringPtrForUpdate("Task with Past Due Date"),
				Description: stringPtrForUpdate("This should fail validation"),
				Priority:    priorityPtrForUpdate(PriorityMedium),
				DueDate:     func() *time.Time { t := time.Now().AddDate(0, 0, -1); return &t }(),
			},
			wantErr: true,
		},
		{
			name: "request with due date exactly now should fail validation",
			req: UpdateTaskRequest{
				Title:       stringPtrForUpdate("Task with Current Due Date"),
				Description: stringPtrForUpdate("This should fail validation"),
				Priority:    priorityPtrForUpdate(PriorityMedium),
				DueDate:     func() *time.Time { t := time.Now(); return &t }(),
			},
			wantErr: true,
		},
		{
			name: "request with no due date should pass validation",
			req: UpdateTaskRequest{
				Title:       stringPtrForUpdate("Task without Due Date"),
				Description: stringPtrForUpdate("This should pass validation"),
				Priority:    priorityPtrForUpdate(PriorityMedium),
				DueDate:     nil,
			},
			wantErr: false,
		},
		{
			name: "request with only title update should pass validation",
			req: UpdateTaskRequest{
				Title: stringPtrForUpdate("Updated Title"),
			},
			wantErr: false,
		},
		{
			name: "request with only status update should pass validation",
			req: UpdateTaskRequest{
				Status: statusPtrForUpdate(StatusDoing),
			},
			wantErr: false,
		},
		{
			name: "request with only priority update should pass validation",
			req: UpdateTaskRequest{
				Priority: priorityPtrForUpdate(PriorityHigh),
			},
			wantErr: false,
		},
		{
			name: "request with only description update should pass validation",
			req: UpdateTaskRequest{
				Description: stringPtrForUpdate("Updated description"),
			},
			wantErr: false,
		},
		{
			name: "request with future due date and other fields should pass validation",
			req: UpdateTaskRequest{
				Title:       stringPtrForUpdate("Updated Task"),
				Description: stringPtrForUpdate("Updated description"),
				Status:      statusPtrForUpdate(StatusDoing),
				Priority:    priorityPtrForUpdate(PriorityHigh),
				DueDate:     func() *time.Time { t := time.Now().AddDate(0, 0, 2); return &t }(),
			},
			wantErr: false,
		},
		{
			name: "request with past due date and other fields should fail validation",
			req: UpdateTaskRequest{
				Title:       stringPtrForUpdate("Updated Task"),
				Description: stringPtrForUpdate("Updated description"),
				Status:      statusPtrForUpdate(StatusDoing),
				Priority:    priorityPtrForUpdate(PriorityHigh),
				DueDate:     func() *time.Time { t := time.Now().AddDate(0, 0, -2); return &t }(),
			},
			wantErr: true,
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

func TestUpdateTaskRequestDueDateValidationEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		req     UpdateTaskRequest
		wantErr bool
	}{
		{
			name: "request with due date in the past by hours should fail validation",
			req: UpdateTaskRequest{
				Title:   stringPtrForUpdate("Task with Past Due Date"),
				DueDate: func() *time.Time { t := time.Now().Add(-2 * time.Hour); return &t }(),
			},
			wantErr: true,
		},
		{
			name: "request with due date in the past by minutes should fail validation",
			req: UpdateTaskRequest{
				Title:   stringPtrForUpdate("Task with Past Due Date"),
				DueDate: func() *time.Time { t := time.Now().Add(-30 * time.Minute); return &t }(),
			},
			wantErr: true,
		},
		{
			name: "request with due date in the future by minutes should pass validation",
			req: UpdateTaskRequest{
				Title:   stringPtrForUpdate("Task with Future Due Date"),
				DueDate: func() *time.Time { t := time.Now().Add(30 * time.Minute); return &t }(),
			},
			wantErr: false,
		},
		{
			name: "request with due date in the future by hours should pass validation",
			req: UpdateTaskRequest{
				Title:   stringPtrForUpdate("Task with Future Due Date"),
				DueDate: func() *time.Time { t := time.Now().Add(2 * time.Hour); return &t }(),
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

// Helper functions for UpdateTaskRequest tests
func stringPtrForUpdate(s string) *string {
	return &s
}

func priorityPtrForUpdate(p TaskPriority) *TaskPriority {
	return &p
}

func statusPtrForUpdate(s TaskStatus) *TaskStatus {
	return &s
}
