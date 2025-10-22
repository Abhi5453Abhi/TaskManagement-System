package service

import (
	"errors"
	"time"
	"task-manager/internal/domain"
	"task-manager/internal/repo"
)

type CategoryService interface {
	CreateCategory(req *domain.CreateCategoryRequest) (*domain.Category, error)
	GetCategory(id int64) (*domain.Category, error)
	GetAllCategories() ([]domain.Category, error)
	UpdateCategory(id int64, req *domain.UpdateCategoryRequest) (*domain.Category, error)
	DeleteCategory(id int64) error
}

type categoryService struct {
	categoryRepo repo.CategoryRepository
}

func NewCategoryService(categoryRepo repo.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

var ErrCategoryNotFound = errors.New("category not found")

func (s *categoryService) CreateCategory(req *domain.CreateCategoryRequest) (*domain.Category, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	now := time.Now()
	category := &domain.Category{
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) GetCategory(id int64) (*domain.Category, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *categoryService) GetAllCategories() ([]domain.Category, error) {
	categories, err := s.categoryRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *categoryService) UpdateCategory(id int64, req *domain.UpdateCategoryRequest) (*domain.Category, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, ErrCategoryNotFound
	}

	// Update fields if provided
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

	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *categoryService) DeleteCategory(id int64) error {
	// Check if category exists
	_, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return ErrCategoryNotFound
	}

	return s.categoryRepo.Delete(id)
}
