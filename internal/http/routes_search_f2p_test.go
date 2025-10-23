package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"task-manager/internal/domain"
	"testing"
)

func TestTaskSearch_F2P(t *testing.T) {
	mockTaskService := newMockTaskService()
	mockCategoryService := newMockCategoryService()
	handler := &Handler{
		taskService:     mockTaskService,
		categoryService: mockCategoryService,
	}
	router := handler.SetupRoutes()

	// Create test tasks with different titles and descriptions
	testTasks := []*domain.Task{
		{
			ID:          1,
			Title:       "Urgent Meeting Preparation",
			Description: "Prepare for the urgent client meeting tomorrow",
			Status:      domain.StatusTodo,
			Priority:    domain.PriorityHigh,
		},
		{
			ID:          2,
			Title:       "Code Review",
			Description: "Review the urgent bug fix implementation",
			Status:      domain.StatusDoing,
			Priority:    domain.PriorityMedium,
		},
		{
			ID:          3,
			Title:       "Documentation Update",
			Description: "Update project documentation",
			Status:      domain.StatusDone,
			Priority:    domain.PriorityLow,
		},
		{
			ID:          4,
			Title:       "Team Meeting",
			Description: "Weekly team standup meeting",
			Status:      domain.StatusTodo,
			Priority:    domain.PriorityMedium,
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
			name:           "search_by_title_keyword",
			url:            "/v1/tasks?search=urgent",
			expectedCount:  2,
			expectedStatus: http.StatusOK,
			description:    "Should return tasks with 'urgent' in title or description",
		},
		{
			name:           "search_by_description_keyword",
			url:            "/v1/tasks?search=meeting",
			expectedCount:  2,
			expectedStatus: http.StatusOK,
			description:    "Should return tasks with 'meeting' in title or description",
		},
		{
			name:           "search_by_partial_keyword",
			url:            "/v1/tasks?search=urg",
			expectedCount:  2,
			expectedStatus: http.StatusOK,
			description:    "Should return tasks with partial keyword match",
		},
		{
			name:           "search_case_insensitive",
			url:            "/v1/tasks?search=URGENT",
			expectedCount:  2,
			expectedStatus: http.StatusOK,
			description:    "Should return tasks with case-insensitive search",
		},
		{
			name:           "search_no_results",
			url:            "/v1/tasks?search=nonexistent",
			expectedCount:  0,
			expectedStatus: http.StatusOK,
			description:    "Should return empty array when no matches found",
		},
		{
			name:           "search_empty_string",
			url:            "/v1/tasks?search=",
			expectedCount:  4,
			expectedStatus: http.StatusOK,
			description:    "Should return all tasks when search is empty",
		},
		{
			name:           "search_combined_with_status_filter",
			url:            "/v1/tasks?search=urgent&status=todo",
			expectedCount:  1,
			expectedStatus: http.StatusOK,
			description:    "Should return urgent todo tasks only",
		},
		{
			name:           "search_combined_with_priority_filter",
			url:            "/v1/tasks?search=meeting&priority=high",
			expectedCount:  1,
			expectedStatus: http.StatusOK,
			description:    "Should return high priority meeting tasks only",
		},
		{
			name:           "search_combined_with_multiple_filters",
			url:            "/v1/tasks?search=urgent&status=todo,doing&priority=high,medium",
			expectedCount:  2,
			expectedStatus: http.StatusOK,
			description:    "Should return urgent tasks with specified status and priority",
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

func TestTaskSearchIntegration_F2P(t *testing.T) {
	mockTaskService := newMockTaskService()
	mockCategoryService := newMockCategoryService()
	handler := &Handler{
		taskService:     mockTaskService,
		categoryService: mockCategoryService,
	}
	router := handler.SetupRoutes()

	// Create a task with specific content
	task := &domain.Task{
		ID:          1,
		Title:       "Important Project Task",
		Description: "This is an important task that needs urgent attention",
		Status:      domain.StatusTodo,
		Priority:    domain.PriorityHigh,
	}
	mockTaskService.tasks[task.ID] = task

	t.Run("search_should_return_correct_task", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/v1/tasks?search=important", nil)
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

		if response[0].Title != "Important Project Task" {
			t.Errorf("Expected task title 'Important Project Task', got '%s'", response[0].Title)
		}
	})
}

func TestTaskSearchEdgeCases_F2P(t *testing.T) {
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
			name:           "search_with_special_characters",
			url:            "/v1/tasks?search=test@email.com",
			expectedStatus: http.StatusOK,
			description:    "Should handle special characters in search",
		},
		{
			name:           "search_with_whitespace",
			url:            "/v1/tasks?search=urgent",
			expectedStatus: http.StatusOK,
			description:    "Should trim whitespace from search terms",
		},
		{
			name:           "search_with_multiple_words",
			url:            "/v1/tasks?search=urgent",
			expectedStatus: http.StatusOK,
			description:    "Should handle multiple words in search",
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
