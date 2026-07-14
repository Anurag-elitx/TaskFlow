package repository

import (
	"context"
	"errors"
	"fmt"
	"taskflow/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Task, error)
	GetByTeam(ctx context.Context, teamID uuid.UUID, filters models.TaskFilters) ([]*models.Task, int, error)
	Update(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type taskRepo struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) TaskRepository {
	return &taskRepo{db: db}
}

func (r *taskRepo) Create(ctx context.Context, task *models.Task) error {
	query := `
		INSERT INTO tasks (title, description, status, priority, assignee_id, team_id, due_date, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`
	
	return r.db.QueryRow(ctx, query, task.Title, task.Description, task.Status, task.Priority, 
		task.AssigneeID, task.TeamID, task.DueDate, task.CreatedBy).
		Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
}

func (r *taskRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	query := `
		SELECT id, title, description, status, priority, assignee_id, team_id, due_date, created_by, created_at, updated_at
		FROM tasks WHERE id = $1`
	var t models.Task
	err := r.db.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.AssigneeID, 
		&t.TeamID, &t.DueDate, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *taskRepo) GetByTeam(ctx context.Context, teamID uuid.UUID, filters models.TaskFilters) ([]*models.Task, int, error) {
	query := `SELECT id, title, description, status, priority, assignee_id, team_id, due_date, created_by, created_at, updated_at FROM tasks WHERE team_id = $1`
	countQuery := `SELECT COUNT(*) FROM tasks WHERE team_id = $1`
	
	args := []interface{}{teamID}
	argCount := 2

	if filters.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		countQuery += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, *filters.Status)
		argCount++
	}
	if filters.Priority != nil {
		query += fmt.Sprintf(" AND priority = $%d", argCount)
		countQuery += fmt.Sprintf(" AND priority = $%d", argCount)
		args = append(args, *filters.Priority)
		argCount++
	}
	if filters.AssigneeID != nil {
		query += fmt.Sprintf(" AND assignee_id = $%d", argCount)
		countQuery += fmt.Sprintf(" AND assignee_id = $%d", argCount)
		args = append(args, *filters.AssigneeID)
		argCount++
	}

	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (filters.Page - 1) * filters.Limit
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, filters.Limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.AssigneeID, &t.TeamID, &t.DueDate, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, &t)
	}
	return tasks, total, nil
}

func (r *taskRepo) Update(ctx context.Context, task *models.Task) error {
	query := `
		UPDATE tasks 
		SET title = $1, description = $2, status = $3, priority = $4, assignee_id = $5, due_date = $6, updated_at = NOW()
		WHERE id = $7 RETURNING updated_at`
	return r.db.QueryRow(ctx, query, task.Title, task.Description, task.Status, task.Priority, task.AssigneeID, task.DueDate, task.ID).Scan(&task.UpdatedAt)
}

func (r *taskRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	return err
}
