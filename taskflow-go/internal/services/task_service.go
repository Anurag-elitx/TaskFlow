package services

import (
	"context"
	"errors"
	"taskflow/internal/models"
	"taskflow/internal/repository"

	"github.com/google/uuid"
)

type TaskService struct {
	taskRepo repository.TaskRepository
	teamRepo repository.TeamRepository
}

func NewTaskService(taskRepo repository.TaskRepository, teamRepo repository.TeamRepository) *TaskService {
	return &TaskService{taskRepo: taskRepo, teamRepo: teamRepo}
}

func (s *TaskService) CreateTask(ctx context.Context, task *models.Task) error {
	// Validate membership
	isMember, err := s.teamRepo.IsMember(ctx, task.TeamID, task.CreatedBy)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.New("user is not a member of this team")
	}
	if task.Status == "" {
		task.Status = "todo"
	}
	if task.Priority == "" {
		task.Priority = "medium"
	}
	return s.taskRepo.Create(ctx, task)
}

func (s *TaskService) GetTasksByTeam(ctx context.Context, teamID, userID uuid.UUID, filters models.TaskFilters) ([]*models.Task, int, error) {
	isMember, err := s.teamRepo.IsMember(ctx, teamID, userID)
	if err != nil {
		return nil, 0, err
	}
	if !isMember {
		return nil, 0, errors.New("user is not a member of this team")
	}
	if filters.Page <= 0 {
		filters.Page = 1
	}
	if filters.Limit <= 0 {
		filters.Limit = 20
	}
	return s.taskRepo.GetByTeam(ctx, teamID, filters)
}

func (s *TaskService) GetTaskByID(ctx context.Context, id, userID uuid.UUID) (*models.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, errors.New("task not found")
	}
	
	isMember, err := s.teamRepo.IsMember(ctx, task.TeamID, userID)
	if err != nil || !isMember {
		return nil, errors.New("access denied")
	}
	return task, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, taskID, userID uuid.UUID, updates *models.Task) (*models.Task, error) {
	existing, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("task not found")
	}

	isAssignee := existing.AssigneeID != nil && *existing.AssigneeID == userID
	isCreator := existing.CreatedBy == userID

	if !isCreator && !isAssignee {
		return nil, errors.New("only creator or assignee can update task")
	}

	if updates.Title != "" {
		existing.Title = updates.Title
	}
	if updates.Description != "" {
		existing.Description = updates.Description
	}
	if updates.Status != "" {
		existing.Status = updates.Status
	}
	if updates.Priority != "" {
		existing.Priority = updates.Priority
	}
	if updates.AssigneeID != nil {
		existing.AssigneeID = updates.AssigneeID
	}
	if updates.DueDate != nil {
		existing.DueDate = updates.DueDate
	}

	err = s.taskRepo.Update(ctx, existing)
	return existing, err
}

func (s *TaskService) DeleteTask(ctx context.Context, taskID, userID uuid.UUID) error {
	existing, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("task not found")
	}
	if existing.CreatedBy != userID {
		return errors.New("only creator can delete task")
	}
	return s.taskRepo.Delete(ctx, taskID)
}
