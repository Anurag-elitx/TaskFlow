# TaskFlow API

A production-grade Task Management REST API built with Go.

## Overview
This API provides robust task management capabilities, including user authentication, team collaboration, and task tracking. It was built using Go for its excellent performance, concurrency model, and type safety, ensuring a highly scalable backend.

## Architecture
We use a Clean Architecture pattern:
```
Handlers (HTTP) -> Services (Business Logic) -> Repository (Data Access) -> DB
```

## Setup with Docker
1. Clone the repository.
2. Run `docker-compose up -d --build`.
3. The API is available at `http://localhost:8080`.

## Setup without Docker
1. Install Go 1.22+.
2. Set up a PostgreSQL instance and set the `DATABASE_URL` in `.env`.
3. Run `go mod tidy`.
4. Run migrations using golang-migrate: `migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5432/taskflow?sslmode=disable" up`
5. Run `go run cmd/main.go` (or use `air` for live reload).

## Endpoints
- POST `/api/v1/auth/register`: Register a new user
- POST `/api/v1/auth/login`: Login
- GET `/api/v1/users/me`: Get current user
- PATCH `/api/v1/users/me`: Update user
- POST `/api/v1/teams`: Create a team
- GET `/api/v1/teams`: Get user's teams
- POST `/api/v1/teams/:id/members`: Add a team member
- POST `/api/v1/teams/:teamId/tasks`: Create a task
- GET `/api/v1/teams/:teamId/tasks`: Get tasks
- GET `/api/v1/tasks/:id`: Get task by ID
- PATCH `/api/v1/tasks/:id`: Update task
- DELETE `/api/v1/tasks/:id`: Delete task
