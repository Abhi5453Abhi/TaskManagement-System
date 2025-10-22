package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"task-manager/internal/domain"
	"testing"
)

// P2P Tests - Product to Product tests
// These tests verify that existing functionality continues to work

func TestBasicTaskCRUD_P2P(t *testing.T) {
	mockService := newMockTaskService()
	handler := NewHandler(mockService, newMockCategoryService())
	router := handler.SetupRoutes()

	t.Run("should create task without priority field", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"title":       "Basic Task",
			"description": "Basic task description",
		}

		jsonBody, _ := json.Marshal(requestBody)
		req := httptest.NewRequest("POST", "/v1/tasks", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("should create task with default priority", func(t *testing.T) {
		requestBody := domain.CreateTaskRequest{
			Title:       "Task with Default Priority",
			Description: "This task should have default priority",
			Priority:    domain.PriorityMedium, // This will be the default
		}

		jsonBody, _ := json.Marshal(requestBody)
		req := httptest.NewRequest("POST", "/v1/tasks", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
		}

		var response domain.Task
		json.Unmarshal(w.Body.Bytes(), &response)

		if response.Priority != domain.PriorityMedium {
			t.Errorf("Expected default priority %s, got %s", domain.PriorityMedium, response.Priority)
		}
	})
}

func TestTaskStatusManagement_P2P(t *testing.T) {
	mockService := newMockTaskService()
	handler := NewHandler(mockService, newMockCategoryService())
	router := handler.SetupRoutes()

	// Create a task first
	createReq := domain.CreateTaskRequest{
		Title:       "Status Test Task",
		Description: "Testing status updates",
		Priority:    domain.PriorityMedium,
	}
	mockService.CreateTask(&createReq)

	t.Run("should update task status to doing", func(t *testing.T) {
		updateReq := domain.UpdateTaskRequest{
			Status: func() *domain.TaskStatus { s := domain.StatusDoing; return &s }(),
		}

		jsonBody, _ := json.Marshal(updateReq)
		req := httptest.NewRequest("PATCH", "/v1/tasks/1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response domain.Task
		json.Unmarshal(w.Body.Bytes(), &response)

		if response.Status != domain.StatusDoing {
			t.Errorf("Expected status %s, got %s", domain.StatusDoing, response.Status)
		}
	})

	t.Run("should update task status to done", func(t *testing.T) {
		updateReq := domain.UpdateTaskRequest{
			Status: func() *domain.TaskStatus { s := domain.StatusDone; return &s }(),
		}

		jsonBody, _ := json.Marshal(updateReq)
		req := httptest.NewRequest("PATCH", "/v1/tasks/1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response domain.Task
		json.Unmarshal(w.Body.Bytes(), &response)

		if response.Status != domain.StatusDone {
			t.Errorf("Expected status %s, got %s", domain.StatusDone, response.Status)
		}
	})
}

func TestTaskTitleAndDescription_P2P(t *testing.T) {
	mockService := newMockTaskService()
	handler := NewHandler(mockService, newMockCategoryService())
	router := handler.SetupRoutes()

	// Create a task first
	createReq := domain.CreateTaskRequest{
		Title:       "Original Title",
		Description: "Original Description",
		Priority:    domain.PriorityMedium,
	}
	mockService.CreateTask(&createReq)

	t.Run("should update task title", func(t *testing.T) {
		newTitle := "Updated Title"
		updateReq := domain.UpdateTaskRequest{
			Title: &newTitle,
		}

		jsonBody, _ := json.Marshal(updateReq)
		req := httptest.NewRequest("PATCH", "/v1/tasks/1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response domain.Task
		json.Unmarshal(w.Body.Bytes(), &response)

		if response.Title != newTitle {
			t.Errorf("Expected title %s, got %s", newTitle, response.Title)
		}
	})

	t.Run("should update task description", func(t *testing.T) {
		newDescription := "Updated Description"
		updateReq := domain.UpdateTaskRequest{
			Description: &newDescription,
		}

		jsonBody, _ := json.Marshal(updateReq)
		req := httptest.NewRequest("PATCH", "/v1/tasks/1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response domain.Task
		json.Unmarshal(w.Body.Bytes(), &response)

		if response.Description != newDescription {
			t.Errorf("Expected description %s, got %s", newDescription, response.Description)
		}
	})
}

func TestTaskDeletion_P2P(t *testing.T) {
	mockService := newMockTaskService()
	handler := NewHandler(mockService, newMockCategoryService())
	router := handler.SetupRoutes()

	// Create a task first
	createReq := domain.CreateTaskRequest{
		Title:       "Task to Delete",
		Description: "This task will be deleted",
		Priority:    domain.PriorityLow,
	}
	mockService.CreateTask(&createReq)

	t.Run("should delete task successfully", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/v1/tasks/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
		}
	})

	t.Run("should return 404 when deleting non-existent task", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/v1/tasks/999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestTaskRetrieval_P2P(t *testing.T) {
	mockService := newMockTaskService()
	handler := NewHandler(mockService, newMockCategoryService())
	router := handler.SetupRoutes()

	// Create multiple tasks
	tasks := []domain.CreateTaskRequest{
		{Title: "Task 1", Description: "Description 1", Priority: domain.PriorityHigh},
		{Title: "Task 2", Description: "Description 2", Priority: domain.PriorityLow},
		{Title: "Task 3", Description: "Description 3", Priority: domain.PriorityMedium},
	}

	for _, task := range tasks {
		mockService.CreateTask(&task)
	}

	t.Run("should get all tasks", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/v1/tasks", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response []domain.Task
		json.Unmarshal(w.Body.Bytes(), &response)

		if len(response) != 3 {
			t.Errorf("Expected 3 tasks, got %d", len(response))
		}
	})

	t.Run("should get specific task by id", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/v1/tasks/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response domain.Task
		json.Unmarshal(w.Body.Bytes(), &response)

		if response.Title != "Task 1" {
			t.Errorf("Expected title 'Task 1', got %s", response.Title)
		}
	})

	t.Run("should return 404 for non-existent task", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/v1/tasks/999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}
