package domain

import (
	"testing"
)

func TestTaskFiltersSearch_Validate(t *testing.T) {
	tests := []struct {
		name        string
		filters     TaskFilters
		expectError bool
		description string
	}{
		{
			name: "valid_filters_with_search",
			filters: TaskFilters{
				Statuses: []TaskStatus{StatusTodo, StatusDoing},
				Search:   "urgent",
			},
			expectError: false,
			description: "Valid filters with search should pass validation",
		},
		{
			name: "valid_filters_with_only_search",
			filters: TaskFilters{
				Search: "meeting",
			},
			expectError: false,
			description: "Valid filters with only search should pass validation",
		},
		{
			name: "valid_filters_with_empty_search",
			filters: TaskFilters{
				Statuses: []TaskStatus{StatusTodo},
				Search:   "",
			},
			expectError: false,
			description: "Valid filters with empty search should pass validation",
		},
		{
			name: "invalid_status_with_search",
			filters: TaskFilters{
				Statuses: []TaskStatus{"invalid"},
				Search:   "urgent",
			},
			expectError: true,
			description: "Invalid status with search should fail validation",
		},
		{
			name: "invalid_priority_with_search",
			filters: TaskFilters{
				Priorities: []TaskPriority{"invalid"},
				Search:     "urgent",
			},
			expectError: true,
			description: "Invalid priority with search should fail validation",
		},
		{
			name: "search_with_special_characters",
			filters: TaskFilters{
				Search: "test@email.com",
			},
			expectError: false,
			description: "Search with special characters should pass validation",
		},
		{
			name: "search_with_whitespace",
			filters: TaskFilters{
				Search: "  urgent  ",
			},
			expectError: false,
			description: "Search with whitespace should pass validation",
		},
		{
			name: "search_with_multiple_words",
			filters: TaskFilters{
				Search: "urgent meeting",
			},
			expectError: false,
			description: "Search with multiple words should pass validation",
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

func TestTaskFiltersSearch_EmptySearch(t *testing.T) {
	filters := TaskFilters{
		Search: "",
	}

	err := filters.Validate()
	if err != nil {
		t.Errorf("Empty search should not cause validation error, got: %v", err)
	}
}

func TestTaskFiltersSearch_OnlySearch(t *testing.T) {
	filters := TaskFilters{
		Search: "urgent",
	}

	err := filters.Validate()
	if err != nil {
		t.Errorf("Search-only filters should pass validation, got: %v", err)
	}
}

func TestTaskFiltersSearch_CombinedFilters(t *testing.T) {
	filters := TaskFilters{
		Statuses:   []TaskStatus{StatusTodo, StatusDoing},
		Priorities: []TaskPriority{PriorityHigh, PriorityCritical},
		Search:     "urgent",
	}

	err := filters.Validate()
	if err != nil {
		t.Errorf("Combined filters with search should pass validation, got: %v", err)
	}
}

func TestTaskFiltersSearch_SearchWithInvalidStatus(t *testing.T) {
	filters := TaskFilters{
		Statuses: []TaskStatus{StatusTodo, "invalid"},
		Search:   "urgent",
	}

	err := filters.Validate()
	if err == nil {
		t.Error("Search with invalid status should cause validation error")
	}
}

func TestTaskFiltersSearch_SearchWithInvalidPriority(t *testing.T) {
	filters := TaskFilters{
		Priorities: []TaskPriority{PriorityHigh, "invalid"},
		Search:     "urgent",
	}

	err := filters.Validate()
	if err == nil {
		t.Error("Search with invalid priority should cause validation error")
	}
}
