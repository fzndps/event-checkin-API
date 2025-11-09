package repository

import (
	"context"

	"github.com/fzndps/eventcheck/internal/domain"
)

// ParticipantRepository adalah interface untuk akses data participant
type ParticipantRepository interface {
	// Create menyimpan satu participant
	Create(ctx context.Context, participant *domain.Participant) error

	// BulkCreate menyimpan banyak participants sekaligus (untuk CSV upload)
	// Menggunakan transaction untuk atomicity: all or nothing
	BulkCreate(ctx context.Context, participants []*domain.Participant) error

	// GetByEventID mencari semua participant di event tertentu
	GetByEventID(ctx context.Context, eventID string) ([]*domain.Participant, error)

	// GetByQRToken mencari participant berdasarkan QR token (untuk check-in)
	GetByQRToken(ctx context.Context, qrToken string) (*domain.Participant, error)

	// CountByEventID menghitung jumlah participant di event
	CountByEventID(ctx context.Context, eventID string) (int, error)

	// CountCheckedInByEventID menghitung jumlah participant yang sudah check-in
	CountCheckedInByEventID(ctx context.Context, eventID string) (int, error)

	// UpdateCheckIn mengupdate status check-in participant
	UpdateCheckIn(ctx context.Context, participantID int64) error

	// DeleteByEventID menghapus semua participant di event (cascade delete)
	DeleteByEventID(ctx context.Context, eventID string) error
}
