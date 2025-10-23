package service

import (
	"errors"
	"strings"
	"task-manager/internal/domain"
	"testing"
	"time"
)

type mockTaskRepository struct {
	tasks  map[int64]*domain.Task
	nextID int64
}

func newMockTaskRepository() *mockTaskRepository {
	return &mockTaskRepository{
		tasks:  make(map[int64]*domain.Task),
		nextID: 1,
	}
}

type mockCategoryRepository struct {
	categories map[int64]*domain.Category
	nextID     int64
}

func newMockCategoryRepository() *mockCategoryRepository {
	return &mockCategoryRepository{
		categories: make(map[int64]*domain.Category),
		nextID:     1,
	}
}

func (m *mockCategoryRepository) Create(category *domain.Category) error {
	category.ID = m.nextID
	m.categories[m.nextID] = category
	m.nextID++
	return nil
}

func (m *mockCategoryRepository) GetByID(id int64) (*domain.Category, error) {
	category, exists := m.categories[id]
	if !exists {
		return nil, errors.New("category not found")
	}
	return category, nil
}

func (m *mockCategoryRepository) GetAll() ([]domain.Category, error) {
	var categories []domain.Category
	for _, category := range m.categories {
		categories = append(categories, *category)
	}
	return categories, nil
}

func (m *mockCategoryRepository) Update(category *domain.Category) error {
	_, exists := m.categories[category.ID]
	if !exists {
		return errors.New("category not found")
	}
	m.categories[category.ID] = category
	return nil
}

func (m *mockCategoryRepository) Delete(id int64) error {
	_, exists := m.categories[id]
	if !exists {
		return errors.New("category not found")
	}
	delete(m.categories, id)
	return nil
}

func (m *mockCategoryRepository) GetByTaskID(taskID int64) ([]domain.Category, error) {
	// Simple implementation for testing
	return []domain.Category{}, nil
}

func (m *mockCategoryRepository) AddTaskCategory(taskID, categoryID int64) error {
	// Simple implementation for testing
	return nil
}

func (m *mockCategoryRepository) RemoveTaskCategory(taskID, categoryID int64) error {
	// Simple implementation for testing
	return nil
}

func (m *mockCategoryRepository) RemoveAllTaskCategories(taskID int64) error {
	// Simple implementation for testing
	return nil
}

func (m *mockTaskRepository) Create(task *domain.Task) error {
	task.ID = m.nextID
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	m.tasks[m.nextID] = task
	m.nextID++
	return nil
}

func (m *mockTaskRepository) GetByID(id int64) (*domain.Task, error) {
	task, exists := m.tasks[id]
	if !exists {
		return nil, errors.New("task not found")
	}
	return task, nil
}

