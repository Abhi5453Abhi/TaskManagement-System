package main

import (
	"database/sql"
	"fmt"
	"log"
	"task-manager/internal/domain"
	"task-manager/internal/repo"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "tasks.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if err := repo.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	taskRepo := repo.NewTaskRepository(db)

	// Seed some sample tasks
	sampleTasks := []domain.CreateTaskRequest{
		{
			Title:       "Complete project documentation",
			Description: "Write comprehensive documentation for the task manager API",
			Priority:    domain.PriorityHigh,
		},
		{
			Title:       "Implement user authentication",
			Description: "Add JWT-based authentication system",
			Priority:    domain.PriorityCritical,
		},
		{
			Title:       "Add task categories",
			Description: "Allow users to organize tasks by categories",
			Priority:    domain.PriorityMedium,
		},
		{
			Title:       "Optimize database queries",
			Description: "Review and optimize slow database queries",
			Priority:    domain.PriorityLow,
		},
		{
			Title:       "Add task due dates",
			Description: "Implement due date functionality for tasks",
			Priority:    domain.PriorityHigh,
		},
	}

	for i, taskReq := range sampleTasks {
		task := &domain.Task{
			Title:       taskReq.Title,
			Description: taskReq.Description,
			Status:      domain.StatusTodo,
			Priority:    taskReq.Priority,
		}

		if err := taskRepo.Create(task); err != nil {
			log.Printf("Failed to create task %d: %v", i+1, err)
		} else {
			fmt.Printf("Created task: %s (Priority: %s)\n", task.Title, task.Priority)
		}
	}

	fmt.Println("Database seeded successfully!")
}
