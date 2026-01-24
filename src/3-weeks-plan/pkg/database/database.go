package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/smaruf/go-lang-study/src/3-weeks-plan/pkg/models"
)

// DB wraps the database connection
type DB struct {
	conn *sql.DB
}

// New creates a new database connection
func New(dbPath string) (*DB, error) {
	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db := &DB{conn: conn}
	if err := db.migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

// migrate creates the necessary tables
func (db *DB) migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		status TEXT NOT NULL DEFAULT 'pending',
		priority TEXT NOT NULL DEFAULT 'medium',
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		due_date DATETIME
	);

	CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
	CREATE INDEX IF NOT EXISTS idx_tasks_priority ON tasks(priority);
	CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at);
	`

	_, err := db.conn.Exec(query)
	return err
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// CreateTask creates a new task
func (db *DB) CreateTask(req *models.CreateTaskRequest) (*models.Task, error) {
	now := time.Now()
	priority := req.Priority
	if priority == "" {
		priority = "medium"
	}

	query := `
	INSERT INTO tasks (title, description, status, priority, created_at, updated_at, due_date)
	VALUES (?, ?, 'pending', ?, ?, ?, ?)
	`

	result, err := db.conn.Exec(query, req.Title, req.Description, priority, now, now, req.DueDate)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return db.GetTask(int(id))
}

// GetTask retrieves a task by ID
func (db *DB) GetTask(id int) (*models.Task, error) {
	query := `
	SELECT id, title, description, status, priority, created_at, updated_at, due_date
	FROM tasks
	WHERE id = ?
	`

	task := &models.Task{}
	var dueDate sql.NullTime

	err := db.conn.QueryRow(query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.CreatedAt,
		&task.UpdatedAt,
		&dueDate,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	if dueDate.Valid {
		task.DueDate = &dueDate.Time
	}

	return task, nil
}

// ListTasks retrieves all tasks with optional filtering
func (db *DB) ListTasks(status, priority string) ([]*models.Task, error) {
	query := `
	SELECT id, title, description, status, priority, created_at, updated_at, due_date
	FROM tasks
	WHERE 1=1
	`
	args := []interface{}{}

	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}
	if priority != "" {
		query += " AND priority = ?"
		args = append(args, priority)
	}

	query += " ORDER BY created_at DESC"

	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}
	defer rows.Close()

	tasks := []*models.Task{}
	for rows.Next() {
		task := &models.Task{}
		var dueDate sql.NullTime

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.CreatedAt,
			&task.UpdatedAt,
			&dueDate,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		if dueDate.Valid {
			task.DueDate = &dueDate.Time
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// UpdateTask updates an existing task
func (db *DB) UpdateTask(id int, req *models.UpdateTaskRequest) (*models.Task, error) {
	// Build dynamic update query
	query := "UPDATE tasks SET updated_at = ?"
	args := []interface{}{time.Now()}

	if req.Title != nil {
		query += ", title = ?"
		args = append(args, *req.Title)
	}
	if req.Description != nil {
		query += ", description = ?"
		args = append(args, *req.Description)
	}
	if req.Status != nil {
		query += ", status = ?"
		args = append(args, *req.Status)
	}
	if req.Priority != nil {
		query += ", priority = ?"
		args = append(args, *req.Priority)
	}
	if req.DueDate != nil {
		query += ", due_date = ?"
		args = append(args, *req.DueDate)
	}

	query += " WHERE id = ?"
	args = append(args, id)

	_, err := db.conn.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return db.GetTask(id)
}

// DeleteTask deletes a task by ID
func (db *DB) DeleteTask(id int) error {
	query := "DELETE FROM tasks WHERE id = ?"
	_, err := db.conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}

// GetStats retrieves task statistics
func (db *DB) GetStats() (*models.TaskStats, error) {
	stats := &models.TaskStats{
		ByStatus:   make(map[string]int),
		ByPriority: make(map[string]int),
	}

	// Get total count
	err := db.conn.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&stats.Total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	// Get counts by status
	rows, err := db.conn.Query("SELECT status, COUNT(*) FROM tasks GROUP BY status")
	if err != nil {
		return nil, fmt.Errorf("failed to get status counts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		stats.ByStatus[status] = count
	}

	// Get counts by priority
	rows, err = db.conn.Query("SELECT priority, COUNT(*) FROM tasks GROUP BY priority")
	if err != nil {
		return nil, fmt.Errorf("failed to get priority counts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var priority string
		var count int
		if err := rows.Scan(&priority, &count); err != nil {
			return nil, err
		}
		stats.ByPriority[priority] = count
	}

	return stats, nil
}
