package repository

import (
	"context"

	"github.com/fzndps/eventcheck/internal/domain"
)

type OrganizerRepository interface {
	Create(ctx context.Context, organizer *domain.Organizer) error
	GetByEmail(ctx context.Context, email string) (*domain.Organizer, error)
	GetByID(ctx context.Context, id int64) (*domain.Organizer, error)
}
