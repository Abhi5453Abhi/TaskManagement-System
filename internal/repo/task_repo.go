package repo

import (
	"database/sql"
	"fmt"
	"task-manager/internal/domain"

	_ "modernc.org/sqlite"
)

type TaskRepository interface {
	Create(task *domain.Task) error
	GetByID(id int64) (*domain.Task, error)
	GetAll() ([]*domain.Task, error)
	Update(task *domain.Task) error
	Delete(id int64) error
}

type taskRepository struct {
	db               *sql.DB
	categoryRepo     CategoryRepository
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepository{
		db:           db,
		categoryRepo: NewCategoryRepository(db),
	}
}

func Migrate(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		status TEXT NOT NULL DEFAULT 'todo',
		priority TEXT NOT NULL DEFAULT 'medium',
		due_date DATE,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(100) NOT NULL UNIQUE,
		description TEXT,
		color VARCHAR(7),
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE TABLE IF NOT EXISTS task_categories (
		task_id INTEGER NOT NULL,
		category_id INTEGER NOT NULL,
		PRIMARY KEY (task_id, category_id),
		FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
		FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
	);
	
	CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
	CREATE INDEX IF NOT EXISTS idx_tasks_priority ON tasks(priority);
	CREATE INDEX IF NOT EXISTS idx_tasks_due_date ON tasks(due_date);
	CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at);
	CREATE INDEX IF NOT EXISTS idx_categories_name ON categories(name);
	CREATE INDEX IF NOT EXISTS idx_task_categories_task_id ON task_categories(task_id);
	CREATE INDEX IF NOT EXISTS idx_task_categories_category_id ON task_categories(category_id);
	`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	// Add due_date column if it doesn't exist (for existing databases)
	alterQuery := `ALTER TABLE tasks ADD COLUMN due_date DATE;`
	db.Exec(alterQuery) // Ignore error if column already exists

	return nil
}

func (r *taskRepository) Create(task *domain.Task) error {
	query := `
		INSERT INTO tasks (title, description, status, priority, due_date, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query, task.Title, task.Description, task.Status, task.Priority, task.DueDate, task.CreatedAt, task.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	task.ID = id
	return nil
}

func (r *taskRepository) GetByID(id int64) (*domain.Task, error) {
	query := `
		SELECT id, title, description, status, priority, due_date, created_at, updated_at
		FROM tasks
		WHERE id = ?
	`

	task := &domain.Task{}
	err := r.db.QueryRow(query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	// Load categories for the task
	categories, err := r.categoryRepo.GetByTaskID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task categories: %w", err)
	}
	task.Categories = categories

	return task, nil
}

func (r *taskRepository) GetAll() ([]*domain.Task, error) {
	query := `
		SELECT id, title, description, status, priority, due_date, created_at, updated_at
		FROM tasks
		ORDER BY 
			CASE WHEN due_date IS NULL THEN 1 ELSE 0 END,
			due_date ASC,
			CASE priority 
				WHEN 'critical' THEN 4
				WHEN 'high' THEN 3
				WHEN 'medium' THEN 2
				WHEN 'low' THEN 1
				ELSE 0
			END DESC,
			created_at ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*domain.Task
	for rows.Next() {
		task := &domain.Task{}
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.DueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		
		// Load categories for the task
		categories, err := r.categoryRepo.GetByTaskID(task.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get task categories: %w", err)
		}
		task.Categories = categories
		
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return tasks, nil
}

func (r *taskRepository) Update(task *domain.Task) error {
	query := `
		UPDATE tasks 
		SET title = ?, description = ?, status = ?, priority = ?, due_date = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.Exec(query, task.Title, task.Description, task.Status, task.Priority, task.DueDate, task.UpdatedAt, task.ID)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func (r *taskRepository) Delete(id int64) error {
	query := `DELETE FROM tasks WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}
