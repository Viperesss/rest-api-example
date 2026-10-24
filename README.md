# REST API Example

A simple REST API for managing TODO tasks, written in Go.

The project demonstrates RESTful HTTP endpoints, JSON request and response handling, task management, and separation of HTTP handlers from business logic.

## Features

- Create, retrieve, and delete tasks
- Retrieve all tasks or filter uncompleted tasks
- Mark tasks as completed or uncompleted
- Validate incoming task data
- Handle errors and return appropriate HTTP status codes
- Protect in-memory task storage with `sync.RWMutex`

## Tech Stack

- **Language:** Go
- **Routing:** Gorilla Mux
- **Data format:** JSON
- **Storage:** In-memory map

## Getting Started

### Prerequisites

- Go 1.27.2 or compatible version

### Installation

Clone the repository:

```bash
git clone https://github.com/Viperesss/rest-api-example.git
cd rest-api-example
```

Download dependencies:

```bash
go mod download
```

### Running the Application

```bash
go run .
```

The server starts on port `8080`.

Base URL: `http://localhost:8080`

## API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/tasks` | Create a task |
| `GET` | `/tasks` | Retrieve all tasks |
| `GET` | `/tasks/{title}` | Retrieve a task by title |
| `GET` | `/tasks?completed=false` | Retrieve uncompleted tasks |
| `PATCH` | `/tasks/{title}` | Change task completion status |
| `DELETE` | `/tasks/{title}` | Delete a task |

### Create a Task

```http
POST /tasks
Content-Type: application/json
```

Request body:

```json
{
  "Title": "Learn Go",
  "Description": "Practice REST API development"
}
```

A successful request returns `201 Created` and the created task.

### Get All Tasks

```bash
curl http://localhost:8080/tasks
```

Returns a JSON object containing the tasks.

### Get a Task by Title

```bash
curl http://localhost:8080/tasks/Learn%20Go
```

Returns the requested task or `404 Not Found` if it does not exist.

### Get Uncompleted Tasks

```bash
curl "http://localhost:8080/tasks?completed=false"
```

Returns tasks that have not been completed.

### Update Task Completion Status

```http
PATCH /tasks/Learn%20Go
Content-Type: application/json
```

Request body:

```json
{
  "Complete": true
}
```

Set `Complete` to `false` to mark the task as uncompleted.

### Delete a Task

```bash
curl -X DELETE http://localhost:8080/tasks/Learn%20Go
```

A successful request returns `204 No Content`.

## Task Model

Each task contains the following fields:

| Field | Type | Description |
|---|---|---|
| `Title` | `string` | Task title |
| `Description` | `string` | Task description |
| `Completed` | `bool` | Completion status |
| `CreatedAt` | `time.Time` | Creation timestamp |
| `CompletedAt` | `*time.Time` | Completion timestamp, or `null` |

## Project Structure

```text
.
├── http/       # HTTP server, handlers, and DTOs
├── todo/       # Task model and business logic
├── main.go     # Application entry point
├── go.mod      # Module definition and dependencies
└── go.sum      # Dependency checksums
```

## Data Storage

Tasks are stored in memory using a Go map, with task titles serving as keys. A read-write mutex protects concurrent access to the task collection.

**Note:** All tasks are lost when the application stops or restarts.
