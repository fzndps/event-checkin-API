package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/fzndps/eventcheck/internal/domain"
	"github.com/fzndps/eventcheck/internal/domain/repository"
)

type organizerRepositoryImpl struct {
	db *sql.DB
}

func NewOrganizerRepositoryImpl(db *sql.DB) repository.OrganizerRepository {
	return &organizerRepositoryImpl{
		db: db,
	}
}

func (r *organizerRepositoryImpl) Create(ctx context.Context, organizer *domain.Organizer) error {
	query := "INSERT INTO organizers (email, name, password_hash, created_at) VALUE (?, ?, ?, NOW())"

	result, err := r.db.ExecContext(ctx, query, organizer.Email, organizer.Name, organizer.PasswordHash)

	if err != nil {
		if isDuplicateKeyError(err) {
			return fmt.Errorf("email already exists: %v", err)
		}
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id organizer: %w", err)
	}

	organizer.ID = int64(id)

	return nil
}

func (r *organizerRepositoryImpl) GetByEmail(ctx context.Context, email string) (*domain.Organizer, error) {
	organizer := &domain.Organizer{}
	query := "SELECT id, email, name, password_hash, created_at FROM organizers WHERE email = ?"

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&organizer.ID,
		&organizer.Email,
		&organizer.Name,
		&organizer.PasswordHash,
		&organizer.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get organizer: %w", err)
	}

	return organizer, nil
}

func (r *organizerRepositoryImpl) GetByID(ctx context.Context, id int64) (*domain.Organizer, error) {
	organizer := &domain.Organizer{}
	query := "SELECT id, email, name, password_hash, created_at FROM organizers WHERE id = ?"

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&organizer.ID,
		&organizer.Email,
		&organizer.Name,
		&organizer.PasswordHash,
		&organizer.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return organizer, nil
}
