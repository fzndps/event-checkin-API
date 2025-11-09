package mysql

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/fzndps/eventcheck/config"
	"github.com/fzndps/eventcheck/internal/domain"
	"github.com/fzndps/eventcheck/internal/infrastructure/database"
	"github.com/google/uuid"
)

func setupTestParticipantRepo(t *testing.T) (*participantRepository, string) {
	cfg, err := config.LoadConfig("../../../.env")
	if err != nil {
		t.Fatal("Failed to load config:", err)
	}

	db, err := database.InitDB(cfg)
	if err != nil {
		t.Fatal("Failed to connect to database:", err)
	}

	// Create test event first
	eventID := uuid.New().String()
	_, err = db.Exec(`
		INSERT INTO events (id, organizer_id, name, slug, date, venue, participant_count, total_price, payment_status, scanner_pin)
		VALUES (?, 1, 'Test Event', ?, NOW(), 'Test Venue', 100, 450000, 'pending', '1234')
	`, eventID, "test-"+eventID)

	if err != nil {
		t.Fatal("Failed to create test event:", err)
	}

	return &participantRepository{db: db}, eventID
}

func cleanupTestEvent(t *testing.T, repo *participantRepository, eventID string) {
	// Delete participants first (foreign key)
	repo.DeleteByEventID(context.Background(), eventID)

	// Delete event
	_, err := repo.db.Exec("DELETE FROM events WHERE id = ?", eventID)
	if err != nil {
		t.Log("Warning: Failed to cleanup event:", err)
	}

	repo.db.Close()
}

func TestParticipantRepository_Create(t *testing.T) {
	repo, eventID := setupTestParticipantRepo(t)
	defer cleanupTestEvent(t, repo, eventID)

	participant := &domain.Participant{
		EventID: eventID,
		Name:    "John Doe",
		Email:   "john@example.com",
		Phone:   "08123456789",
		QRToken: "test-token-" + uuid.New().String(),
	}

	err := repo.Create(context.Background(), participant)
	if err != nil {
		t.Fatal("Failed to create participant:", err)
	}

	if participant.ID == 0 {
		t.Error("Participant ID should be set after create")
	}

	t.Log("✅ Participant created successfully with ID:", participant.ID)
}

func TestParticipantRepository_BulkCreate(t *testing.T) {
	repo, eventID := setupTestParticipantRepo(t)
	defer cleanupTestEvent(t, repo, eventID)

	// Create 100 participants
	participants := make([]*domain.Participant, 100)
	for i := 0; i < 100; i++ {
		participants[i] = &domain.Participant{
			EventID: eventID,
			Name:    "Participant " + string(rune('0'+i%10)),
			Email:   "user" + string(rune('0'+i%10)) + "@example.com",
			Phone:   "0812345678" + string(rune('0'+i%10)),
			QRToken: uuid.New().String(),
		}
	}

	err := repo.BulkCreate(context.Background(), participants)
	if err != nil {
		t.Fatal("Failed to bulk create participants:", err)
	}

	// Verify count
	count, err := repo.CountByEventID(context.Background(), eventID)
	if err != nil {
		t.Fatal("Failed to count participants:", err)
	}

	b, err := json.MarshalIndent(participants, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))

	if count != 100 {
		t.Errorf("Expected 100 participants, got %d", count)
	}

	t.Log("✅ Bulk created 100 participants successfully")
}

func TestParticipantRepository_GetByEventID(t *testing.T) {
	repo, eventID := setupTestParticipantRepo(t)
	defer cleanupTestEvent(t, repo, eventID)

	// Create 3 participants
	for i := 0; i < 3; i++ {
		p := &domain.Participant{
			EventID: eventID,
			Name:    "User " + string(rune('A'+i)),
			Email:   "user" + string(rune('a'+i)) + "@example.com",
			Phone:   "0812345678" + string(rune('0'+i)),
			QRToken: uuid.New().String(),
		}
		repo.Create(context.Background(), p)
	}

	// Get by event ID
	participants, err := repo.GetByEventID(context.Background(), eventID)
	if err != nil {
		t.Fatal("Failed to get participants:", err)
	}

	if len(participants) != 3 {
		t.Errorf("Expected 3 participants, got %d", len(participants))
	}

	t.Log("✅ Retrieved participants by event ID successfully")
}

func TestParticipantRepository_GetByQRToken(t *testing.T) {
	repo, eventID := setupTestParticipantRepo(t)
	defer cleanupTestEvent(t, repo, eventID)

	qrToken := uuid.New().String()

	participant := &domain.Participant{
		EventID: eventID,
		Name:    "John Doe",
		Email:   "john@example.com",
		Phone:   "08123456789",
		QRToken: qrToken,
	}

	repo.Create(context.Background(), participant)

	// Get by QR token
	found, err := repo.GetByQRToken(context.Background(), qrToken)
	if err != nil {
		t.Fatal("Failed to get participant by QR token:", err)
	}

	if found.QRToken != qrToken {
		t.Errorf("Expected QR token %s, got %s", qrToken, found.QRToken)
	}
	if found.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", found.Name)
	}

	t.Log("✅ Retrieved participant by QR token successfully")
}

