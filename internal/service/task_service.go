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
	UpdateTask(id int64, req *domain.UpdateTaskRequest) (*domain.Task, error)
	DeleteTask(id int64) error
}

type taskService struct {
	taskRepo repo.TaskRepository
}

func NewTaskService(taskRepo repo.TaskRepository) TaskService {
	return &taskService{
		taskRepo: taskRepo,
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

	return task, nil
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

	return existingTask, nil
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
