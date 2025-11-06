package mysql

import (
	"context"
	"testing"
	"time"

	"github.com/fzndps/eventcheck/config"
	"github.com/fzndps/eventcheck/internal/domain"
	"github.com/fzndps/eventcheck/internal/infrastructure/database"
	"github.com/google/uuid"
)

func setupTestEventRepo(t *testing.T) *eventRepository {
	cfg, err := config.LoadConfig("../../../.env")
	if err != nil {
		t.Fatal("Failed to load config:", err)
	}

	db, err := database.InitDB(cfg)
	if err != nil {
		t.Fatal("Failed to connect to database:", err)
	}

	// Note: Don't close DB here, let test handle it
	return &eventRepository{db: db}
}

func TestEventRepository_Create(t *testing.T) {
	repo := setupTestEventRepo(t)
	defer repo.db.Close()

	event := &domain.Event{
		ID:               uuid.New().String(),
		OrganizerID:      1,
		Name:             "Test Event",
		Slug:             "test-event-" + time.Now().Format("20060102150405"),
		Date:             time.Now().Add(24 * time.Hour),
		Venue:            "Test Venue",
		ParticipantCount: 100,
		TotalPrice:       450000,
		PaymentStatus:    domain.PaymentStatusPending,
		ScannerPIN:       "1234",
	}

	err := repo.Create(context.Background(), event)
	if err != nil {
		t.Fatal("Failed to create event:", err)
	}

	// Cleanup
	repo.Delete(context.Background(), event.ID)

	t.Log("✅ Event created successfully")
}

func TestEventRepository_GetByID(t *testing.T) {
	repo := setupTestEventRepo(t)
	defer repo.db.Close()

	// Create event first
	event := &domain.Event{
		ID:               uuid.New().String(),
		OrganizerID:      1,
		Name:             "Test Event",
		Slug:             "test-event-" + time.Now().Format("20060102150405"),
		Date:             time.Now().Add(24 * time.Hour),
		Venue:            "Test Venue",
		ParticipantCount: 100,
		TotalPrice:       450000,
		PaymentStatus:    domain.PaymentStatusPending,
		ScannerPIN:       "1234",
	}

	repo.Create(context.Background(), event)
	defer repo.Delete(context.Background(), event.ID)

	// Get by ID
	found, err := repo.GetByID(context.Background(), event.ID)
	if err != nil {
		t.Fatal("Failed to get event:", err)
	}

	if found.ID != event.ID {
		t.Errorf("Expected ID %s, got %s", event.ID, found.ID)
	}
	if found.Name != event.Name {
		t.Errorf("Expected name %s, got %s", event.Name, found.Name)
	}

	t.Log("✅ Event retrieved successfully")
}

func TestEventRepository_GetBySlug(t *testing.T) {
	repo := setupTestEventRepo(t)
	defer repo.db.Close()

	slug := "test-slug-" + time.Now().Format("20060102150405")

	event := &domain.Event{
		ID:               uuid.New().String(),
		OrganizerID:      1,
		Name:             "Test Event",
		Slug:             slug,
		Date:             time.Now().Add(24 * time.Hour),
		Venue:            "Test Venue",
		ParticipantCount: 100,
		TotalPrice:       450000,
		PaymentStatus:    domain.PaymentStatusPending,
		ScannerPIN:       "1234",
	}

	repo.Create(context.Background(), event)
	defer repo.Delete(context.Background(), event.ID)

	// Get by slug
	found, err := repo.GetBySlug(context.Background(), slug)
	if err != nil {
		t.Fatal("Failed to get event by slug:", err)
	}

	if found.Slug != slug {
		t.Errorf("Expected slug %s, got %s", slug, found.Slug)
	}

	t.Log("✅ Event retrieved by slug successfully")
}

func TestEventRepository_GetByOrganizerID(t *testing.T) {
	repo := setupTestEventRepo(t)
	defer repo.db.Close()

	organizer := &domain.Organizer{
		ID:           999,
		Email:        "test-org-for-get@example.com",
		PasswordHash: "hashedpassword",
		Name:         "Test Org",
		CreatedAt:    time.Now(),
	}

	organizerRepo := NewOrganizerRepositoryImpl(repo.db)
	err := organizerRepo.Create(context.Background(), organizer)
	if err != nil {
		t.Fatalf("Failed to create prerequisite organizer: %v", err)
	}

	organizerID := int64(organizer.ID) // Use unique ID for test

	// Create 3 events
	eventIDs := make([]string, 3)
	for i := 0; i < 3; i++ {
		event := &domain.Event{
			ID:               uuid.New().String(),
			OrganizerID:      organizerID,
			Name:             "Test Event",
			Slug:             "test-event-" + time.Now().Format("20060102150405") + string(rune('a'+i)),
			Date:             time.Now().Add(24 * time.Hour),
			Venue:            "Test Venue",
			ParticipantCount: 100,
			TotalPrice:       450000,
			PaymentStatus:    domain.PaymentStatusPending,
			ScannerPIN:       "1234",
		}
		repo.Create(context.Background(), event)
		eventIDs[i] = event.ID
	}

	// Get by organizer ID with pagination
	events, total, err := repo.GetByOrganizerID(context.Background(), organizerID, 10, 0)
	if err != nil {
		t.Fatal("Failed to get events:", err)
	}

	if len(events) != 3 {
		t.Errorf("Expected 3 events, got %d", len(events))
	}

	if total != 3 {
		t.Errorf("Expected total 3, got %d", total)
	}

	// Cleanup
	for _, id := range eventIDs {
		repo.Delete(context.Background(), id)
	}

	t.Log("✅ Events retrieved by organizer ID successfully")
}

