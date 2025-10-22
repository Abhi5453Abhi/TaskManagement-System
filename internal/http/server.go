package http

import (
	"net/http"
	"task-manager/internal/service"

	"github.com/gorilla/mux"
)

type Server struct {
	taskService service.TaskService
	router      *mux.Router
}

func NewServer(taskService service.TaskService) *Server {
	handler := NewHandler(taskService)
	router := handler.SetupRoutes()

	return &Server{
		taskService: taskService,
		router:      router,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
