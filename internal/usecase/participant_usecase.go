package usecase

import (
	"context"
	"fmt"
	"io"

	"github.com/fzndps/eventcheck/internal/domain"
	"github.com/fzndps/eventcheck/internal/domain/repository"
	"github.com/fzndps/eventcheck/pkg/csv"
	"github.com/fzndps/eventcheck/pkg/random"
)

type ParticipantUsecase struct {
	eventRepo      repository.EventRepository
	participanRepo repository.ParticipantRepository
}

func NewParticipantUsecase(eventRepo repository.EventRepository, participanRepo repository.ParticipantRepository) *ParticipantUsecase {
	return &ParticipantUsecase{
		eventRepo:      eventRepo,
		participanRepo: participanRepo,
	}
}

// Menangani untuk upload CSV participants
func (u *ParticipantUsecase) UploadParticipants(
	ctx context.Context,
	organizerID int64,
	eventID string,
	csvReader io.Reader,
) (*domain.UploadParticipantsResponse, error) {
	// Cek authorization untuk memastikan event milik organizer
	isOwned, err := u.eventRepo.IsOwnedBy(ctx, eventID, organizerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get owner event: %w", err)
	}

	if !isOwned {
		return nil, domain.ErrUnauthorizedAccess
	}

	// Parse CSV
	participants, parseError, err := csv.ParseParticipants(csvReader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse csd: %w", err)
	}

	// Generate unique QR code untuk setiap participant
	for _, p := range participants {
		token, err := random.GenerateToken()
		if err != nil {
			return nil, fmt.Errorf("failed to generate token: %w", err)
		}
		p.QRToken = token
		p.EventID = eventID
	}

	// bulk insert ke database
	var successCount int
	var failedReason []string

	if len(participants) > 0 {
		err = u.participanRepo.BulkCreate(ctx, participants)
		if err != nil {
			return nil, fmt.Errorf("failed to bulk insert: %w", err)
		}
		successCount = len(participants)
	}

	res := &domain.UploadParticipantsResponse{
		Success:       successCount,
		Failed:        len(parseError),
		FailedReasons: failedReason,
	}

	return res, nil
}

// Menangani list partisipan
func (u *ParticipantUsecase) GetParticipantsByEvent(
	ctx context.Context,
	organizerID int64,
	eventID string,
) ([]*domain.Participant, error) {
	isOwned, err := u.eventRepo.IsOwnedBy(ctx, eventID, organizerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get owner event: %w", err)
	}

	if !isOwned {
		return nil, domain.ErrUnauthorizedAccess
	}

	participants, err := u.participanRepo.GetByEventID(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get participants by event ID %s: %w", eventID, err)
	}

	return participants, nil
}
