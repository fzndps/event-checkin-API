package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/fzndps/eventcheck/internal/domain"
	"github.com/fzndps/eventcheck/internal/domain/repository"
)

type participantRepository struct {
	db *sql.DB
}

func NewParticipantRepository(db *sql.DB) repository.ParticipantRepository {
	return &participantRepository{
		db: db,
	}
}

// Create menyimpan satu participant
func (r *participantRepository) Create(ctx context.Context, participant *domain.Participant) error {
	query := `
		INSERT INTO participants
			(event_id, name, email, phone, qr_token,
			checked_in, checked_in_at, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		participant.EventID,
		participant.Name,
		participant.Email,
		participant.Phone,
		participant.QRToken,
		participant.CheckedIn,
		participant.CheckedInAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create participant: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %v", err)
	}

	participant.ID = id
	return nil
}

// BulkCreate menyimpan banyak participants sekaligus (untuk CSV upload)
// Menggunakan transaction untuk atomicity: all or nothing
func (r *participantRepository) BulkCreate(ctx context.Context, participants []*domain.Participant) error {
	if len(participants) == 0 {
		return nil
	}

	// Mulai transaction menggunakan tx
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback() // Rollback semua jika eerror

	// bulk insert query
	valueStrings := make([]string, 0, len(participants))
	valueArgs := make([]interface{}, 0, len(participants)*7) // 7 kolom

	for _, p := range participants {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, NOW())")
		valueArgs = append(valueArgs,
			p.EventID,
			p.Name,
			p.Email,
			p.Phone,
			p.QRToken,
			false, // checked in default dibuat false
			nil,   // checked in at default dibuat null
		)
	}

	query := fmt.Sprintf(`
		INSERT INTO participants (
		event_id, name, email, phone, qr_token,
		checked_in, checked_in_at, created_at
		) VALUES %s
	`, strings.Join(valueStrings, ","))

	_, err = tx.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return err
	}

	// Commit transaction jika berhasil
	return tx.Commit()
}

// GetByID mencari participant berdasarkan ID
func (r *participantRepository) GetByID(ctx context.Context, id int64) (*domain.Participant, error) {
	query := `
		SELECT 
			id, event_id, name, email, phone, qr_token,
			checked_in, checked_in_at, qr_sent, qr_sent_at, created_at
		FROM participants
		WHERE id = ?
	`

	p := &domain.Participant{}
	var checkinAt, qrSentAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID,
		&p.EventID,
		&p.Name,
		&p.Email,
		&p.Phone,
		&p.QRToken,
		&p.CheckedIn,
		&checkinAt,
		&p.QRSent,
		&qrSentAt,
		&p.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	if checkinAt.Valid {
		p.CheckedInAt = &checkinAt.Time
	}

	if qrSentAt.Valid {
		p.QRSentAt = &qrSentAt.Time
	}

	return p, nil
}

// GetByEventID mencari semua participant di event tertentu
func (r *participantRepository) GetByEventID(ctx context.Context, eventID string) ([]*domain.Participant, error) {
	query := `
		SELECT
			id, event_id, name, email, phone, qr_token,
			checked_in, checked_in_at, created_at
		FROM participants
		WHERE event_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get participant: %v", err)
	}

	defer rows.Close()

	var participants []*domain.Participant
	for rows.Next() {
		participant := &domain.Participant{}
		var checkedInAt sql.NullTime
		err := rows.Scan(
			&participant.ID,
			&participant.EventID,
			&participant.Name,
			&participant.Email,
			&participant.Phone,
			&participant.QRToken,
			&participant.CheckedIn,
			&checkedInAt,
			&participant.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		if checkedInAt.Valid {
			participant.CheckedInAt = &checkedInAt.Time
		}

		participants = append(participants, participant)
	}

	return participants, nil

}

// GetByQRToken mencari participant berdasarkan QR token (untuk check-in)
func (r *participantRepository) GetByQRToken(ctx context.Context, qrToken string) (*domain.Participant, error) {
	query := `
		SELECT
			id, event_id, name, email, phone, qr_token,
			checked_in, checked_in_at, created_at
		FROM participants
		WHERE qr_token = ?
	`

	participant := &domain.Participant{}
	var checkedInAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, qrToken).Scan(
		&participant.ID,
		&participant.EventID,
		&participant.Name,
		&participant.Email,
		&participant.Phone,
		&participant.QRToken,
		&participant.CheckedIn,
		&checkedInAt,
		&participant.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to get participant by token %s: %v", qrToken, err)
		}
		return nil, err
	}

	if checkedInAt.Valid {
		participant.CheckedInAt = &checkedInAt.Time
	}

	return participant, nil

}

// CountByEventID menghitung jumlah participant di event
func (r *participantRepository) CountByEventID(ctx context.Context, eventID string) (int, error) {
	query := `SELECT COUNT(*) FROM participants WHERE event_id = ?`

	var count int
	err := r.db.QueryRowContext(ctx, query, eventID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// CountCheckedInByEventID menghitung jumlah participant yang sudah check-in
func (r *participantRepository) CountCheckedInByEventID(ctx context.Context, eventID string) (int, error) {
	query := `SELECT COUNT(*) FROM participants WHERE event_id = ? AND checked_in = TRUE`

	var count int
	err := r.db.QueryRowContext(ctx, query, eventID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// UpdateCheckIn mengupdate status check-in participant
func (r *participantRepository) UpdateCheckIn(ctx context.Context, participantID int64) error {
	query := `
		UPDATE participants
		SET checked_in = TRUE, checked_in_at = NOW()
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query, participantID)
	if err != nil {
		return fmt.Errorf("failed to update checkin: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("update not found: %v", err)
	}

	return nil
}

func (r *participantRepository) MarkQRSent(ctx context.Context, participantID int64) error {
	query := `
		UPDATE participants
		SET qr_sent = true, qr_sent_at = NOW()
		WHERE id  = ?
	`

	result, err := r.db.ExecContext(ctx, query, participantID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *participantRepository) GetPendingQR(ctx context.Context, eventID string) ([]*domain.Participant, error) {
	query := `
		SELECT 
			id, event_id, name, email, phone, qr_token,
			checked_in, checked_in_at, qr_sent, qr_sent_at, created_at
		FROM participants
		WHERE event_id = ? AND qr_sent = FALSE
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []*domain.Participant
	for rows.Next() {
		p := &domain.Participant{}
		var checkedInAt, qrSentAt sql.NullTime

		err := rows.Scan(
			&p.ID,
			&p.EventID,
			&p.Name,
			&p.Email,
			&p.Phone,
			&p.QRToken,
			&p.CheckedIn,
			&checkedInAt,
			&p.QRSent,
			&qrSentAt,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if checkedInAt.Valid {
			p.CheckedInAt = &checkedInAt.Time
		}
		if qrSentAt.Valid {
			p.QRSentAt = &qrSentAt.Time
		}

		participants = append(participants, p)
	}

	return participants, nil
}

// DeleteByEventID menghapus semua participant di event (cascade delete)
func (r *participantRepository) DeleteByEventID(ctx context.Context, eventID string) error {
	query := `DELETE FROM participants WHERE event_id = ?`
	_, err := r.db.ExecContext(ctx, query, eventID)

	return err
}