func TestParticipantRepository_CountByEventID(t *testing.T) {
	repo, eventID := setupTestParticipantRepo(t)
	defer cleanupTestEvent(t, repo, eventID)

	// Create 5 participants
	for i := 0; i < 5; i++ {
		p := &domain.Participant{
			EventID: eventID,
			Name:    "User",
			Email:   "user@example.com",
			Phone:   "08123456789",
			QRToken: uuid.New().String(),
		}
		repo.Create(context.Background(), p)
	}

	count, err := repo.CountByEventID(context.Background(), eventID)
	if err != nil {
		t.Fatal("Failed to count participants:", err)
	}

	if count != 5 {
		t.Errorf("Expected count 5, got %d", count)
	}

	t.Log("✅ Count participants working correctly")
}

func TestParticipantRepository_UpdateCheckIn(t *testing.T) {
	repo, eventID := setupTestParticipantRepo(t)
	defer cleanupTestEvent(t, repo, eventID)

	participant := &domain.Participant{
		EventID: eventID,
		Name:    "John Doe",
		Email:   "john@example.com",
		Phone:   "08123456789",
		QRToken: uuid.New().String(),
	}

	repo.Create(context.Background(), participant)

	// Update check-in
	err := repo.UpdateCheckIn(context.Background(), participant.ID)
	if err != nil {
		t.Fatal("Failed to update check-in:", err)
	}

	// Verify check-in
	found, _ := repo.GetByQRToken(context.Background(), participant.QRToken)
	if !found.CheckedIn {
		t.Error("Expected checked_in to be true")
	}
	if found.CheckedInAt == nil {
		t.Error("Expected checked_in_at to be set")
	}

	t.Log("✅ Check-in updated successfully")
}

func TestParticipantRepository_CountCheckedInByEventID(t *testing.T) {
	repo, eventID := setupTestParticipantRepo(t)
	defer cleanupTestEvent(t, repo, eventID)

	// Create 5 participants, check-in 2 of them
	for i := 0; i < 5; i++ {
		p := &domain.Participant{
			EventID: eventID,
			Name:    "User",
			Email:   "user@example.com",
			Phone:   "08123456789",
			QRToken: uuid.New().String(),
		}
		repo.Create(context.Background(), p)

		if i < 2 {
			repo.UpdateCheckIn(context.Background(), p.ID)
		}
	}

	count, err := repo.CountCheckedInByEventID(context.Background(), eventID)
	if err != nil {
		t.Fatal("Failed to count checked-in participants:", err)
	}

	if count != 2 {
		t.Errorf("Expected 2 checked-in, got %d", count)
	}

	t.Log("✅ Count checked-in participants working correctly")
}

func TestParticipantRepository_DeleteByEventID(t *testing.T) {
	repo, eventID := setupTestParticipantRepo(t)
	defer cleanupTestEvent(t, repo, eventID)

	// Create 3 participants
	for i := 0; i < 3; i++ {
		p := &domain.Participant{
			EventID: eventID,
			Name:    "User",
			Email:   "user@example.com",
			Phone:   "08123456789",
			QRToken: uuid.New().String(),
		}
		repo.Create(context.Background(), p)
	}

	// Delete all
	err := repo.DeleteByEventID(context.Background(), eventID)
	if err != nil {
		t.Fatal("Failed to delete participants:", err)
	}

	// Verify deleted
	count, _ := repo.CountByEventID(context.Background(), eventID)
	if count != 0 {
		t.Errorf("Expected 0 participants after delete, got %d", count)
	}

	t.Log("✅ Participants deleted by event ID successfully")
}

func BenchmarkBulkCreate(b *testing.B) {
	cfg, _ := config.LoadConfig("../../../.env")
	db, _ := database.InitDB(cfg)
	repo := &participantRepository{db: db}

	// Create test event
	eventID := uuid.New().String()
	db.Exec(`
		INSERT INTO events (id, organizer_id, name, slug, date, venue, participant_count, total_price, payment_status, scanner_pin)
		VALUES (?, 1, 'Bench Event', ?, NOW(), 'Test Venue', 1000, 4500000, 'pending', '1234')
	`, eventID, "bench-"+eventID)

	defer func() {
		repo.DeleteByEventID(context.Background(), eventID)
		db.Exec("DELETE FROM events WHERE id = ?", eventID)
		db.Close()
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create 100 participants
		participants := make([]*domain.Participant, 100)
		for j := 0; j < 100; j++ {
			participants[j] = &domain.Participant{
				EventID: eventID,
				Name:    "User",
				Email:   "user@example.com",
				Phone:   "08123456789",
				QRToken: uuid.New().String(),
			}
		}

		repo.BulkCreate(context.Background(), participants)
		repo.DeleteByEventID(context.Background(), eventID) // Cleanup for next iteration
	}
}
