package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/joho/godotenv"
)

// Task represents a single task
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Priority    string    `json:"priority"` // low, medium, high
	Tags        []string  `json:"tags"`
}

// TaskManager manages tasks and persistence
type TaskManager struct {
	tasks    []Task
	dataFile string
	nextID   int
}

// NewTaskManager creates a new task manager
func NewTaskManager(dataFile string) *TaskManager {
	tm := &TaskManager{
		tasks:    make([]Task, 0),
		dataFile: dataFile,
		nextID:   1,
	}
	tm.loadTasks()
	return tm
}

// loadTasks loads tasks from the data file
func (tm *TaskManager) loadTasks() error {
	if _, err := os.Stat(tm.dataFile); os.IsNotExist(err) {
		return nil // File doesn't exist, start with empty tasks
	}

	data, err := ioutil.ReadFile(tm.dataFile)
	if err != nil {
		return fmt.Errorf("failed to read tasks file: %v", err)
	}

	if len(data) == 0 {
		return nil // Empty file
	}

	if err := json.Unmarshal(data, &tm.tasks); err != nil {
		return fmt.Errorf("failed to parse tasks file: %v", err)
	}

	// Find the highest ID to set nextID
	for _, task := range tm.tasks {
		if task.ID >= tm.nextID {
			tm.nextID = task.ID + 1
		}
	}

	return nil
}