func (m *mockTaskRepository) GetAll() ([]*domain.Task, error) {
	var tasks []*domain.Task
	for _, task := range m.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (m *mockTaskRepository) GetWithFilters(filters *domain.TaskFilters) ([]*domain.Task, error) {
	if filters == nil {
		return m.GetAll()
	}

	var filteredTasks []*domain.Task
	for _, task := range m.tasks {
		// Check status filter
		if len(filters.Statuses) > 0 {
			statusMatch := false
			for _, status := range filters.Statuses {
				if task.Status == status {
					statusMatch = true
					break
				}
			}
			if !statusMatch {
				continue
			}
		}

		// Check priority filter
		if len(filters.Priorities) > 0 {
			priorityMatch := false
			for _, priority := range filters.Priorities {
				if task.Priority == priority {
					priorityMatch = true
					break
				}
			}
			if !priorityMatch {
				continue
			}
		}

		// Check search filter
		if filters.Search != "" {
			searchMatch := false
			searchLower := strings.ToLower(filters.Search)
			titleLower := strings.ToLower(task.Title)
			descriptionLower := strings.ToLower(task.Description)

			if strings.Contains(titleLower, searchLower) || strings.Contains(descriptionLower, searchLower) {
				searchMatch = true
			}

			if !searchMatch {
				continue
			}
		}

		filteredTasks = append(filteredTasks, task)
	}

	return filteredTasks, nil
}

func (m *mockTaskRepository) Update(task *domain.Task) error {
	if _, exists := m.tasks[task.ID]; !exists {
		return errors.New("task not found")
	}
	task.UpdatedAt = time.Now()
	m.tasks[task.ID] = task
	return nil
}

func (m *mockTaskRepository) Delete(id int64) error {
	if _, exists := m.tasks[id]; !exists {
		return errors.New("task not found")
	}
	delete(m.tasks, id)
	return nil
}

func TestTaskService_CreateTask(t *testing.T) {
	mockRepo := newMockTaskRepository()
	mockCategoryRepo := newMockCategoryRepository()
	service := NewTaskService(mockRepo, mockCategoryRepo)

	tests := []struct {
		name    string
		req     *domain.CreateTaskRequest
		wantErr bool
	}{
		{
			name: "valid task creation",
			req: &domain.CreateTaskRequest{
				Title:       "Test Task",
				Description: "Test Description",
				Priority:    domain.PriorityHigh,
			},
			wantErr: false,
		},
		{
			name: "invalid request",
			req: &domain.CreateTaskRequest{
				Title:       "",
				Description: "Test Description",
				Priority:    domain.PriorityHigh,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := service.CreateTask(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if task == nil {
					t.Error("CreateTask() returned nil task")
					return
				}
				if task.Title != tt.req.Title {
					t.Errorf("CreateTask() title = %v, want %v", task.Title, tt.req.Title)
				}
				if task.Status != domain.StatusTodo {
					t.Errorf("CreateTask() status = %v, want %v", task.Status, domain.StatusTodo)
				}
			}
		})
	}
}

func TestTaskService_GetTask(t *testing.T) {
	mockRepo := newMockTaskRepository()
	mockCategoryRepo := newMockCategoryRepository()
	service := NewTaskService(mockRepo, mockCategoryRepo)

	// Create a test task
	req := &domain.CreateTaskRequest{
		Title:       "Test Task",
		Description: "Test Description",
		Priority:    domain.PriorityMedium,
	}
	createdTask, _ := service.CreateTask(req)

	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{
			name:    "valid task id",
			id:      createdTask.ID,
			wantErr: false,
		},
		{
			name:    "invalid task id",
			id:      999,
			wantErr: true,
		},
		{
			name:    "zero task id",
			id:      0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := service.GetTask(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && task == nil {
				t.Error("GetTask() returned nil task")
			}
		})
	}
}

func TestTaskService_UpdateTask(t *testing.T) {
	mockRepo := newMockTaskRepository()
	mockCategoryRepo := newMockCategoryRepository()
	service := NewTaskService(mockRepo, mockCategoryRepo)

	// Create a test task
	req := &domain.CreateTaskRequest{
		Title:       "Test Task",
		Description: "Test Description",
		Priority:    domain.PriorityMedium,
	}
	createdTask, _ := service.CreateTask(req)

	tests := []struct {
		name    string
		id      int64
		req     *domain.UpdateTaskRequest
		wantErr bool
	}{
		{
			name: "valid update",
			id:   createdTask.ID,
			req: &domain.UpdateTaskRequest{
				Title:  stringPtr("Updated Task"),
				Status: func() *domain.TaskStatus { s := domain.StatusDone; return &s }(),
			},
			wantErr: false,
		},
		{
			name: "invalid task id",
			id:   999,
			req: &domain.UpdateTaskRequest{
				Title: stringPtr("Updated Task"),
			},
			wantErr: true,
		},
		{
			name: "invalid update data",
			id:   createdTask.ID,
			req: &domain.UpdateTaskRequest{
				Title: stringPtr(""),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := service.UpdateTask(tt.id, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && task == nil {
				t.Error("UpdateTask() returned nil task")
			}
		})
	}
}

func TestTaskService_DeleteTask(t *testing.T) {
	mockRepo := newMockTaskRepository()
	mockCategoryRepo := newMockCategoryRepository()
	service := NewTaskService(mockRepo, mockCategoryRepo)

	// Create a test task
	req := &domain.CreateTaskRequest{
		Title:       "Test Task",
		Description: "Test Description",
		Priority:    domain.PriorityMedium,
	}
	createdTask, _ := service.CreateTask(req)

	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{
			name:    "valid delete",
			id:      createdTask.ID,
			wantErr: false,
		},
		{
			name:    "invalid task id",
			id:      999,
			wantErr: true,
		},
		{
			name:    "zero task id",
			id:      0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.DeleteTask(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
