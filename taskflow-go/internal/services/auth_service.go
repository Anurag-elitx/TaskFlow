package services

import (
	"context"
	"errors"
	"taskflow/internal/models"
	"taskflow/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewAuthService(repo repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtSecret}
}

func (s *AuthService) Register(ctx context.Context, email, password, fullName string) (*models.User, string, error) {
	// Check if user exists
	existingUser, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}
	if existingUser != nil {
		return nil, "", errors.New("email already in use")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, "", err
	}

	user := &models.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		FullName:     fullName,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, "", err
	}

	token, err := s.generateToken(user.ID.String())
	return user, token, err
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*models.User, string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}
	if user == nil {
		return nil, "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	token, err := s.generateToken(user.ID.String())
	return user, token, err
}

func (s *AuthService) generateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
