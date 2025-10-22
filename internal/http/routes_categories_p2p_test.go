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

func TestTaskWithCategories_P2P(t *testing.T) {
	mockTaskService := newMockTaskService()
	mockCategoryService := newMockCategoryService()
	handler := &Handler{
		taskService:     mockTaskService,
		categoryService: mockCategoryService,
	}

	router := handler.SetupRoutes()

	// Create test categories first
	categories := []*domain.Category{
		{ID: 1, Name: "Work", Description: stringPtr("Work tasks"), Color: stringPtr("#FF0000"), CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, Name: "Personal", Description: stringPtr("Personal tasks"), Color: stringPtr("#00FF00"), CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	for _, cat := range categories {
		mockCategoryService.categories[cat.ID] = cat
	}

	t.Run("create task with categories", func(t *testing.T) {
		requestBody := domain.CreateTaskRequest{
			Title:       "Task with Categories",
			Description: "This task has categories",
			Priority:    domain.PriorityHigh,
			CategoryIDs: []int64{1, 2},
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

		if len(response.Categories) != 2 {
			t.Errorf("Expected 2 categories, got %d", len(response.Categories))
		}
	})

	t.Run("create task without categories", func(t *testing.T) {
		requestBody := domain.CreateTaskRequest{
			Title:       "Task without Categories",
			Description: "This task has no categories",
			Priority:    domain.PriorityMedium,
			CategoryIDs: []int64{},
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

		if len(response.Categories) != 0 {
			t.Errorf("Expected 0 categories, got %d", len(response.Categories))
		}
	})

	t.Run("update task categories", func(t *testing.T) {
		// First create a task
		createReq := domain.CreateTaskRequest{
			Title:       "Task to Update",
			Description: "This task will be updated",
			Priority:    domain.PriorityLow,
			CategoryIDs: []int64{1},
		}

		jsonBody, _ := json.Marshal(createReq)
		req := httptest.NewRequest("POST", "/v1/tasks", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var createdTask domain.Task
		json.Unmarshal(w.Body.Bytes(), &createdTask)

		// Now update the task with different categories
		updateReq := domain.UpdateTaskRequest{
			CategoryIDs: &[]int64{2},
		}

		jsonBody, _ = json.Marshal(updateReq)
		req = httptest.NewRequest("PATCH", "/v1/tasks/1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response domain.Task
		json.Unmarshal(w.Body.Bytes(), &response)

		if len(response.Categories) != 1 {
			t.Errorf("Expected 1 category, got %d", len(response.Categories))
		}

		if response.Categories[0].ID != 2 {
			t.Errorf("Expected category ID 2, got %d", response.Categories[0].ID)
		}
	})

	t.Run("get task with categories", func(t *testing.T) {
		// Create a task with categories
		createReq := domain.CreateTaskRequest{
			Title:       "Task with Categories",
			Description: "This task has categories",
			Priority:    domain.PriorityHigh,
			CategoryIDs: []int64{1, 2},
		}

		jsonBody, _ := json.Marshal(createReq)
		req := httptest.NewRequest("POST", "/v1/tasks", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var createdTask domain.Task
		json.Unmarshal(w.Body.Bytes(), &createdTask)

		// Now get the task
		req = httptest.NewRequest("GET", "/v1/tasks/1", nil)

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response domain.Task
		json.Unmarshal(w.Body.Bytes(), &response)

		if len(response.Categories) != 2 {
			t.Errorf("Expected 2 categories, got %d", len(response.Categories))
		}
	})

	t.Run("get all tasks with categories", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/v1/tasks", nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response []domain.Task
		json.Unmarshal(w.Body.Bytes(), &response)

		if len(response) == 0 {
			t.Errorf("Expected at least one task, got %d", len(response))
		}

		// Check that all tasks have categories loaded
		for _, task := range response {
			if task.Categories == nil {
				t.Errorf("Expected categories to be loaded for task %d", task.ID)
			}
		}
	})
}

func TestCategoryManagement_P2P(t *testing.T) {
	mockCategoryService := newMockCategoryService()
	handler := &Handler{
		categoryService: mockCategoryService,
	}

	router := handler.SetupRoutes()

	t.Run("create and retrieve category", func(t *testing.T) {
		// Create a category
		createReq := domain.CreateCategoryRequest{
			Name:        "Test Category",
			Description: stringPtr("Test description"),
			Color:       stringPtr("#FF0000"),
		}

		jsonBody, _ := json.Marshal(createReq)
		req := httptest.NewRequest("POST", "/v1/categories", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
		}

		var createdCategory domain.Category
		json.Unmarshal(w.Body.Bytes(), &createdCategory)

		// Now retrieve the category
		req = httptest.NewRequest("GET", "/v1/categories/1", nil)

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response domain.Category
		json.Unmarshal(w.Body.Bytes(), &response)

		if response.Name != createdCategory.Name {
			t.Errorf("Expected name %s, got %s", createdCategory.Name, response.Name)
		}
	})

	t.Run("update category preserves other fields", func(t *testing.T) {
		// Create a category
		createReq := domain.CreateCategoryRequest{
			Name:        "Original Name",
			Description: stringPtr("Original description"),
			Color:       stringPtr("#FF0000"),
		}

		jsonBody, _ := json.Marshal(createReq)
		req := httptest.NewRequest("POST", "/v1/categories", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var createdCategory domain.Category
		json.Unmarshal(w.Body.Bytes(), &createdCategory)

		// Update only the name
		updateReq := domain.UpdateCategoryRequest{
			Name: stringPtr("Updated Name"),
		}

		jsonBody, _ = json.Marshal(updateReq)
		req = httptest.NewRequest("PUT", "/v1/categories/1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response domain.Category
		json.Unmarshal(w.Body.Bytes(), &response)

		if response.Name != "Updated Name" {
			t.Errorf("Expected name 'Updated Name', got %s", response.Name)
		}

		// Description and color should be preserved
		if response.Description == nil || *response.Description != "Original description" {
			t.Errorf("Expected description to be preserved")
		}

		if response.Color == nil || *response.Color != "#FF0000" {
			t.Errorf("Expected color to be preserved")
		}
	})

	t.Run("delete category", func(t *testing.T) {
		// Create a category
		createReq := domain.CreateCategoryRequest{
			Name: "Category to Delete",
		}

		jsonBody, _ := json.Marshal(createReq)
		req := httptest.NewRequest("POST", "/v1/categories", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Delete the category
		req = httptest.NewRequest("DELETE", "/v1/categories/1", nil)

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
		}

		// Try to get the deleted category
		req = httptest.NewRequest("GET", "/v1/categories/1", nil)

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

