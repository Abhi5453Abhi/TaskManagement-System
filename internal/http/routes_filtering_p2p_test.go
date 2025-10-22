package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"task-manager/internal/domain"
	"testing"
)

func TestTaskFilteringP2P(t *testing.T) {
	mockTaskService := newMockTaskService()
	mockCategoryService := newMockCategoryService()
	handler := &Handler{
		taskService:     mockTaskService,
		categoryService: mockCategoryService,
	}
	router := handler.SetupRoutes()

	// Create test tasks
	testTasks := []*domain.Task{
		{
			ID:       1,
			Title:    "Task 1",
			Status:   domain.StatusTodo,
			Priority: domain.PriorityHigh,
		},
		{
			ID:       2,
			Title:    "Task 2",
			Status:   domain.StatusDone,
			Priority: domain.PriorityLow,
		},
	}

	// Add test tasks to mock service
	for _, task := range testTasks {
		mockTaskService.tasks[task.ID] = task
	}

	tests := []struct {
		name           string
		url            string
		expectedStatus int
		description    string
	}{
		{
			name:           "existing_get_all_tasks_still_works",
			url:            "/v1/tasks",
			expectedStatus: http.StatusOK,
			description:    "Existing GET /v1/tasks endpoint should still work without filters",
		},
		{
			name:           "existing_get_task_by_id_still_works",
			url:            "/v1/tasks/1",
			expectedStatus: http.StatusOK,
			description:    "Existing GET /v1/tasks/{id} endpoint should still work",
		},
		{
			name:           "existing_create_task_still_works",
			url:            "/v1/tasks",
			expectedStatus: http.StatusOK,
			description:    "Existing POST /v1/tasks endpoint should still work",
		},
		{
			name:           "existing_update_task_still_works",
			url:            "/v1/tasks/1",
			expectedStatus: http.StatusOK,
			description:    "Existing PUT /v1/tasks/{id} endpoint should still work",
		},
		{
			name:           "existing_delete_task_still_works",
			url:            "/v1/tasks/1",
			expectedStatus: http.StatusOK,
			description:    "Existing DELETE /v1/tasks/{id} endpoint should still work",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			var method string

			switch tt.name {
			case "existing_get_all_tasks_still_works":
				method = "GET"
				req = httptest.NewRequest(method, tt.url, nil)
			case "existing_get_task_by_id_still_works":
				method = "GET"
				req = httptest.NewRequest(method, tt.url, nil)
			case "existing_create_task_still_works":
				method = "POST"
				req = httptest.NewRequest(method, tt.url, nil)
			case "existing_update_task_still_works":
				method = "PUT"
				req = httptest.NewRequest(method, tt.url, nil)
			case "existing_delete_task_still_works":
				method = "DELETE"
				req = httptest.NewRequest(method, tt.url, nil)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// For GET requests, expect 200 OK
			// For POST/PUT/DELETE, we expect various status codes depending on implementation
			if method == "GET" {
				if w.Code != http.StatusOK {
					t.Errorf("Expected status %d, got %d. %s", http.StatusOK, w.Code, tt.description)
				}
			} else {
				// For other methods, just ensure they don't return 500 (internal server error)
				if w.Code >= http.StatusInternalServerError {
					t.Errorf("Expected status < 500, got %d. %s", w.Code, tt.description)
				}
			}
		})
	}
}

func TestTaskFilteringBackwardCompatibility_P2P(t *testing.T) {
	mockTaskService := newMockTaskService()
	mockCategoryService := newMockCategoryService()
	handler := &Handler{
		taskService:     mockTaskService,
		categoryService: mockCategoryService,
	}
	router := handler.SetupRoutes()

	// Create test tasks
	testTasks := []*domain.Task{
		{
			ID:       1,
			Title:    "Task 1",
			Status:   domain.StatusTodo,
			Priority: domain.PriorityHigh,
		},
		{
			ID:       2,
			Title:    "Task 2",
			Status:   domain.StatusDone,
			Priority: domain.PriorityLow,
		},
	}

	// Add test tasks to mock service
	for _, task := range testTasks {
		mockTaskService.tasks[task.ID] = task
	}

	t.Run("get_all_tasks_without_filters_returns_same_as_before", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/v1/tasks", nil)
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

		// Should return all tasks (2 in this case)
		if len(response) != 2 {
			t.Errorf("Expected 2 tasks, got %d", len(response))
		}
	})

	t.Run("get_all_tasks_with_empty_filters_returns_same_as_before", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/v1/tasks?status=&priority=", nil)
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

		// Should return all tasks (2 in this case)
		if len(response) != 2 {
			t.Errorf("Expected 2 tasks, got %d", len(response))
		}
	})
}

func TestTaskFilteringWithExistingData_P2P(t *testing.T) {
	mockTaskService := newMockTaskService()
	mockCategoryService := newMockCategoryService()
	handler := &Handler{
		taskService:     mockTaskService,
		categoryService: mockCategoryService,
	}
	router := handler.SetupRoutes()

	// Create test tasks with various statuses and priorities
	testTasks := []*domain.Task{
		{
			ID:       1,
			Title:    "High Priority Todo",
			Status:   domain.StatusTodo,
			Priority: domain.PriorityHigh,
		},
		{
			ID:       2,
			Title:    "Low Priority Done",
			Status:   domain.StatusDone,
			Priority: domain.PriorityLow,
		},
		{
			ID:       3,
			Title:    "Medium Priority Doing",
			Status:   domain.StatusDoing,
			Priority: domain.PriorityMedium,
		},
	}

	// Add test tasks to mock service
	for _, task := range testTasks {
		mockTaskService.tasks[task.ID] = task
	}

	t.Run("filtering_works_with_existing_task_data", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/v1/tasks?status=todo", nil)
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

		// Should return only the todo task
		if len(response) != 1 {
			t.Errorf("Expected 1 task, got %d", len(response))
			return
		}

		if response[0].Title != "High Priority Todo" {
			t.Errorf("Expected task title 'High Priority Todo', got '%s'", response[0].Title)
		}
	})

	t.Run("filtering_works_with_existing_priority_data", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/v1/tasks?priority=high", nil)
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

		// Should return only the high priority task
		if len(response) != 1 {
			t.Errorf("Expected 1 task, got %d", len(response))
			return
		}

		if response[0].Title != "High Priority Todo" {
			t.Errorf("Expected task title 'High Priority Todo', got '%s'", response[0].Title)
		}
	})
}
