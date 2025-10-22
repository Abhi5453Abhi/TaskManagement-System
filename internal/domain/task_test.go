package domain

import (
	"testing"
)

func TestTask_Validate(t *testing.T) {
	tests := []struct {
		name    string
		task    Task
		wantErr bool
	}{
		{
			name: "valid task",
			task: Task{
				Title:       "Valid Task",
				Description: "Valid description",
				Status:      StatusTodo,
				Priority:    PriorityMedium,
			},
			wantErr: false,
		},
		{
			name: "empty title",
			task: Task{
				Title:       "",
				Description: "Valid description",
				Status:      StatusTodo,
				Priority:    PriorityMedium,
			},
			wantErr: true,
		},
		{
			name: "title too long",
			task: Task{
				Title:       string(make([]byte, 201)),
				Description: "Valid description",
				Status:      StatusTodo,
				Priority:    PriorityMedium,
			},
			wantErr: true,
		},
		{
			name: "description too long",
			task: Task{
				Title:       "Valid Task",
				Description: string(make([]byte, 1001)),
				Status:      StatusTodo,
				Priority:    PriorityMedium,
			},
			wantErr: true,
		},
		{
			name: "invalid status",
			task: Task{
				Title:       "Valid Task",
				Description: "Valid description",
				Status:      "invalid",
				Priority:    PriorityMedium,
			},
			wantErr: true,
		},
		{
			name: "invalid priority",
			task: Task{
				Title:       "Valid Task",
				Description: "Valid description",
				Status:      StatusTodo,
				Priority:    "invalid",
			},
			wantErr: true,
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

func TestCreateTaskRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateTaskRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CreateTaskRequest{
				Title:       "Valid Task",
				Description: "Valid description",
				Priority:    PriorityMedium,
			},
			wantErr: false,
		},
		{
			name: "empty title",
			req: CreateTaskRequest{
				Title:       "",
				Description: "Valid description",
				Priority:    PriorityMedium,
			},
			wantErr: true,
		},
		{
			name: "invalid priority",
			req: CreateTaskRequest{
				Title:       "Valid Task",
				Description: "Valid description",
				Priority:    "invalid",
			},
			wantErr: true,
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

func TestGetPriorityOrder(t *testing.T) {
	tests := []struct {
		priority TaskPriority
		want     int
	}{
		{PriorityCritical, 4},
		{PriorityHigh, 3},
		{PriorityMedium, 2},
		{PriorityLow, 1},
		{"invalid", 0},
	}

	for _, tt := range tests {
		t.Run(string(tt.priority), func(t *testing.T) {
			if got := GetPriorityOrder(tt.priority); got != tt.want {
				t.Errorf("GetPriorityOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidStatus(t *testing.T) {
	tests := []struct {
		status TaskStatus
		want   bool
	}{
		{StatusTodo, true},
		{StatusDoing, true},
		{StatusDone, true},
		{"invalid", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if got := isValidStatus(tt.status); got != tt.want {
				t.Errorf("isValidStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidPriority(t *testing.T) {
	tests := []struct {
		priority TaskPriority
		want     bool
	}{
		{PriorityLow, true},
		{PriorityMedium, true},
		{PriorityHigh, true},
		{PriorityCritical, true},
		{"invalid", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.priority), func(t *testing.T) {
			if got := isValidPriority(tt.priority); got != tt.want {
				t.Errorf("isValidPriority() = %v, want %v", got, tt.want)
			}
		})
	}
}
