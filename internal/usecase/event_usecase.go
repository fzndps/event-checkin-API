package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/fzndps/eventcheck/internal/domain"
	"github.com/fzndps/eventcheck/internal/domain/repository"
	"github.com/fzndps/eventcheck/pkg/random"
	"github.com/fzndps/eventcheck/pkg/slug"
	"github.com/google/uuid"
)

type EventUsecase struct {
	eventRepo      repository.EventRepository
	participanRepo repository.ParticipantRepository
}

func NewEventUsecase(eventRepo repository.EventRepository, participantRepo repository.ParticipantRepository) *EventUsecase {
	return &EventUsecase{
		eventRepo:      eventRepo,
		participanRepo: participantRepo,
	}
}

// Menangani untuk create event
func (u *EventUsecase) CreateEvent(
	ctx context.Context,
	organizerID int64,
	req *domain.CreateEventRequest,
) (*domain.Event, error) {
	// validasi input
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// generate UUID untuk event ID
	eventID := uuid.New().String()

	// Generate slug dari event name
	eventSlug := slug.Generate(req.Name)

	// Cek slug apakah sudah ada, jika ya generate unique slug
	existingEvent, err := u.eventRepo.GetBySlug(ctx, eventSlug)
	if err != nil && !errors.Is(err, domain.ErrEventNotFound) {
		return nil, fmt.Errorf("failed to get slug, slug not found: %w", err)
	}

	if existingEvent != nil {
		// Slug sudah ada, generate unique slug dengan timestamp
		eventSlug = slug.GenerateUnique(req.Name)
	}

	// Gnenerate random 4 digit scanner PIN dengan crypto/rand karena lebih secure
	scannerPIN, err := random.GeneratePIN()
	if err != nil {
		return nil, fmt.Errorf("failed generate random PIN: %w", err)
	}

	// Kalkulasi total price berdasarkan partisipan
	totalPrice := domain.CalculatePrice(req.ParticipantCount)

	// buat object event
	event := &domain.Event{
		ID:               eventID,
		OrganizerID:      int64(organizerID),
		Name:             req.Name,
		Slug:             eventSlug,
		Date:             req.Date.Time,
		Venue:            req.Venue,
		ParticipantCount: req.ParticipantCount,
		TotalPrice:       totalPrice,
		PaymentStatus:    domain.PaymentStatusPending,
		ScannerPIN:       scannerPIN,
	}

	if err := u.eventRepo.Create(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	return event, nil

}

// Menangani list events dengan pagination
func (u *EventUsecase) GetEventByOrganizer(
	ctx context.Context,
	organizerID int64,
	page,
	limit int,
) (*domain.EventListResponse, error) {
	// Memvalidari paginaton
	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 100 {
		limit = 10 //default 10 ituems per page
	}

	// kalkulasi offset
	offset := (page - 1) * limit

	// Get events dari repo
	events, total, err := u.eventRepo.GetByOrganizerID(ctx, organizerID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get event by organizerID: %w", err)
	}

	// kalkulasi total pages
	totalPgaes := total / limit
	if total%limit != 0 {
		totalPgaes++
	}

	// response
	res := &domain.EventListResponse{
		Events:    events,
		Total:     total,
		Page:      page,
		Limit:     limit,
		TotalPage: totalPgaes,
	}

	return res, nil
}

// Menangani get evebt detail dengan partisipan
func (u *EventUsecase) GetEventDetail(
	ctx context.Context,
	organizerID int64,
	eventID string,
) (*domain.EventDetailResponse, error) {
	// Dapatkan event dari repo
	event, err := u.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	// Pastikan event milik organizer
	if event.OrganizerID != organizerID {
		return nil, domain.ErrUnauthorizedAccess
	}

	// dapatkan partisipan dari repo
	participant, err := u.participanRepo.GetByEventID(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get all participant: %w", err)
	}

	// Hitung statistik
	participantRegistered, err := u.participanRepo.CountByEventID(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get total participant in event: %w", err)
	}

	participantCheckedIn, err := u.participanRepo.CountCheckedInByEventID(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to count all participant: %w", err)
	}

	// return response
	res := &domain.EventDetailResponse{
		Event:                 event,
		Participants:          participant,
		ParticipantRegistered: participantRegistered,
		ParticipantCheckedIn:  participantCheckedIn,
	}

	return res, nil
}

func (u *EventUsecase) UpdateEvent(
	ctx context.Context,
	organizerID int64,
	eventID string,
	req *domain.UpdateEventRequest,
) (*domain.Event, error) {
	// Cek event apakah ada
	event, err := u.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	// Cek authorization
	if event.OrganizerID != organizerID {
		return nil, domain.ErrUnauthorizedAccess
	}

	// Update fields
	if req.Name != "" {
		event.Name = req.Name
		// generate ulang slug jika nama berubah
		newSlug := slug.Generate(req.Name)
		if newSlug != event.Slug {
			// cek apakah slug baru sudah ada
			existingEvent, err := u.eventRepo.GetBySlug(ctx, newSlug)
			if err != nil && !errors.Is(err, domain.ErrEventNotFound) {
				return nil, fmt.Errorf("failed to get new slug: %w", err)
			}

			if existingEvent != nil && existingEvent.ID != eventID {
				// Slug sudah dipakai event lain
				newSlug = slug.GenerateUnique(req.Name)
			}
			event.Slug = newSlug
		}
	}

	if !req.Date.IsZero() {
		event.Date = req.Date.Time
	}

	if req.Venue != "" {
		event.Venue = req.Venue
	}

	// save update ke database
	if err := u.eventRepo.Update(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to update event: %w", err)
	}

	return event, nil
}

func (u *EventUsecase) DeleteEvent(
	ctx context.Context,
	organizerID int64,
	eventID string,
) error {
	// check authorization dengan method IsOwnedBy
	isOwned, err := u.eventRepo.IsOwnedBy(ctx, eventID, organizerID)
	if err != nil {
		return fmt.Errorf("failed to get owner event: %w", err)
	}

	if !isOwned {
		return domain.ErrUnauthorizedAccess
	}

	if err := u.eventRepo.Delete(ctx, eventID); err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	return nil

}
