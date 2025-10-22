package repo

import (
	"database/sql"
	"fmt"
	"task-manager/internal/domain"
)

type CategoryRepository interface {
	Create(category *domain.Category) error
	GetByID(id int64) (*domain.Category, error)
	GetAll() ([]domain.Category, error)
	Update(category *domain.Category) error
	Delete(id int64) error
	GetByTaskID(taskID int64) ([]domain.Category, error)
	AddTaskCategory(taskID, categoryID int64) error
	RemoveTaskCategory(taskID, categoryID int64) error
	RemoveAllTaskCategories(taskID int64) error
}

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *domain.Category) error {
	query := `
		INSERT INTO categories (name, description, color, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, category.Name, category.Description, category.Color, category.CreatedAt, category.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get category ID: %w", err)
	}

	category.ID = id
	return nil
}

func (r *categoryRepository) GetByID(id int64) (*domain.Category, error) {
	query := `
		SELECT id, name, description, color, created_at, updated_at
		FROM categories
		WHERE id = ?
	`
	category := &domain.Category{}
	err := r.db.QueryRow(query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.Color,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("category not found")
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return category, nil
}

func (r *categoryRepository) GetAll() ([]domain.Category, error) {
	query := `
		SELECT id, name, description, color, created_at, updated_at
		FROM categories
		ORDER BY name ASC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var category domain.Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.Color,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating categories: %w", err)
	}

	return categories, nil
}

func (r *categoryRepository) Update(category *domain.Category) error {
	query := `
		UPDATE categories
		SET name = ?, description = ?, color = ?, updated_at = ?
		WHERE id = ?
	`
	result, err := r.db.Exec(query, category.Name, category.Description, category.Color, category.UpdatedAt, category.ID)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}

func (r *categoryRepository) Delete(id int64) error {
	query := `DELETE FROM categories WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}

func (r *categoryRepository) GetByTaskID(taskID int64) ([]domain.Category, error) {
	query := `
		SELECT c.id, c.name, c.description, c.color, c.created_at, c.updated_at
		FROM categories c
		INNER JOIN task_categories tc ON c.id = tc.category_id
		WHERE tc.task_id = ?
		ORDER BY c.name ASC
	`
	rows, err := r.db.Query(query, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories for task: %w", err)
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var category domain.Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.Color,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating categories: %w", err)
	}

	return categories, nil
}

func (r *categoryRepository) AddTaskCategory(taskID, categoryID int64) error {
	query := `INSERT INTO task_categories (task_id, category_id) VALUES (?, ?)`
	_, err := r.db.Exec(query, taskID, categoryID)
	if err != nil {
		return fmt.Errorf("failed to add task category: %w", err)
	}
	return nil
}

func (r *categoryRepository) RemoveTaskCategory(taskID, categoryID int64) error {
	query := `DELETE FROM task_categories WHERE task_id = ? AND category_id = ?`
	result, err := r.db.Exec(query, taskID, categoryID)
	if err != nil {
		return fmt.Errorf("failed to remove task category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("task category relationship not found")
	}

	return nil
}

func (r *categoryRepository) RemoveAllTaskCategories(taskID int64) error {
	query := `DELETE FROM task_categories WHERE task_id = ?`
	_, err := r.db.Exec(query, taskID)
	if err != nil {
		return fmt.Errorf("failed to remove all task categories: %w", err)
	}
	return nil
}
