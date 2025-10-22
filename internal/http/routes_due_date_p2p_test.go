package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"task-manager/internal/domain"
	"testing"
	"time"
)

// P2P Tests for Due Date Feature - Product to Product tests
// These tests verify that existing functionality continues to work with due date feature

func TestBasicTaskCRUDWithDueDate_P2P(t *testing.T) {
	mockService := newMockTaskService()
	handler := NewHandler(mockService)
	router := handler.SetupRoutes()

	t.Run("should create task with all fields including due date", func(t *testing.T) {
		tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
		requestBody := domain.CreateTaskRequest{
			Title:       "Complete Task with Due Date",
			Description: "This task has all fields including due date",
			Priority:    domain.PriorityHigh,
			DueDate:     func() *time.Time { t, _ := time.Parse("2006-01-02", tomorrow); return &t }(),
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

		if response.Title != requestBody.Title {
			t.Errorf("Expected title %s, got %s", requestBody.Title, response.Title)
		}
		if response.Description != requestBody.Description {
			t.Errorf("Expected description %s, got %s", requestBody.Description, response.Description)
		}
		if response.Priority != requestBody.Priority {
			t.Errorf("Expected priority %s, got %s", requestBody.Priority, response.Priority)
		}
		if response.Status != domain.StatusTodo {
			t.Errorf("Expected status %s, got %s", domain.StatusTodo, response.Status)
		}
		if response.DueDate == nil {
			t.Errorf("Expected due date, got nil")
		} else {
			actualDate := response.DueDate.Format("2006-01-02")
			if actualDate != tomorrow {
				t.Errorf("Expected due date %s, got %s", tomorrow, actualDate)
			}
		}
	})
}

func TestTaskStatusManagementWithDueDate_P2P(t *testing.T) {
	mockService := newMockTaskService()
	handler := NewHandler(mockService)
	router := handler.SetupRoutes()

	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	// Create a task with due date
	createReq := domain.CreateTaskRequest{
		Title:       "Status Test Task with Due Date",
		Description: "Testing status updates with due date",
		Priority:    domain.PriorityMedium,
		DueDate:     func() *time.Time { t, _ := time.Parse("2006-01-02", tomorrow); return &t }(),
	}
	mockService.CreateTask(&createReq)

	t.Run("should update task status while preserving due date", func(t *testing.T) {
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

		// Due date should be preserved
		if response.DueDate == nil {
			t.Errorf("Expected due date to be preserved, got nil")
		} else {
			actualDate := response.DueDate.Format("2006-01-02")
			if actualDate != tomorrow {
				t.Errorf("Expected due date %s to be preserved, got %s", tomorrow, actualDate)
			}
		}
	})

	t.Run("should update task status to done while preserving due date", func(t *testing.T) {
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

		// Due date should still be preserved
		if response.DueDate == nil {
			t.Errorf("Expected due date to be preserved, got nil")
		}
	})
}

func TestTaskPriorityManagementWithDueDate_P2P(t *testing.T) {
	mockService := newMockTaskService()
	handler := NewHandler(mockService)
	router := handler.SetupRoutes()

	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	// Create a task with due date
	createReq := domain.CreateTaskRequest{
		Title:       "Priority Test Task with Due Date",
		Description: "Testing priority updates with due date",
		Priority:    domain.PriorityMedium,
		DueDate:     func() *time.Time { t, _ := time.Parse("2006-01-02", tomorrow); return &t }(),
	}
	mockService.CreateTask(&createReq)

	t.Run("should update task priority while preserving due date", func(t *testing.T) {
		updateReq := domain.UpdateTaskRequest{
			Priority: func() *domain.TaskPriority { p := domain.PriorityCritical; return &p }(),
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

		if response.Priority != domain.PriorityCritical {
			t.Errorf("Expected priority %s, got %s", domain.PriorityCritical, response.Priority)
		}

		// Due date should be preserved
		if response.DueDate == nil {
			t.Errorf("Expected due date to be preserved, got nil")
		} else {
			actualDate := response.DueDate.Format("2006-01-02")
			if actualDate != tomorrow {
				t.Errorf("Expected due date %s to be preserved, got %s", tomorrow, actualDate)
			}
		}
	})
}

func TestTaskTitleAndDescriptionWithDueDate_P2P(t *testing.T) {
	mockService := newMockTaskService()
	handler := NewHandler(mockService)
	router := handler.SetupRoutes()

	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	// Create a task with due date
	createReq := domain.CreateTaskRequest{
		Title:       "Original Title with Due Date",
		Description: "Original Description with Due Date",
		Priority:    domain.PriorityMedium,
		DueDate:     func() *time.Time { t, _ := time.Parse("2006-01-02", tomorrow); return &t }(),
	}
	mockService.CreateTask(&createReq)

	t.Run("should update task title while preserving due date", func(t *testing.T) {
		newTitle := "Updated Title with Due Date"
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

		// Due date should be preserved
		if response.DueDate == nil {
			t.Errorf("Expected due date to be preserved, got nil")
		}
	})

	t.Run("should update task description while preserving due date", func(t *testing.T) {
		newDescription := "Updated Description with Due Date"
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

		// Due date should be preserved
		if response.DueDate == nil {
			t.Errorf("Expected due date to be preserved, got nil")
		}
	})
}

func TestTaskDeletionWithDueDate_P2P(t *testing.T) {
	mockService := newMockTaskService()
	handler := NewHandler(mockService)
	router := handler.SetupRoutes()

	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	// Create a task with due date
	createReq := domain.CreateTaskRequest{
		Title:       "Task to Delete with Due Date",
		Description: "This task will be deleted and has a due date",
		Priority:    domain.PriorityLow,
		DueDate:     func() *time.Time { t, _ := time.Parse("2006-01-02", tomorrow); return &t }(),
	}
	mockService.CreateTask(&createReq)

	t.Run("should delete task with due date successfully", func(t *testing.T) {
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

func TestTaskRetrievalWithDueDate_P2P(t *testing.T) {
	mockService := newMockTaskService()
	handler := NewHandler(mockService)
	router := handler.SetupRoutes()

	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	// Create multiple tasks with and without due dates
	tasks := []domain.CreateTaskRequest{
		{Title: "Task 1 with Due Date", Description: "Description 1", Priority: domain.PriorityHigh, DueDate: func() *time.Time { t, _ := time.Parse("2006-01-02", tomorrow); return &t }()},
		{Title: "Task 2 without Due Date", Description: "Description 2", Priority: domain.PriorityLow, DueDate: nil},
		{Title: "Task 3 with Due Date", Description: "Description 3", Priority: domain.PriorityMedium, DueDate: func() *time.Time { t, _ := time.Parse("2006-01-02", tomorrow); return &t }()},
	}

	for _, task := range tasks {
		mockService.CreateTask(&task)
	}

	t.Run("should get all tasks including due date information", func(t *testing.T) {
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

		// Check that tasks with due dates come first
		if response[0].DueDate == nil {
			t.Errorf("Expected first task to have due date, got nil")
		}
		if response[1].DueDate == nil {
			t.Errorf("Expected second task to have due date, got nil")
		}
		if response[2].DueDate != nil {
			t.Errorf("Expected third task to have no due date, got %s", response[2].DueDate.Format("2006-01-02"))
		}
	})

	t.Run("should get specific task by id with due date", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/v1/tasks/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response domain.Task
		json.Unmarshal(w.Body.Bytes(), &response)

		if response.Title != "Task 1 with Due Date" {
			t.Errorf("Expected title 'Task 1 with Due Date', got %s", response.Title)
		}

		if response.DueDate == nil {
			t.Errorf("Expected due date, got nil")
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
