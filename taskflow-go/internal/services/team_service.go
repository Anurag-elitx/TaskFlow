package services

import (
	"context"
	"errors"
	"taskflow/internal/models"
	"taskflow/internal/repository"

	"github.com/google/uuid"
)

type TeamService struct {
	repo repository.TeamRepository
}

func NewTeamService(repo repository.TeamRepository) *TeamService {
	return &TeamService{repo: repo}
}

func (s *TeamService) CreateTeam(ctx context.Context, name, description string, ownerID uuid.UUID) (*models.Team, error) {
	team := &models.Team{
		Name:        name,
		Description: description,
		OwnerID:     ownerID,
	}
	if err := s.repo.Create(ctx, team); err != nil {
		return nil, err
	}
	if err := s.repo.AddMember(ctx, team.ID, ownerID, "owner"); err != nil {
		return nil, err
	}
	return team, nil
}

func (s *TeamService) GetUserTeams(ctx context.Context, userID uuid.UUID) ([]*models.Team, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *TeamService) AddMember(ctx context.Context, teamID, ownerID, userID uuid.UUID, role string) error {
	team, err := s.repo.GetByID(ctx, teamID)
	if err != nil {
		return err
	}
	if team == nil {
		return errors.New("team not found")
	}
	if team.OwnerID != ownerID {
		return errors.New("only team owner can add members")
	}
	return s.repo.AddMember(ctx, teamID, userID, role)
}
