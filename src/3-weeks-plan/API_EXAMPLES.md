# API Testing Examples

This file contains sample curl commands to test the Task Management API.

## Health Check

```bash
curl http://localhost:8080/health
```

## Create Tasks

### Create a high priority task
```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Complete 3-weeks Go plan",
    "description": "Learn Go in 3 weeks",
    "priority": "high"
  }'
```

### Create a medium priority task
```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Build REST API",
    "description": "Implement task management API",
    "priority": "medium"
  }'
```

### Create a task with due date
```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Submit project",
    "description": "Final project submission",
    "priority": "high",
    "due_date": "2024-12-31T23:59:59Z"
  }'
```

## List Tasks

### List all tasks
```bash
curl http://localhost:8080/api/tasks
```

### List tasks by status
```bash
curl http://localhost:8080/api/tasks?status=pending
curl http://localhost:8080/api/tasks?status=in_progress
curl http://localhost:8080/api/tasks?status=completed
```

### List tasks by priority
```bash
curl http://localhost:8080/api/tasks?priority=high
curl http://localhost:8080/api/tasks?priority=medium
curl http://localhost:8080/api/tasks?priority=low
```

### Combine filters
```bash
curl http://localhost:8080/api/tasks?status=pending&priority=high
```

## Get Task

### Get a specific task (replace {id} with actual task ID)
```bash
curl http://localhost:8080/api/tasks/1
```

## Update Task

### Update task status
```bash
curl -X PUT http://localhost:8080/api/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{
    "status": "in_progress"
  }'
```

### Update task to completed
```bash
curl -X PUT http://localhost:8080/api/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{
    "status": "completed"
  }'
```

### Update multiple fields
```bash
curl -X PUT http://localhost:8080/api/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated title",
    "description": "Updated description",
    "status": "completed",
    "priority": "low"
  }'
```

## Delete Task

### Delete a task (replace {id} with actual task ID)
```bash
curl -X DELETE http://localhost:8080/api/tasks/1
```

## Statistics

### Get task statistics
```bash
curl http://localhost:8080/api/stats
```

## Pretty Print with jq

Add `| jq .` to any command for pretty-printed JSON:

```bash
curl http://localhost:8080/api/tasks | jq .
curl http://localhost:8080/api/stats | jq .
```

## Complete Workflow Example

```bash
# 1. Check health
curl http://localhost:8080/health

# 2. Create tasks
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Task 1","priority":"high"}'

curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Task 2","priority":"medium"}'

# 3. List all tasks
curl http://localhost:8080/api/tasks | jq .

# 4. Update first task
curl -X PUT http://localhost:8080/api/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"status":"completed"}'

# 5. Get statistics
curl http://localhost:8080/api/stats | jq .

# 6. Delete a task
curl -X DELETE http://localhost:8080/api/tasks/2
```
