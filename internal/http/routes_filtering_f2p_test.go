package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"task-manager/internal/domain"
	"testing"
)

func TestTaskFiltering_F2P(t *testing.T) {
	mockTaskService := newMockTaskService()
	mockCategoryService := newMockCategoryService()
	handler := &Handler{
		taskService:     mockTaskService,
		categoryService: mockCategoryService,
	}
	router := handler.SetupRoutes()

	// Create test tasks with different statuses and priorities
	testTasks := []*domain.Task{
		{
			ID:       1,
			Title:    "High Priority Todo",
			Status:   domain.StatusTodo,
			Priority: domain.PriorityHigh,
		},
		{
			ID:       2,
			Title:    "Critical Doing",
			Status:   domain.StatusDoing,
			Priority: domain.PriorityCritical,
		},
		{
			ID:       3,
			Title:    "Low Priority Done",
			Status:   domain.StatusDone,
			Priority: domain.PriorityLow,
		},
		{
			ID:       4,
			Title:    "Medium Priority Todo",
			Status:   domain.StatusTodo,
			Priority: domain.PriorityMedium,
		},
	}

	// Add test tasks to mock service
	for _, task := range testTasks {
		mockTaskService.tasks[task.ID] = task
	}

	tests := []struct {
		name           string
		url            string
		expectedCount  int
		expectedStatus int
		description    string
	}{
		{
			name:           "filter_by_status_todo",
			url:            "/v1/tasks?status=todo",
			expectedCount:  2,
			expectedStatus: http.StatusOK,
			description:    "Should return only todo tasks",
		},
		{
			name:           "filter_by_priority_high",
			url:            "/v1/tasks?priority=high",
			expectedCount:  1,
			expectedStatus: http.StatusOK,
			description:    "Should return only high priority tasks",
		},
		{
			name:           "filter_by_status_and_priority",
			url:            "/v1/tasks?status=todo&priority=high",
			expectedCount:  1,
			expectedStatus: http.StatusOK,
			description:    "Should return tasks matching both filters",
		},
		{
			name:           "filter_by_multiple_statuses",
			url:            "/v1/tasks?status=todo,doing",
			expectedCount:  3,
			expectedStatus: http.StatusOK,
			description:    "Should return tasks with todo or doing status",
		},
		{
			name:           "filter_by_multiple_priorities",
			url:            "/v1/tasks?priority=high,critical",
			expectedCount:  2,
			expectedStatus: http.StatusOK,
			description:    "Should return tasks with high or critical priority",
		},
		{
			name:           "filter_by_invalid_status",
			url:            "/v1/tasks?status=invalid",
			expectedCount:  0,
			expectedStatus: http.StatusBadRequest,
			description:    "Should return 400 error for invalid status",
		},
		{
			name:           "filter_by_invalid_priority",
			url:            "/v1/tasks?priority=invalid",
			expectedCount:  0,
			expectedStatus: http.StatusBadRequest,
			description:    "Should return 400 error for invalid priority",
		},
		{
			name:           "no_filters_returns_all",
			url:            "/v1/tasks",
			expectedCount:  4,
			expectedStatus: http.StatusOK,
			description:    "Should return all tasks when no filters provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. %s", tt.expectedStatus, w.Code, tt.description)
				return
			}

			if tt.expectedStatus == http.StatusOK {
				var response []domain.Task
				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
					return
				}

				if len(response) != tt.expectedCount {
					t.Errorf("Expected %d tasks, got %d. %s", tt.expectedCount, len(response), tt.description)
				}
			}
		})
	}
}

func TestTaskFilteringIntegration_F2P(t *testing.T) {
	mockTaskService := newMockTaskService()
	mockCategoryService := newMockCategoryService()
	handler := &Handler{
		taskService:     mockTaskService,
		categoryService: mockCategoryService,
	}
	router := handler.SetupRoutes()

	// Create a task with specific status and priority
	task := &domain.Task{
		ID:       1,
		Title:    "Test Task",
		Status:   domain.StatusTodo,
		Priority: domain.PriorityHigh,
	}
	mockTaskService.tasks[task.ID] = task

	t.Run("filter_should_return_correct_task", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/v1/tasks?status=todo&priority=high", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
			return
		}

		var response []domain.Task
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
			return
		}

		if len(response) != 1 {
			t.Errorf("Expected 1 task, got %d", len(response))
			return
		}

		if response[0].Title != "Test Task" {
			t.Errorf("Expected task title 'Test Task', got '%s'", response[0].Title)
		}
	})
}

func TestTaskFilteringEdgeCases_F2P(t *testing.T) {
	mockTaskService := newMockTaskService()
	mockCategoryService := newMockCategoryService()
	handler := &Handler{
		taskService:     mockTaskService,
		categoryService: mockCategoryService,
	}
	router := handler.SetupRoutes()

	tests := []struct {
		name           string
		url            string
		expectedStatus int
		description    string
	}{
		{
			name:           "empty_status_filter",
			url:            "/v1/tasks?status=",
			expectedStatus: http.StatusOK,
			description:    "Empty status filter should be ignored",
		},
		{
			name:           "empty_priority_filter",
			url:            "/v1/tasks?priority=",
			expectedStatus: http.StatusOK,
			description:    "Empty priority filter should be ignored",
		},
		{
			name:           "whitespace_in_filters",
			url:            "/v1/tasks?status=todo,doing",
			expectedStatus: http.StatusOK,
			description:    "Whitespace in filters should be trimmed",
		},
		{
			name:           "comma_separated_values",
			url:            "/v1/tasks?status=todo,,doing",
			expectedStatus: http.StatusOK,
			description:    "Empty values in comma-separated list should be ignored",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. %s", tt.expectedStatus, w.Code, tt.description)
			}
		})
	}
}
