package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"task-manager/internal/domain"
	"task-manager/internal/service"

	"github.com/gorilla/mux"
)

type Handler struct {
	taskService     service.TaskService
	categoryService service.CategoryService
}

func NewHandler(taskService service.TaskService, categoryService service.CategoryService) *Handler {
	return &Handler{
		taskService:     taskService,
		categoryService: categoryService,
	}
}

func (h *Handler) SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.Use(corsMiddleware)
	r.Use(loggingMiddleware)

	api := r.PathPrefix("/v1").Subrouter()

	api.HandleFunc("/tasks", h.createTask).Methods("POST")
	api.HandleFunc("/tasks", h.getAllTasks).Methods("GET")
	api.HandleFunc("/tasks/{id}", h.getTask).Methods("GET")
	api.HandleFunc("/tasks/{id}", h.updateTask).Methods("PATCH")
	api.HandleFunc("/tasks/{id}", h.deleteTask).Methods("DELETE")

	// Category endpoints
	api.HandleFunc("/categories", h.createCategory).Methods("POST")
	api.HandleFunc("/categories", h.getAllCategories).Methods("GET")
	api.HandleFunc("/categories/{id}", h.getCategory).Methods("GET")
	api.HandleFunc("/categories/{id}", h.updateCategory).Methods("PUT")
	api.HandleFunc("/categories/{id}", h.deleteCategory).Methods("DELETE")

	return r
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	task, err := h.taskService.CreateTask(&req)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusCreated, task)
}

func (h *Handler) getTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	task, err := h.taskService.GetTask(id)
	if err != nil {
		writeErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, task)
}

func (h *Handler) getAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.taskService.GetAllTasks()
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, tasks)
}

func (h *Handler) updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var req domain.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	task, err := h.taskService.UpdateTask(id, &req)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, task)
}

func (h *Handler) deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	if err := h.taskService.DeleteTask(id); err != nil {
		writeErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Category handlers
func (h *Handler) createCategory(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	category, err := h.categoryService.CreateCategory(&req)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusCreated, category)
}

func (h *Handler) getCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	category, err := h.categoryService.GetCategory(id)
	if err != nil {
		writeErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, category)
}

func (h *Handler) getAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.categoryService.GetAllCategories()
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, categories)
}

func (h *Handler) updateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	var req domain.UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	category, err := h.categoryService.UpdateCategory(id, &req)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, category)
}

func (h *Handler) deleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	if err := h.categoryService.DeleteCategory(id); err != nil {
		writeErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
