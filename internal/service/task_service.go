package service

import (
	"errors"
	"fmt"
	"task-manager/internal/domain"
	"task-manager/internal/repo"
	"time"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskService interface {
	CreateTask(req *domain.CreateTaskRequest) (*domain.Task, error)
	GetTask(id int64) (*domain.Task, error)
	GetAllTasks() ([]*domain.Task, error)
	GetTasksWithFilters(filters *domain.TaskFilters) ([]*domain.Task, error)
	UpdateTask(id int64, req *domain.UpdateTaskRequest) (*domain.Task, error)
	DeleteTask(id int64) error
}

type taskService struct {
	taskRepo     repo.TaskRepository
	categoryRepo repo.CategoryRepository
}

func NewTaskService(taskRepo repo.TaskRepository, categoryRepo repo.CategoryRepository) TaskService {
	return &taskService{
		taskRepo:     taskRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *taskService) CreateTask(req *domain.CreateTaskRequest) (*domain.Task, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	now := time.Now()
	task := &domain.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      domain.StatusTodo,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.taskRepo.Create(task); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	// Add categories to the task
	for _, categoryID := range req.CategoryIDs {
		if err := s.categoryRepo.AddTaskCategory(task.ID, categoryID); err != nil {
			return nil, fmt.Errorf("failed to add category to task: %w", err)
		}
	}

	// Reload the task with categories
	return s.taskRepo.GetByID(task.ID)
}

func (s *taskService) GetTask(id int64) (*domain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid task id")
	}

	task, err := s.taskRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return task, nil
}

func (s *taskService) GetAllTasks() ([]*domain.Task, error) {
	tasks, err := s.taskRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	return tasks, nil
}

func (s *taskService) GetTasksWithFilters(filters *domain.TaskFilters) ([]*domain.Task, error) {
	if filters == nil {
		return s.GetAllTasks()
	}

	if err := filters.Validate(); err != nil {
		return nil, fmt.Errorf("invalid filters: %w", err)
	}

	tasks, err := s.taskRepo.GetWithFilters(filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks with filters: %w", err)
	}

	return tasks, nil
}

func (s *taskService) UpdateTask(id int64, req *domain.UpdateTaskRequest) (*domain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid task id")
	}

	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	existingTask, err := s.taskRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing task: %w", err)
	}

	if req.Title != nil {
		existingTask.Title = *req.Title
	}
	if req.Description != nil {
		existingTask.Description = *req.Description
	}
	if req.Status != nil {
		existingTask.Status = *req.Status
	}
	if req.Priority != nil {
		existingTask.Priority = *req.Priority
	}
	if req.DueDate != nil {
		existingTask.DueDate = req.DueDate
	}

	existingTask.UpdatedAt = time.Now()

	if err := existingTask.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if err := s.taskRepo.Update(existingTask); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	// Handle category updates
	if req.CategoryIDs != nil {
		// Remove all existing categories
		if err := s.categoryRepo.RemoveAllTaskCategories(id); err != nil {
			return nil, fmt.Errorf("failed to remove existing categories: %w", err)
		}

		// Add new categories
		for _, categoryID := range *req.CategoryIDs {
			if err := s.categoryRepo.AddTaskCategory(id, categoryID); err != nil {
				return nil, fmt.Errorf("failed to add category to task: %w", err)
			}
		}
	}

	// Reload the task with updated categories
	return s.taskRepo.GetByID(id)
}

func (s *taskService) DeleteTask(id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid task id")
	}

	if err := s.taskRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}
