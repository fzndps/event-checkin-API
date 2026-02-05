package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/fzndps/eventcheck/internal/domain"
	"github.com/fzndps/eventcheck/internal/domain/repository"
)

type CheckInUsecase struct {
	EventRepo      repository.EventRepository
	ParticipanRepo repository.ParticipantRepository
}

func NewCheckInUsecase(eventRepo repository.EventRepository, participantRepo repository.ParticipantRepository) *CheckInUsecase {
	return &CheckInUsecase{
		EventRepo:      eventRepo,
		ParticipanRepo: participantRepo,
	}
}

func (u *CheckInUsecase) VerifyPIN(ctx context.Context, eventSlug, scannerPIN string) (*domain.VerifyPINResponse, error) {
	// Get event dari slug
	event, err := u.EventRepo.GetBySlug(ctx, eventSlug)
	if err != nil {
		if errors.Is(err, domain.ErrEventNotFound) {
			return &domain.VerifyPINResponse{
				Valid:   false,
				Message: "Event not found",
			}, nil
		}

		return nil, err
	}

	// verifikasi PIN
	if event.ScannerPIN != scannerPIN {
		return &domain.VerifyPINResponse{
			Valid:   false,
			Message: "Invalid PIN scanner",
		}, nil
	}

	// return success response
	return &domain.VerifyPINResponse{
		Valid:     true,
		EventID:   event.ID,
		EventName: event.Name,
		Message:   "Valid PIN, Scanner authorized",
	}, nil
}

func (u *CheckInUsecase) CheckIn(ctx context.Context, qrToken string) (*domain.CheckInResponse, error) {

	// validasi qr token format
	if len(qrToken) != 32 {
		return &domain.CheckInResponse{
			Success: false,
			Message: "Invalid QR code (Wrong format)",
		}, nil
	}

	// Get partisipan dari qr token
	participant, err := u.ParticipanRepo.GetByQRToken(ctx, qrToken)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return &domain.CheckInResponse{
				Success: false,
				Message: "QR code not found, Make sure the QR code is correct.",
			}, nil
		}

		return nil, err
	}

	// cek apakah partisipan sudah pernah check in
	if participant.CheckedIn {
		return &domain.CheckInResponse{
			Success:          false,
			Message:          "Participants have already checked in",
			Participant:      participant,
			AlreadyCheckedIn: true,
			CheckedInAt:      participant.CheckedInAt,
		}, nil
	}

	err = u.ParticipanRepo.UpdateCheckIn(ctx, participant.ID)
	if err != nil {
		return nil, err
	}

	// Update participant objek
	now := time.Now()
	participant.CheckedIn = true
	participant.CheckedInAt = &now

	// return sukses response
	return &domain.CheckInResponse{
		Success:          true,
		Message:          "Checkin success! Welcome " + participant.Name,
		Participant:      participant,
		AlreadyCheckedIn: false,
		CheckedInAt:      participant.CheckedInAt,
	}, nil

}

func (u *CheckInUsecase) GetEventStats(ctx context.Context, eventSlug string) (*domain.EventStatsResponse, error) {
	// get event by slug
	event, err := u.EventRepo.GetBySlug(ctx, eventSlug)
	if err != nil {
		return nil, err
	}

	// get total partisipan
	totalParticipants, err := u.ParticipanRepo.CountByEventID(ctx, event.ID)
	if err != nil {
		return nil, err
	}

	// get checked in count
	checkedInCount, err := u.ParticipanRepo.CountCheckedInByEventID(ctx, event.ID)
	if err != nil {
		return nil, err
	}

	// hitung presentase
	var percentage float64
	if totalParticipants > 0 {
		percentage = (float64(checkedInCount) / float64(totalParticipants)) * 100
	}

	// dapatkan check in terbaru (10 terakhir)
	allParticipant, err := u.ParticipanRepo.GetByEventID(ctx, event.ID)
	if err != nil {
		return nil, err
	}

	// filter checked-in partisipan dan dapatkan 10 terakhir
	var recentCheckIns []*domain.Participant
	var lastCheckInAt *time.Time

	for _, p := range allParticipant {
		if p.CheckedIn {
			if lastCheckInAt == nil || p.CheckedInAt.After(*lastCheckInAt) {
				lastCheckInAt = p.CheckedInAt
			}

			if len(recentCheckIns) < 10 {
				recentCheckIns = append(recentCheckIns, p)
			}
		}
	}

	return &domain.EventStatsResponse{
		EventID:           event.ID,
		EventName:         event.Name,
		TotalParticipants: totalParticipants,
		CheckedInCount:    checkedInCount,
		NotCheckedInCount: totalParticipants - checkedInCount,
		CheckInPercentage: percentage,
		LastCheckInAt:     lastCheckInAt,
		RecentCheckIns:    recentCheckIns,
	}, nil
}

func (u *CheckInUsecase) GetScanPageData(ctx context.Context, eventSlug string) (*domain.ScanPageData, error) {
	// get event by slug
	event, err := u.EventRepo.GetBySlug(ctx, eventSlug)
	if err != nil {
		return nil, err
	}

	return &domain.ScanPageData{
		EventSlug: event.Slug,
		EventName: event.Name,
		EventDate: event.Date,
		Venue:     event.Venue,
	}, nil
}
