package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/smaruf/go-lang-study/src/3-weeks-plan/pkg/database"
	"github.com/smaruf/go-lang-study/src/3-weeks-plan/pkg/handlers"
	"github.com/smaruf/go-lang-study/src/3-weeks-plan/pkg/models"
)

func setupTestDB(t *testing.T) *database.DB {
	dbPath := "/tmp/test_tasks.db"
	// Remove existing test database
	os.Remove(dbPath)

	db, err := database.New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	return db
}

func TestCreateTask(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	h := handlers.New(db)
	router := mux.NewRouter()
	router.HandleFunc("/api/tasks", h.CreateTask).Methods("POST")

	reqBody := models.CreateTaskRequest{
		Title:       "Test Task",
		Description: "This is a test task",
		Priority:    "high",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/tasks", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", rec.Code)
	}

	var response handlers.Response
	json.NewDecoder(rec.Body).Decode(&response)

	if !response.Success {
		t.Errorf("Expected success to be true, got false")
	}

	taskData, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Response data is not a task")
	}

	if taskData["title"] != "Test Task" {
		t.Errorf("Expected title 'Test Task', got '%v'", taskData["title"])
	}

	if taskData["priority"] != "high" {
		t.Errorf("Expected priority 'high', got '%v'", taskData["priority"])
	}
}

func TestGetTask(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create a task first
	task, err := db.CreateTask(&models.CreateTaskRequest{
		Title:       "Test Task",
		Description: "Test Description",
		Priority:    "medium",
	})
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	h := handlers.New(db)
	router := mux.NewRouter()
	router.HandleFunc("/api/tasks/{id}", h.GetTask).Methods("GET")

	req := httptest.NewRequest("GET", "/api/tasks/1", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var response handlers.Response
	json.NewDecoder(rec.Body).Decode(&response)

	if !response.Success {
		t.Errorf("Expected success to be true, got false")
	}

	taskData, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Response data is not a task")
	}

	if int(taskData["id"].(float64)) != task.ID {
		t.Errorf("Expected task ID %d, got %v", task.ID, taskData["id"])
	}
}

func TestListTasks(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create multiple tasks
	tasks := []models.CreateTaskRequest{
		{Title: "Task 1", Priority: "high"},
		{Title: "Task 2", Priority: "medium"},
		{Title: "Task 3", Priority: "low"},
	}

	for _, task := range tasks {
		_, err := db.CreateTask(&task)
		if err != nil {
			t.Fatalf("Failed to create task: %v", err)
		}
	}

	h := handlers.New(db)
	router := mux.NewRouter()
	router.HandleFunc("/api/tasks", h.ListTasks).Methods("GET")

	req := httptest.NewRequest("GET", "/api/tasks", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var response handlers.Response
	json.NewDecoder(rec.Body).Decode(&response)

	if !response.Success {
		t.Errorf("Expected success to be true, got false")
	}

	taskList, ok := response.Data.([]interface{})
	if !ok {
		t.Fatal("Response data is not a list")
	}

	if len(taskList) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(taskList))
	}
}

func TestUpdateTask(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create a task first
	task, err := db.CreateTask(&models.CreateTaskRequest{
		Title:    "Original Title",
		Priority: "medium",
	})
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	h := handlers.New(db)
	router := mux.NewRouter()
	router.HandleFunc("/api/tasks/{id}", h.UpdateTask).Methods("PUT")

	newTitle := "Updated Title"
	newStatus := "completed"
	updateReq := models.UpdateTaskRequest{
		Title:  &newTitle,
		Status: &newStatus,
	}

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/api/tasks/1", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var response handlers.Response
	json.NewDecoder(rec.Body).Decode(&response)

	if !response.Success {
		t.Errorf("Expected success to be true, got false")
	}

	taskData, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Response data is not a task")
	}

	if taskData["title"] != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got '%v'", taskData["title"])
	}

	if taskData["status"] != "completed" {
		t.Errorf("Expected status 'completed', got '%v'", taskData["status"])
	}

	if int(taskData["id"].(float64)) != task.ID {
		t.Errorf("Expected task ID %d, got %v", task.ID, taskData["id"])
	}
}

func TestDeleteTask(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create a task first
	_, err := db.CreateTask(&models.CreateTaskRequest{
		Title:    "Task to Delete",
		Priority: "low",
	})
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	h := handlers.New(db)
	router := mux.NewRouter()
	router.HandleFunc("/api/tasks/{id}", h.DeleteTask).Methods("DELETE")

	req := httptest.NewRequest("DELETE", "/api/tasks/1", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var response handlers.Response
	json.NewDecoder(rec.Body).Decode(&response)

	if !response.Success {
		t.Errorf("Expected success to be true, got false")
	}

	// Verify task is deleted
	deletedTask, err := db.GetTask(1)
	if err != nil {
		t.Fatalf("Failed to check deleted task: %v", err)
	}
	if deletedTask != nil {
		t.Error("Expected task to be deleted, but it still exists")
	}
}

func TestGetStats(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create tasks with different statuses and priorities
	tasks := []models.CreateTaskRequest{
		{Title: "Task 1", Priority: "high"},
		{Title: "Task 2", Priority: "high"},
		{Title: "Task 3", Priority: "medium"},
	}

	for _, task := range tasks {
		_, err := db.CreateTask(&task)
		if err != nil {
			t.Fatalf("Failed to create task: %v", err)
		}
	}

	// Update one task to completed
	status := "completed"
	_, err := db.UpdateTask(1, &models.UpdateTaskRequest{Status: &status})
	if err != nil {
		t.Fatalf("Failed to update task: %v", err)
	}

	h := handlers.New(db)
	router := mux.NewRouter()
	router.HandleFunc("/api/stats", h.GetStats).Methods("GET")

	req := httptest.NewRequest("GET", "/api/stats", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var response handlers.Response
	json.NewDecoder(rec.Body).Decode(&response)

	if !response.Success {
		t.Errorf("Expected success to be true, got false")
	}

	statsData, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Response data is not stats")
	}

	total := int(statsData["total"].(float64))
	if total != 3 {
		t.Errorf("Expected total 3, got %d", total)
	}
}
