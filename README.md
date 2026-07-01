# TaskFlow

TaskFlow is a robust, production-grade Task Management system featuring a high-performance REST API backend built in Go. It supports user authentication, team collaboration, and comprehensive task tracking.

## Project Structure

- `taskflow-go/`: The core Go backend application.

## Features

- **User Authentication**: Secure registration and login.
- **Team Collaboration**: Create teams and manage team members.
- **Task Management**: Create, view, update, and delete tasks within teams.
- **Clean Architecture**: Follows a scalable Handlers -> Services -> Repository -> DB architecture.

## Tech Stack

- **Language**: Go 1.22+
- **Database**: PostgreSQL
- **Infrastructure**: Docker & Docker Compose
- **Migrations**: golang-migrate
- **Live Reload**: Air

## Quick Start (Docker)

The easiest way to run the application is using Docker.

1. Clone the repository:
   ```bash
   git clone https://github.com/Anurag-elitx/TaskFlow.git
   cd TaskFlow/taskflow-go
   ```

2. Start the services:
   ```bash
   docker-compose up -d --build
   ```

3. The API will be available at `http://localhost:8080`.

## Manual Setup

If you prefer to run the application without Docker:

1. Install Go 1.22+.
2. Set up a PostgreSQL database.
3. Navigate to the `taskflow-go` directory:
   ```bash
   cd taskflow-go
   ```
4. Copy the environment variables example and configure your database connection:
   ```bash
   cp .env.example .env
   ```
5. Install dependencies:
   ```bash
   go mod tidy
   ```
6. Run database migrations:
   ```bash
   migrate -path db/migrations -database "postgres://postgres:postgres@localhost:5432/taskflow?sslmode=disable" up
   ```
7. Start the server:
   ```bash
   go run cmd/main.go
   ```
   *(Alternatively, use `air` for hot-reloading during development).*

## API Documentation

The REST API exposes the following primary endpoints under `/api/v1/`:

**Auth**
- `POST /auth/register` - Register a new user
- `POST /auth/login` - Login

**Users**
- `GET /users/me` - Get current user profile
- `PATCH /users/me` - Update profile

**Teams**
- `POST /teams` - Create a team
- `GET /teams` - List user's teams
- `POST /teams/:id/members` - Add a member to a team

**Tasks**
- `POST /teams/:teamId/tasks` - Create a new task
- `GET /teams/:teamId/tasks` - List team tasks
- `GET /tasks/:id` - Get a specific task
- `PATCH /tasks/:id` - Update a task
- `DELETE /tasks/:id` - Delete a task

## License

This project is open-source and available under the MIT License.