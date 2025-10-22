package domain

import (
	"testing"
)

func TestTaskFilters_Validate(t *testing.T) {
	tests := []struct {
		name        string
		filters     TaskFilters
		expectError bool
		description string
	}{
		{
			name: "valid_filters_with_status",
			filters: TaskFilters{
				Statuses: []TaskStatus{StatusTodo, StatusDoing},
			},
			expectError: false,
			description: "Valid status filters should pass validation",
		},
		{
			name: "valid_filters_with_priority",
			filters: TaskFilters{
				Priorities: []TaskPriority{PriorityHigh, PriorityCritical},
			},
			expectError: false,
			description: "Valid priority filters should pass validation",
		},
		{
			name: "valid_filters_with_both",
			filters: TaskFilters{
				Statuses:   []TaskStatus{StatusTodo},
				Priorities: []TaskPriority{PriorityHigh},
			},
			expectError: false,
			description: "Valid status and priority filters should pass validation",
		},
		{
			name: "invalid_status_filter",
			filters: TaskFilters{
				Statuses: []TaskStatus{"invalid"},
			},
			expectError: true,
			description: "Invalid status should fail validation",
		},
		{
			name: "invalid_priority_filter",
			filters: TaskFilters{
				Priorities: []TaskPriority{"invalid"},
			},
			expectError: true,
			description: "Invalid priority should fail validation",
		},
		{
			name: "empty_filters",
			filters: TaskFilters{
				Statuses:   []TaskStatus{},
				Priorities: []TaskPriority{},
			},
			expectError: false,
			description: "Empty filters should pass validation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.filters.Validate()
			hasError := err != nil

			if hasError != tt.expectError {
				t.Errorf("TaskFilters.Validate() error = %v, want error = %v. %s", 
					hasError, tt.expectError, tt.description)
			}
		})
	}
}

func TestTaskFilters_EmptyFilters(t *testing.T) {
	filters := TaskFilters{}

	err := filters.Validate()
	if err != nil {
		t.Errorf("Empty filters should not cause validation error, got: %v", err)
	}
}

func TestTaskFilters_MixedValidInvalid(t *testing.T) {
	filters := TaskFilters{
		Statuses:   []TaskStatus{StatusTodo, "invalid"},
		Priorities: []TaskPriority{PriorityHigh, "invalid"},
	}

	err := filters.Validate()
	if err == nil {
		t.Error("Mixed valid and invalid filters should cause validation error")
	}
}

func TestTaskFilters_AllValidStatuses(t *testing.T) {
	filters := TaskFilters{
		Statuses: []TaskStatus{StatusTodo, StatusDoing, StatusDone},
	}

	err := filters.Validate()
	if err != nil {
		t.Errorf("All valid statuses should pass validation, got: %v", err)
	}
}

func TestTaskFilters_AllValidPriorities(t *testing.T) {
	filters := TaskFilters{
		Priorities: []TaskPriority{PriorityLow, PriorityMedium, PriorityHigh, PriorityCritical},
	}

	err := filters.Validate()
	if err != nil {
		t.Errorf("All valid priorities should pass validation, got: %v", err)
	}
}
