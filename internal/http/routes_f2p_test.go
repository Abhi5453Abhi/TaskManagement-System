package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"task-manager/internal/domain"
	"testing"
)

// F2P Tests - Feature to Product tests
// These tests verify that the priority feature works end-to-end

func TestCreateTaskWithPriority_F2P(t *testing.T) {
	mockService := newMockTaskService()
	handler := NewHandler(mockService)
	router := handler.SetupRoutes()

	tests := []struct {
		name             string
		requestBody      domain.CreateTaskRequest
		expectedStatus   int
		expectedPriority domain.TaskPriority
	}{
		{
			name: "should create task with critical priority",
			requestBody: domain.CreateTaskRequest{
				Title:       "Critical Task",
				Description: "This is a critical task",
				Priority:    domain.PriorityCritical,
			},
			expectedStatus:   http.StatusCreated,
			expectedPriority: domain.PriorityCritical,
		},
		{
			name: "should create task with high priority",
			requestBody: domain.CreateTaskRequest{
				Title:       "High Priority Task",
				Description: "This is a high priority task",
				Priority:    domain.PriorityHigh,
			},
			expectedStatus:   http.StatusCreated,
			expectedPriority: domain.PriorityHigh,
		},
		{
			name: "should create task with medium priority",
			requestBody: domain.CreateTaskRequest{
				Title:       "Medium Priority Task",
				Description: "This is a medium priority task",
				Priority:    domain.PriorityMedium,
			},
			expectedStatus:   http.StatusCreated,
			expectedPriority: domain.PriorityMedium,
		},
		{
			name: "should create task with low priority",
			requestBody: domain.CreateTaskRequest{
				Title:       "Low Priority Task",
				Description: "This is a low priority task",
				Priority:    domain.PriorityLow,
			},
			expectedStatus:   http.StatusCreated,
			expectedPriority: domain.PriorityLow,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/v1/tasks", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if w.Code == http.StatusCreated {
				var response domain.Task
				json.Unmarshal(w.Body.Bytes(), &response)

				if response.Priority != tt.expectedPriority {
					t.Errorf("Expected priority %s, got %s", tt.expectedPriority, response.Priority)
				}

				if response.Title != tt.requestBody.Title {
					t.Errorf("Expected title %s, got %s", tt.requestBody.Title, response.Title)
				}
			}
		})
	}
}

func TestUpdateTaskPriority_F2P(t *testing.T) {
	mockService := newMockTaskService()
	handler := NewHandler(mockService)
	router := handler.SetupRoutes()

	// Create a task first
	createReq := domain.CreateTaskRequest{
		Title:       "Test Task",
		Description: "Test Description",
		Priority:    domain.PriorityMedium,
	}
	createdTask, _ := mockService.CreateTask(&createReq)

	tests := []struct {
		name             string
		taskID           int64
		updateRequest    domain.UpdateTaskRequest
		expectedStatus   int
		expectedPriority domain.TaskPriority
	}{
		{
			name:   "should update task priority to critical",
			taskID: createdTask.ID,
			updateRequest: domain.UpdateTaskRequest{
				Priority: func() *domain.TaskPriority { p := domain.PriorityCritical; return &p }(),
			},
			expectedStatus:   http.StatusOK,
			expectedPriority: domain.PriorityCritical,
		},
		{
			name:   "should update task priority to high",
			taskID: createdTask.ID,
			updateRequest: domain.UpdateTaskRequest{
				Priority: func() *domain.TaskPriority { p := domain.PriorityHigh; return &p }(),
			},
			expectedStatus:   http.StatusOK,
			expectedPriority: domain.PriorityHigh,
		},
		{
			name:   "should update task priority to low",
			taskID: createdTask.ID,
			updateRequest: domain.UpdateTaskRequest{
				Priority: func() *domain.TaskPriority { p := domain.PriorityLow; return &p }(),
			},
			expectedStatus:   http.StatusOK,
			expectedPriority: domain.PriorityLow,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.updateRequest)
			req := httptest.NewRequest("PATCH", fmt.Sprintf("/v1/tasks/%d", tt.taskID), bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if w.Code == http.StatusOK {
				var response domain.Task
				json.Unmarshal(w.Body.Bytes(), &response)

				if response.Priority != tt.expectedPriority {
					t.Errorf("Expected priority %s, got %s", tt.expectedPriority, response.Priority)
				}
			}
		})
	}
}

func TestGetTaskWithPriority_F2P(t *testing.T) {
	mockService := newMockTaskService()
	handler := NewHandler(mockService)
	router := handler.SetupRoutes()

	// Create tasks with different priorities
	priorities := []domain.TaskPriority{
		domain.PriorityCritical,
		domain.PriorityHigh,
		domain.PriorityMedium,
		domain.PriorityLow,
	}

	for i, priority := range priorities {
		createReq := domain.CreateTaskRequest{
			Title:       "Task " + string(rune(i+1)),
			Description: "Description " + string(rune(i+1)),
			Priority:    priority,
		}
		mockService.CreateTask(&createReq)
	}

	t.Run("should return task with correct priority information", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/v1/tasks/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response domain.Task
		json.Unmarshal(w.Body.Bytes(), &response)

		if response.Priority != domain.PriorityCritical {
			t.Errorf("Expected priority %s, got %s", domain.PriorityCritical, response.Priority)
		}
	})
}
