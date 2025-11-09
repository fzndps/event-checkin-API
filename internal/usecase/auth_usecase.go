package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fzndps/eventcheck/config"
	"github.com/fzndps/eventcheck/internal/domain"
	"github.com/fzndps/eventcheck/internal/domain/repository"
	"github.com/fzndps/eventcheck/pkg/hash"
	"github.com/fzndps/eventcheck/pkg/jwt"
)

type AuthUsecase struct {
	organizerRepo repository.OrganizerRepository
	jwtManager    *jwt.JWTManager
	cfg           *config.Config
}

func NewAuthUsecase(organizerRepo repository.OrganizerRepository, jwtManager *jwt.JWTManager, cfg *config.Config) *AuthUsecase {
	return &AuthUsecase{
		organizerRepo: organizerRepo,
		jwtManager:    jwtManager,
		cfg:           cfg,
	}
}

func (u *AuthUsecase) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.Organizer, error) {
	existingOrganizer, err := u.organizerRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingOrganizer != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := hash.HashPassword(req.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	organizer := &domain.Organizer{
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
	}

	err = u.organizerRepo.Create(ctx, organizer)
	if err != nil {
		return nil, fmt.Errorf("failed to create organizer: %v", err)
	}

	fullOrganizer, err := u.organizerRepo.GetByID(ctx, organizer.ID)
	if err != nil {
		return nil, fmt.Errorf("something went wrong: %v", err)
	}

	return fullOrganizer, nil
}

func (u *AuthUsecase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.OrganizerLoginResponse, error) {
	organizer, err := u.organizerRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !hash.CheckPassword(req.PasswordHash, organizer.PasswordHash) {
		return nil, errors.New("invalid email or password")
	}

	token, err := u.jwtManager.GenerateToken(organizer.ID, organizer.Email, u.cfg.JWT.Expiry)
	if err != nil {
		return nil, err
	}

	res := &domain.OrganizerLoginResponse{
		Token:     token,
		Organizer: organizer,
	}

	return res, nil
}

func (u *AuthUsecase) GetProfileByID(ctx context.Context, organizerID int) (*domain.Organizer, error) {
	organizer, err := u.organizerRepo.GetByID(ctx, organizerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get organizer by id %d: %v", organizerID, err)
	}

	return organizer, nil
}