func TestEventRepository_Update(t *testing.T) {
	repo := setupTestEventRepo(t)
	defer repo.db.Close()

	// Create event
	event := &domain.Event{
		ID:               uuid.New().String(),
		OrganizerID:      1,
		Name:             "Original Name",
		Slug:             "original-slug-" + time.Now().Format("20060102150405"),
		Date:             time.Now().Add(24 * time.Hour),
		Venue:            "Original Venue",
		ParticipantCount: 100,
		TotalPrice:       450000,
		PaymentStatus:    domain.PaymentStatusPending,
		ScannerPIN:       "1234",
	}

	repo.Create(context.Background(), event)
	defer repo.Delete(context.Background(), event.ID)

	// Update event
	event.Name = "Updated Name"
	event.Venue = "Updated Venue"

	err := repo.Update(context.Background(), event)
	if err != nil {
		t.Fatal("Failed to update event:", err)
	}

	// Verify update
	updated, _ := repo.GetByID(context.Background(), event.ID)
	if updated.Name != "Updated Name" {
		t.Errorf("Expected name 'Updated Name', got '%s'", updated.Name)
	}
	if updated.Venue != "Updated Venue" {
		t.Errorf("Expected venue 'Updated Venue', got '%s'", updated.Venue)
	}

	t.Log("✅ Event updated successfully")
}

func TestEventRepository_Delete(t *testing.T) {
	repo := setupTestEventRepo(t)
	defer repo.db.Close()

	// Create event
	event := &domain.Event{
		ID:               uuid.New().String(),
		OrganizerID:      1,
		Name:             "Test Event",
		Slug:             "test-event-" + time.Now().Format("20060102150405"),
		Date:             time.Now().Add(24 * time.Hour),
		Venue:            "Test Venue",
		ParticipantCount: 100,
		TotalPrice:       450000,
		PaymentStatus:    domain.PaymentStatusPending,
		ScannerPIN:       "1234",
	}

	repo.Create(context.Background(), event)

	// Delete event
	err := repo.Delete(context.Background(), event.ID)
	if err != nil {
		t.Fatal("Failed to delete event:", err)
	}

	// Verify deleted
	_, err = repo.GetByID(context.Background(), event.ID)
	if err != domain.ErrEventNotFound {
		t.Error("Expected ErrEventNotFound after delete")
	}

	t.Log("✅ Event deleted successfully")
}

func TestEventRepository_IsOwnedBy(t *testing.T) {
	repo := setupTestEventRepo(t)
	defer repo.db.Close()

	organizer := &domain.Organizer{
		ID:           123,
		Email:        "test-org-for-delete@example.com",
		PasswordHash: "hashedpassword",
		Name:         "Test Org Del",
		CreatedAt:    time.Now(),
	}

	organizerRepo := NewOrganizerRepositoryImpl(repo.db)
	err := organizerRepo.Create(context.Background(), organizer)
	if err != nil {
		t.Fatalf("Failed to create prerequisite organizer: %v", err)
	}

	organizerID := int64(organizer.ID) // Use unique ID for test

	event := &domain.Event{
		ID:               uuid.New().String(),
		OrganizerID:      organizerID,
		Name:             "Test Event",
		Slug:             "test-event-" + time.Now().Format("20060102150405"),
		Date:             time.Now().Add(24 * time.Hour),
		Venue:            "Test Venue",
		ParticipantCount: 100,
		TotalPrice:       450000,
		PaymentStatus:    domain.PaymentStatusPending,
		ScannerPIN:       "1234",
	}

	repo.Create(context.Background(), event)
	defer repo.Delete(context.Background(), event.ID)

	// Test correct owner
	isOwned, err := repo.IsOwnedBy(context.Background(), event.ID, organizerID)
	if err != nil {
		t.Fatal("Failed to check ownership:", err)
	}
	if !isOwned {
		t.Error("Expected true for correct owner")
	}

	// Test wrong owner
	isOwned, err = repo.IsOwnedBy(context.Background(), event.ID, 999)
	if err != nil {
		t.Fatal("Failed to check ownership:", err)
	}
	if isOwned {
		t.Error("Expected false for wrong owner")
	}

	t.Log("✅ Ownership check working correctly")
}
