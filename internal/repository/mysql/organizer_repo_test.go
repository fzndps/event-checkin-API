package mysql

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/fzndps/eventcheck/config"
	"github.com/fzndps/eventcheck/internal/domain"
	"github.com/fzndps/eventcheck/internal/infrastructure/database"
)

func setupTestDB(t *testing.T) *sql.DB {
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal("Failed to load config:", err)
	}

	db, err := database.InitDB(cfg)
	if err != nil {
		t.Fatal("Failed to connect to database:", err)
	}

	return db
}

func TestCreateOrganizer(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewOrganizerRepositoryImpl(db)

	organizer := &domain.Organizer{
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		Name:         "Test User",
		CreatedAt:    time.Now(),
	}

	err := repo.Create(context.Background(), organizer)
	if err != nil {
		t.Fatal("Failed to create organizer:", err)
	}

	if organizer.ID == 0 {
		t.Fatal("Organizer ID should be set")
	}

	t.Log("✅ Organizer created with ID:", organizer.ID)

	// Cleanup
	db.Exec("DELETE FROM organizers WHERE email = ?", organizer.Email)
}

func TestGetByEmail(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewOrganizerRepositoryImpl(db)

	// Create test organizer
	email := "test-get@example.com"
	organizer := &domain.Organizer{
		Email:        email,
		PasswordHash: "hashedpassword",
		Name:         "Test User",
		CreatedAt:    time.Now(),
	}
	repo.Create(context.Background(), organizer)

	// Test GetByEmail
	found, err := repo.GetByEmail(context.Background(), email)
	if err != nil {
		t.Fatal("Failed to get organizer:", err)
	}

	if found.Email != email {
		t.Fatalf("Expected email %s, got %s", email, found.Email)
	}

	t.Log("✅ Organizer found by email")

	// Cleanup
	db.Exec("DELETE FROM organizers WHERE email = ?", email)
}
