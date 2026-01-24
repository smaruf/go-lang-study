package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/smaruf/go-lang-study/src/3-weeks-plan/pkg/database"
	"github.com/smaruf/go-lang-study/src/3-weeks-plan/pkg/models"
)

// Handler holds dependencies for HTTP handlers
type Handler struct {
	db *database.DB
}

// New creates a new Handler
func New(db *database.DB) *Handler {
	return &Handler{db: db}
}

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// sendJSON sends a JSON response
func sendJSON(w http.ResponseWriter, status int, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// CreateTask handles POST /api/tasks
func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	if req.Title == "" {
		sendJSON(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Title is required",
		})
		return
	}

	task, err := h.db.CreateTask(&req)
	if err != nil {
		sendJSON(w, http.StatusInternalServerError, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	sendJSON(w, http.StatusCreated, Response{
		Success: true,
		Data:    task,
	})
}

// GetTask handles GET /api/tasks/{id}
func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendJSON(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid task ID",
		})
		return
	}

	task, err := h.db.GetTask(id)
	if err != nil {
		sendJSON(w, http.StatusInternalServerError, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if task == nil {
		sendJSON(w, http.StatusNotFound, Response{
			Success: false,
			Error:   "Task not found",
		})
		return
	}

	sendJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    task,
	})
}

// ListTasks handles GET /api/tasks
func (h *Handler) ListTasks(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	priority := r.URL.Query().Get("priority")

	tasks, err := h.db.ListTasks(status, priority)
	if err != nil {
		sendJSON(w, http.StatusInternalServerError, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	sendJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    tasks,
	})
}

// UpdateTask handles PUT /api/tasks/{id}
func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendJSON(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid task ID",
		})
		return
	}

	var req models.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	task, err := h.db.UpdateTask(id, &req)
	if err != nil {
		sendJSON(w, http.StatusInternalServerError, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if task == nil {
		sendJSON(w, http.StatusNotFound, Response{
			Success: false,
			Error:   "Task not found",
		})
		return
	}

	sendJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    task,
	})
}

// DeleteTask handles DELETE /api/tasks/{id}
func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendJSON(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid task ID",
		})
		return
	}

	if err := h.db.DeleteTask(id); err != nil {
		sendJSON(w, http.StatusInternalServerError, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	sendJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    map[string]string{"message": "Task deleted successfully"},
	})
}

// GetStats handles GET /api/stats
func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.db.GetStats()
	if err != nil {
		sendJSON(w, http.StatusInternalServerError, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	sendJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    stats,
	})
}

// HealthCheck handles GET /health
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	sendJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    map[string]string{"status": "healthy"},
	})
}
