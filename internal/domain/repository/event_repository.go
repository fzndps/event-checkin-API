package repository

import (
	"context"

	"github.com/fzndps/eventcheck/internal/domain"
)

// EventRepository adalah interface untuk akses data event
type EventRepository interface {
	// Create menyimpan event baru ke database
	Create(ctx context.Context, event *domain.Event) error

	// GetByID mencari event berdasarkan ID
	GetByID(ctx context.Context, id string) (*domain.Event, error)

	// GetBySlug mencari event berdasarkan slug
	GetBySlug(ctx context.Context, slug string) (*domain.Event, error)

	// GetByOrganizerID mencari semua event milik organizer dengan pagination
	// offset = (page - 1) * limit
	GetByOrganizerID(ctx context.Context, organizerID int64, limit, offset int) ([]*domain.Event, int, error)

	// Update mengupdate data event
	Update(ctx context.Context, event *domain.Event) error

	// Delete menghapus event (dan cascade delete participants)
	Delete(ctx context.Context, id string) error

	// IsOwnedBy mengecek apakah event dimiliki oleh organizer
	IsOwnedBy(ctx context.Context, eventID string, organizerID int64) (bool, error)
}
