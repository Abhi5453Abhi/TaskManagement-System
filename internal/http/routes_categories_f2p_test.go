package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"task-manager/internal/domain"
	"testing"
	"time"
)

// Mock category service for testing
type mockCategoryService struct {
	categories map[int64]*domain.Category
	nextID     int64
}

func newMockCategoryService() *mockCategoryService {
	return &mockCategoryService{
		categories: make(map[int64]*domain.Category),
		nextID:     1,
	}
}

func (m *mockCategoryService) CreateCategory(req *domain.CreateCategoryRequest) (*domain.Category, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	category := &domain.Category{
		ID:          m.nextID,
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	m.categories[m.nextID] = category
	m.nextID++
	return category, nil
}

func (m *mockCategoryService) GetCategory(id int64) (*domain.Category, error) {
	category, exists := m.categories[id]
	if !exists {
		return nil, domain.ErrCategoryNotFound
	}
	return category, nil
}

func (m *mockCategoryService) GetAllCategories() ([]domain.Category, error) {
	var categories []domain.Category
	for _, category := range m.categories {
		categories = append(categories, *category)
	}
	return categories, nil
}

func (m *mockCategoryService) UpdateCategory(id int64, req *domain.UpdateCategoryRequest) (*domain.Category, error) {
	category, exists := m.categories[id]
	if !exists {
		return nil, domain.ErrCategoryNotFound
	}

	if err := req.Validate(); err != nil {
		return nil, err
	}

	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Description != nil {
		category.Description = req.Description
	}
	if req.Color != nil {
		category.Color = req.Color
	}
	category.UpdatedAt = time.Now()

	return category, nil
}

func (m *mockCategoryService) DeleteCategory(id int64) error {
	_, exists := m.categories[id]
	if !exists {
		return domain.ErrCategoryNotFound
	}
	delete(m.categories, id)
	return nil
}

func TestCreateCategory_F2P(t *testing.T) {
	mockCategoryService := newMockCategoryService()
	handler := &Handler{
		categoryService: mockCategoryService,
	}

	router := handler.SetupRoutes()

	tests := []struct {
		name           string
		requestBody    domain.CreateCategoryRequest
		expectedStatus int
	}{
		{
			name: "create category with all fields",
			requestBody: domain.CreateCategoryRequest{
				Name:        "Work",
				Description: stringPtr("Work related tasks"),
				Color:       stringPtr("#FF0000"),
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "create category with minimal fields",
			requestBody: domain.CreateCategoryRequest{
				Name: "Personal",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "create category with invalid color",
			requestBody: domain.CreateCategoryRequest{
				Name:  "Invalid",
				Color: stringPtr("red"),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "create category with empty name",
			requestBody: domain.CreateCategoryRequest{
				Name: "",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/v1/categories", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusCreated {
				var response domain.Category
				json.Unmarshal(w.Body.Bytes(), &response)

				if response.Name != tt.requestBody.Name {
					t.Errorf("Expected name %s, got %s", tt.requestBody.Name, response.Name)
				}

				if tt.requestBody.Description != nil {
					if response.Description == nil || *response.Description != *tt.requestBody.Description {
						t.Errorf("Expected description %s, got %v", *tt.requestBody.Description, response.Description)
					}
				}

				if tt.requestBody.Color != nil {
					if response.Color == nil || *response.Color != *tt.requestBody.Color {
						t.Errorf("Expected color %s, got %v", *tt.requestBody.Color, response.Color)
					}
				}
			}
		})
	}
}

func TestGetCategory_F2P(t *testing.T) {
	mockCategoryService := newMockCategoryService()
	handler := &Handler{
		categoryService: mockCategoryService,
	}

	router := handler.SetupRoutes()

	// Create a test category
	testCategory := &domain.Category{
		ID:          1,
		Name:        "Test Category",
		Description: stringPtr("Test description"),
		Color:       stringPtr("#00FF00"),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	mockCategoryService.categories[1] = testCategory

	tests := []struct {
		name           string
		categoryID     int64
		expectedStatus int
	}{
		{
			name:           "get existing category",
			categoryID:     1,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "get non-existent category",
			categoryID:     999,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", fmt.Sprintf("/v1/categories/%d", tt.categoryID), nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var response domain.Category
				json.Unmarshal(w.Body.Bytes(), &response)

				if response.Name != testCategory.Name {
					t.Errorf("Expected name %s, got %s", testCategory.Name, response.Name)
				}
			}
		})
	}
}

func TestGetAllCategories_F2P(t *testing.T) {
	mockCategoryService := newMockCategoryService()
	handler := &Handler{
		categoryService: mockCategoryService,
	}

	router := handler.SetupRoutes()

	// Create test categories
	categories := []*domain.Category{
		{ID: 1, Name: "Work", Description: stringPtr("Work tasks"), Color: stringPtr("#FF0000"), CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, Name: "Personal", Description: stringPtr("Personal tasks"), Color: stringPtr("#00FF00"), CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 3, Name: "Shopping", Description: nil, Color: nil, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	for _, cat := range categories {
		mockCategoryService.categories[cat.ID] = cat
	}

	req := httptest.NewRequest("GET", "/v1/categories", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response []domain.Category
	json.Unmarshal(w.Body.Bytes(), &response)

	if len(response) != 3 {
		t.Errorf("Expected 3 categories, got %d", len(response))
	}
}

func TestUpdateCategory_F2P(t *testing.T) {
	mockCategoryService := newMockCategoryService()
	handler := &Handler{
		categoryService: mockCategoryService,
	}

	router := handler.SetupRoutes()

	// Create a test category
	testCategory := &domain.Category{
		ID:          1,
		Name:        "Original Name",
		Description: stringPtr("Original description"),
		Color:       stringPtr("#FF0000"),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	mockCategoryService.categories[1] = testCategory

	tests := []struct {
		name           string
		categoryID     int64
		requestBody    domain.UpdateCategoryRequest
		expectedStatus int
	}{
		{
			name:       "update all fields",
			categoryID: 1,
			requestBody: domain.UpdateCategoryRequest{
				Name:        stringPtr("Updated Name"),
				Description: stringPtr("Updated description"),
				Color:       stringPtr("#00FF00"),
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:       "update only name",
			categoryID: 1,
			requestBody: domain.UpdateCategoryRequest{
				Name: stringPtr("New Name"),
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:       "update with invalid color",
			categoryID: 1,
			requestBody: domain.UpdateCategoryRequest{
				Color: stringPtr("red"),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:       "update non-existent category",
			categoryID: 999,
			requestBody: domain.UpdateCategoryRequest{
				Name: stringPtr("New Name"),
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("PUT", fmt.Sprintf("/v1/categories/%d", tt.categoryID), bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var response domain.Category
				json.Unmarshal(w.Body.Bytes(), &response)

				if tt.requestBody.Name != nil && response.Name != *tt.requestBody.Name {
					t.Errorf("Expected name %s, got %s", *tt.requestBody.Name, response.Name)
				}
			}
		})
	}
}

func TestDeleteCategory_F2P(t *testing.T) {
	mockCategoryService := newMockCategoryService()
	handler := &Handler{
		categoryService: mockCategoryService,
	}

	router := handler.SetupRoutes()

	// Create a test category
	testCategory := &domain.Category{
		ID:          1,
		Name:        "Test Category",
		Description: stringPtr("Test description"),
		Color:       stringPtr("#FF0000"),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	mockCategoryService.categories[1] = testCategory

	tests := []struct {
	name           string
	categoryID     int64
	expectedStatus int
}{
		{
			name:           "delete existing category",
			categoryID:     1,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "delete non-existent category",
			categoryID:     999,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/v1/categories/1", nil)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

// Helper function
func stringPtr(s string) *string {
	return &s
}
