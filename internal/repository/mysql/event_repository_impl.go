package mysql

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/fzndps/eventcheck/internal/domain"
	"github.com/fzndps/eventcheck/internal/domain/repository"
)

type eventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) repository.EventRepository {
	return &eventRepository{
		db: db,
	}
}

// Create menyimpan event baru ke database
func (r *eventRepository) Create(ctx context.Context, event *domain.Event) error {
	query := `INSERT INTO events (
			id, organizer_id, name, slug, date, venue, 
			participant_count, total_price, payment_status, 
			payment_proof_url, scanner_pin, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`

	_, err := r.db.ExecContext(ctx, query,
		event.ID,
		event.OrganizerID,
		event.Name,
		event.Slug,
		event.Date,
		event.Venue,
		event.ParticipantCount,
		event.TotalPrice,
		event.PaymentStatus,
		event.PaymentProofURL,
		event.ScannerPIN,
	)

	if err != nil {
		if isDuplicateKeyError(err) {
			return domain.ErrSlugAlreadyExists
		}

		return err
	}

	return nil
}

// GetByID mencari event berdasarkan ID
func (r *eventRepository) GetByID(ctx context.Context, id string) (*domain.Event, error) {
	query := `
		SELECT id, organizer_id, name, slug, date, venue, 
		participant_count, total_price, payment_status, 
		payment_proof_url, scanner_pin, created_at
		FROM events WHERE id = ?
	`

	event := &domain.Event{}
	var paymentProofURL sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&event.ID,
		&event.OrganizerID,
		&event.Name,
		&event.Slug,
		&event.Date,
		&event.Venue,
		&event.ParticipantCount,
		&event.TotalPrice,
		&event.PaymentStatus,
		&paymentProofURL,
		&event.ScannerPIN,
		&event.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrEventNotFound
		}
		log.Print("error get by id:", err)
		return nil, err
	}

	if paymentProofURL.Valid {
		event.PaymentProofURL = paymentProofURL.String
	}

	return event, nil
}

// GetBySlug mencari event berdasarkan slug
func (r *eventRepository) GetBySlug(ctx context.Context, slug string) (*domain.Event, error) {
	query := `
		SELECT id, organizer_id, name, slug, date, venue, 
		participant_count, total_price, payment_status, 
		payment_proof_url, scanner_pin, created_at
		FROM events WHERE slug = ?
	`

	event := &domain.Event{}
	var paymentProofURL sql.NullString

	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&event.ID,
		&event.OrganizerID,
		&event.Name,
		&event.Slug,
		&event.Date,
		&event.Venue,
		&event.ParticipantCount,
		&event.TotalPrice,
		&event.PaymentStatus,
		&paymentProofURL,
		&event.ScannerPIN,
		&event.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrEventNotFound
		}

		return nil, err
	}

	if paymentProofURL.Valid {
		event.PaymentProofURL = paymentProofURL.String
	}

	return event, nil
}

// GetByOrganizerID mencari semua event milik organizer dengan pagination
// offset = (page - 1) * limit
func (r *eventRepository) GetByOrganizerID(ctx context.Context, organizerID int64, limit, offset int) ([]*domain.Event, int, error) {
	query := `
		SELECT
			id, organizer_id, name, slug, date, venue, 
			participant_count, total_price, payment_status, 
			payment_proof_url, scanner_pin, created_at
		FROM events 
		WHERE organizer_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, organizerID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var events []*domain.Event
	for rows.Next() {
		event := &domain.Event{}
		var paymentProofURL sql.NullString

		err := rows.Scan(
			&event.ID,
			&event.OrganizerID,
			&event.Name,
			&event.Slug,
			&event.Date,
			&event.Venue,
			&event.ParticipantCount,
			&event.TotalPrice,
			&event.PaymentStatus,
			&paymentProofURL,
			&event.ScannerPIN,
			&event.CreatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		if paymentProofURL.Valid {
			event.PaymentProofURL = paymentProofURL.String
		}

		events = append(events, event)
	}

	countQuery := `SELECT COUNT(*) FROM events WHERE organizer_id = ?`
	var total int
	err = r.db.QueryRowContext(ctx, countQuery, organizerID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

// Update mengupdate data event
func (r *eventRepository) Update(ctx context.Context, event *domain.Event) error {
	query := `
		UPDATE events SET
			name = ?,
			slug = ?,
			date = ?,
			venue = ?,
			participant_count = ?,
			total_price = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		event.Name,
		event.Slug,
		event.Date,
		event.Venue,
		event.ParticipantCount,
		event.TotalPrice,
		event.ID,
	)

	if err != nil {
		if isDuplicateKeyError(err) {
			return domain.ErrSlugAlreadyExists
		}

		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrEventNotFound
	}

	return nil
}

// Delete menghapus event (dan cascade delete participants)
func (r *eventRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM events WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrEventNotFound
	}

	return nil
}

// IsOwnedBy mengecek apakah event dimiliki oleh organizer
func (r *eventRepository) IsOwnedBy(ctx context.Context, eventID string, organizerID int64) (bool, error) {
	query := `SELECT COUNT(*) FROM events WHERE id = ? AND organizer_id = ?`

	var count int
	err := r.db.QueryRowContext(ctx, query, eventID, organizerID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