// saveTasks saves tasks to the data file
func (tm *TaskManager) saveTasks() error {
	// Ensure directory exists
	dir := filepath.Dir(tm.dataFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	data, err := json.MarshalIndent(tm.tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %v", err)
	}

	if err := ioutil.WriteFile(tm.dataFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write tasks file: %v", err)
	}

	return nil
}

// addTask adds a new task
func (tm *TaskManager) addTask(title, description, priority string, tags []string) (*Task, error) {
	task := Task{
		ID:          tm.nextID,
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Priority:    priority,
		Tags:        tags,
	}

	tm.tasks = append(tm.tasks, task)
	tm.nextID++

	if err := tm.saveTasks(); err != nil {
		return nil, err
	}

	return &task, nil
}

// listTasks returns filtered tasks
func (tm *TaskManager) listTasks(showCompleted bool, priority string, tag string) []Task {
	var filtered []Task
	
	for _, task := range tm.tasks {
		// Filter by completion status
		if !showCompleted && task.Completed {
			continue
		}
		
		// Filter by priority
		if priority != "" && task.Priority != priority {
			continue
		}
		
		// Filter by tag
		if tag != "" {
			hasTag := false
			for _, t := range task.Tags {
				if t == tag {
					hasTag = true
					break
				}
			}
			if !hasTag {
				continue
			}
		}
		
		filtered = append(filtered, task)
	}

	return filtered
}

// getTask finds a task by ID
func (tm *TaskManager) getTask(id int) (*Task, error) {
	for i, task := range tm.tasks {
		if task.ID == id {
			return &tm.tasks[i], nil
		}
	}
	return nil, fmt.Errorf("task with ID %d not found", id)
}

// updateTask updates an existing task
func (tm *TaskManager) updateTask(id int, title, description, priority string, tags []string) error {
	task, err := tm.getTask(id)
	if err != nil {
		return err
	}

	if title != "" {
		task.Title = title
	}
	if description != "" {
		task.Description = description
	}
	if priority != "" {
		task.Priority = priority
	}
	if len(tags) > 0 {
		task.Tags = tags
	}
	task.UpdatedAt = time.Now()

	return tm.saveTasks()
}

// completeTask marks a task as completed
func (tm *TaskManager) completeTask(id int) error {
	task, err := tm.getTask(id)
	if err != nil {
		return err
	}

	task.Completed = true
	task.UpdatedAt = time.Now()

	return tm.saveTasks()
}

// deleteTask removes a task
func (tm *TaskManager) deleteTask(id int) error {
	for i, task := range tm.tasks {
		if task.ID == id {
			tm.tasks = append(tm.tasks[:i], tm.tasks[i+1:]...)
			return tm.saveTasks()
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

// Global task manager instance
var taskManager *TaskManager

// Command line flags
var (
	dataFile      string
	showCompleted bool
	priority      string
	tag           string
	description   string
	tags          []string
)

func main() {
	// Load environment variables
	godotenv.Load()

	// Initialize task manager
	defaultDataFile := getEnv("TASK_DATA_FILE", filepath.Join(os.Getenv("HOME"), ".tasks", "tasks.json"))
	taskManager = NewTaskManager(defaultDataFile)

	// Execute root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "task",
	Short: "A modern CLI task management tool",
	Long: `Task Manager is a CLI application for managing your tasks.
It supports creating, listing, updating, and deleting tasks with
priorities, tags, and completion tracking.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks",
	Long:  "List all tasks with optional filtering by completion status, priority, or tags.",
	Run: func(cmd *cobra.Command, args []string) {
		tasks := taskManager.listTasks(showCompleted, priority, tag)
		
		if len(tasks) == 0 {
			fmt.Println("No tasks found.")
			return
		}

		fmt.Printf("Found %d task(s):\n\n", len(tasks))
		
		for _, task := range tasks {
			status := "❌"
			if task.Completed {
				status = "✅"
			}

			priorityColor := getPriorityColor(task.Priority)
			
			fmt.Printf("%s [ID:%d] %s%s%s\n", status, task.ID, priorityColor, task.Title, resetColor())
			
			if task.Description != "" {
				fmt.Printf("   Description: %s\n", task.Description)
			}
			
			if len(task.Tags) > 0 {
				fmt.Printf("   Tags: %s\n", strings.Join(task.Tags, ", "))
			}
			
			fmt.Printf("   Priority: %s | Created: %s\n", 
				task.Priority, task.CreatedAt.Format("2006-01-02 15:04"))
			fmt.Println()
		}
	},
}

var addCmd = &cobra.Command{
	Use:   "add [title]",
	Short: "Add a new task",
	Long:  "Add a new task with optional description, priority, and tags.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		title := strings.Join(args, " ")
		
		if priority == "" {
			priority = "medium"
		}
		
		task, err := taskManager.addTask(title, description, priority, tags)
		if err != nil {
			fmt.Printf("Error adding task: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Task added successfully!\n")
		fmt.Printf("ID: %d\n", task.ID)
		fmt.Printf("Title: %s\n", task.Title)
		if task.Description != "" {
			fmt.Printf("Description: %s\n", task.Description)
		}
		fmt.Printf("Priority: %s\n", task.Priority)
		if len(task.Tags) > 0 {
			fmt.Printf("Tags: %s\n", strings.Join(task.Tags, ", "))
		}
	},
}

var updateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update an existing task",
	Long:  "Update an existing task's title, description, priority, or tags.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Invalid task ID: %s\n", args[0])
			os.Exit(1)
		}

		// Get the title from the --title flag
		title, _ := cmd.Flags().GetString("title")
		
		err = taskManager.updateTask(id, title, description, priority, tags)
		if err != nil {
			fmt.Printf("Error updating task: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Task %d updated successfully!\n", id)
	},
}

var completeCmd = &cobra.Command{
	Use:   "complete [id]",
	Short: "Mark a task as completed",
	Long:  "Mark a task as completed by its ID.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Invalid task ID: %s\n", args[0])
			os.Exit(1)
		}

		err = taskManager.completeTask(id)
		if err != nil {
			fmt.Printf("Error completing task: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Task %d marked as completed!\n", id)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a task",
	Long:  "Delete a task by its ID.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Invalid task ID: %s\n", args[0])
			os.Exit(1)
		}

		err = taskManager.deleteTask(id)
		if err != nil {
			fmt.Printf("Error deleting task: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("🗑️  Task %d deleted successfully!\n", id)
	},
}

func init() {
	// Root command flags
	rootCmd.PersistentFlags().StringVar(&dataFile, "data-file", "", "Path to the tasks data file")

	// List command flags
	listCmd.Flags().BoolVar(&showCompleted, "completed", false, "Show completed tasks")
	listCmd.Flags().StringVar(&priority, "priority", "", "Filter by priority (low, medium, high)")
	listCmd.Flags().StringVar(&tag, "tag", "", "Filter by tag")

	// Add command flags
	addCmd.Flags().StringVar(&description, "description", "", "Task description")
	addCmd.Flags().StringVar(&priority, "priority", "medium", "Task priority (low, medium, high)")
	addCmd.Flags().StringSliceVar(&tags, "tags", []string{}, "Task tags (comma-separated)")

	// Update command flags
	updateCmd.Flags().String("title", "", "New task title")
	updateCmd.Flags().StringVar(&description, "description", "", "New task description")
	updateCmd.Flags().StringVar(&priority, "priority", "", "New task priority (low, medium, high)")
	updateCmd.Flags().StringSliceVar(&tags, "tags", []string{}, "New task tags (comma-separated)")

	// Add commands to root
	rootCmd.AddCommand(listCmd, addCmd, updateCmd, completeCmd, deleteCmd)
}

// Helper functions

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getPriorityColor(priority string) string {
	switch priority {
	case "high":
		return "\033[31m" // Red
	case "medium":
		return "\033[33m" // Yellow
	case "low":
		return "\033[32m" // Green
	default:
		return ""
	}
}

func resetColor() string {
	return "\033[0m"
}