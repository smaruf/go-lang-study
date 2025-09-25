# CLI Task Manager

A modern command-line task management tool built with Go and Cobra framework, demonstrating best practices for CLI application development.

## Features

- **Task Management**: Create, list, update, complete, and delete tasks
- **Filtering**: Filter tasks by completion status, priority, or tags
- **Persistence**: JSON file-based storage with configurable location
- **Rich CLI Interface**: Color-coded output and intuitive commands
- **Environment Configuration**: Environment variable support
- **Data Validation**: Input validation and error handling
- **Cross-platform**: Works on Windows, macOS, and Linux

## Installation

1. **Build from source:**
   ```bash
   go build -o task main.go
   ```

2. **Or run directly:**
   ```bash
   go run main.go [command]
   ```

## Usage

### Basic Commands

```bash
# Show help
./task help

# Add a new task
./task add "Complete the project documentation"

# Add a task with description, priority, and tags
./task add "Fix critical bug" \
  --description "Fix the memory leak in user service" \
  --priority high \
  --tags "bug,critical,backend"

# List all active tasks
./task list

# List all tasks including completed ones
./task list --completed

# Filter tasks by priority
./task list --priority high

# Filter tasks by tag
./task list --tag "bug"

# Mark a task as completed
./task complete 1

# Update a task
./task update 1 \
  --title "Updated task title" \
  --description "New description" \
  --priority medium \
  --tags "updated,modified"

# Delete a task
./task delete 1
```

### Advanced Usage

```bash
# Use custom data file location
./task --data-file /path/to/tasks.json list

# Set data file via environment variable
export TASK_DATA_FILE="/home/user/.config/tasks.json"
./task list

# Chain commands for workflow automation
./task add "Review code" --priority high --tags "review,code"
./task add "Write tests" --priority medium --tags "testing"
./task add "Deploy to staging" --priority low --tags "deployment"
```

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `TASK_DATA_FILE` | `~/.tasks/tasks.json` | Path to the tasks data file |

### Data Storage

Tasks are stored in JSON format with the following structure:

```json
[
  {
    "id": 1,
    "title": "Complete project documentation",
    "description": "Write comprehensive docs for the API",
    "completed": false,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "priority": "high",
    "tags": ["documentation", "api"]
  }
]
```

## Command Reference

### Global Flags

- `--data-file`: Specify custom data file location

### `add` - Add New Task

Add a new task with optional properties.

```bash
task add "Task title" [flags]
```

**Flags:**
- `--description`: Task description
- `--priority`: Task priority (low, medium, high) [default: medium]
- `--tags`: Comma-separated list of tags

**Examples:**
```bash
task add "Buy groceries"
task add "Fix bug #123" --description "Critical memory leak" --priority high --tags "bug,critical"
```

### `list` - List Tasks

Display tasks with optional filtering.

```bash
task list [flags]
```

**Flags:**
- `--completed`: Include completed tasks
- `--priority`: Filter by priority (low, medium, high)
- `--tag`: Filter by specific tag

**Examples:**
```bash
task list
task list --completed
task list --priority high
task list --tag "urgent"
```

### `update` - Update Task

Update an existing task's properties.

```bash
task update [id] [flags]
```

**Flags:**
- `--title`: New task title
- `--description`: New task description
- `--priority`: New task priority
- `--tags`: New comma-separated list of tags

**Examples:**
```bash
task update 1 --title "Updated title"
task update 2 --priority high --tags "urgent,important"
```

### `complete` - Mark Task Complete

Mark a task as completed.

```bash
task complete [id]
```

**Examples:**
```bash
task complete 1
```

### `delete` - Delete Task

Delete a task permanently.

```bash
task delete [id]
```

**Examples:**
```bash
task delete 1
```

## Development

### Project Structure

```
cli-tool/
├── main.go              # Main application with CLI commands
├── go.mod              # Go module dependencies
├── README.md           # This file
└── .env.example        # Environment configuration example
```

### Key Components

1. **Task Structure**: Data model for tasks with JSON serialization
2. **TaskManager**: Business logic for task operations and persistence
3. **Cobra Commands**: CLI command definitions and handlers
4. **Color Output**: Terminal color coding for better UX
5. **Environment Config**: Environment-based configuration
6. **File Management**: JSON file persistence with directory creation

### Testing

```bash
# Run the application
go run main.go list

# Build and test
go build -o task main.go
./task add "Test task" --priority high
./task list
./task complete 1
./task list --completed
```

## Architecture Patterns

### Command Pattern
Uses Cobra framework's command pattern for organized CLI structure:
```go
var rootCmd = &cobra.Command{
    Use:   "task",
    Short: "A modern CLI task management tool",
    // ...
}
```

### Repository Pattern
TaskManager encapsulates data access and business logic:
```go
type TaskManager struct {
    tasks    []Task
    dataFile string
    nextID   int
}
```

### Configuration Management
Environment-based configuration with defaults:
```go
defaultDataFile := getEnv("TASK_DATA_FILE", 
    filepath.Join(os.Getenv("HOME"), ".tasks", "tasks.json"))
```

### Error Handling
Consistent error handling and user feedback:
```go
if err != nil {
    fmt.Printf("Error adding task: %v\n", err)
    os.Exit(1)
}
```

## Key Concepts Demonstrated

1. **CLI Development**: Using Cobra for professional CLI applications
2. **File I/O**: JSON file persistence and directory management
3. **Error Handling**: Proper error propagation and user feedback
4. **Data Modeling**: Struct design with JSON tags
5. **Environment Configuration**: Using environment variables
6. **User Experience**: Color output and intuitive commands
7. **Input Validation**: Command argument and flag validation
8. **Cross-platform Compatibility**: Path handling and file operations

## Production Considerations

For production use, consider adding:
- Configuration file support (YAML, TOML)
- Database backend (SQLite, PostgreSQL)
- Task scheduling and reminders
- Sync with external services
- Import/export functionality
- Backup and restore features
- Plugin system
- Shell completion
- Man page generation

## Learning Objectives

This example teaches:
- Building production-ready CLI applications
- Using the Cobra framework effectively
- File-based data persistence
- Environment configuration patterns
- User experience design for CLI tools
- Error handling in CLI applications
- Cross-platform Go development
- JSON data modeling and serialization