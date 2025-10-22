package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"task-manager/internal/config"
	httpHandler "task-manager/internal/http"
	"task-manager/internal/repo"
	"task-manager/internal/service"

	_ "modernc.org/sqlite"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("sqlite", cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if err := repo.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	taskRepo := repo.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo)

	httpServer := httpHandler.NewServer(taskService)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: httpServer,
	}

	go func() {
		log.Printf("Server starting on port %d", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited")
}
