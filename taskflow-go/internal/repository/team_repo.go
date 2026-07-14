package repository

import (
	"context"
	"taskflow/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamRepository interface {
	Create(ctx context.Context, team *models.Team) error
	AddMember(ctx context.Context, teamID, userID uuid.UUID, role string) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Team, error)
	GetByID(ctx context.Context, teamID uuid.UUID) (*models.Team, error)
	IsMember(ctx context.Context, teamID, userID uuid.UUID) (bool, error)
}

type teamRepo struct {
	db *pgxpool.Pool
}

func NewTeamRepository(db *pgxpool.Pool) TeamRepository {
	return &teamRepo{db: db}
}

func (r *teamRepo) Create(ctx context.Context, team *models.Team) error {
	query := `INSERT INTO teams (name, description, owner_id) VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.db.QueryRow(ctx, query, team.Name, team.Description, team.OwnerID).
		Scan(&team.ID, &team.CreatedAt)
}

func (r *teamRepo) AddMember(ctx context.Context, teamID, userID uuid.UUID, role string) error {
	query := `INSERT INTO team_members (team_id, user_id, role) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(ctx, query, teamID, userID, role)
	return err
}

func (r *teamRepo) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Team, error) {
	query := `
		SELECT t.id, t.name, t.description, t.owner_id, t.created_at 
		FROM teams t 
		JOIN team_members tm ON t.id = tm.team_id 
		WHERE tm.user_id = $1`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []*models.Team
	for rows.Next() {
		var t models.Team
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.OwnerID, &t.CreatedAt); err != nil {
			return nil, err
		}
		teams = append(teams, &t)
	}
	return teams, nil
}

func (r *teamRepo) GetByID(ctx context.Context, teamID uuid.UUID) (*models.Team, error) {
	query := `SELECT id, name, description, owner_id, created_at FROM teams WHERE id = $1`
	var team models.Team
	err := r.db.QueryRow(ctx, query, teamID).Scan(&team.ID, &team.Name, &team.Description, &team.OwnerID, &team.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *teamRepo) IsMember(ctx context.Context, teamID, userID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM team_members WHERE team_id = $1 AND user_id = $2)`
	var exists bool
	err := r.db.QueryRow(ctx, query, teamID, userID).Scan(&exists)
	return exists, err
}
