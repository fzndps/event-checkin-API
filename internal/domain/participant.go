package domain

import "time"

// Participant entity peserta event
type Participant struct {
	ID          int64      `json:"id"`
	EventID     string     `json:"event_id"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	Phone       string     `json:"phone"`
	QRToken     string     `json:"qr_token"`
	CheckedIn   bool       `json:"checked_in"`
	CheckedInAt *time.Time `json:"checked_in_at"`
	CreatedAt   time.Time  `json:"created_at"`

	// QR Code URL (generated, tidak disimpan di DB)
	QRCodeURL string `json:"qr_code_url,omitempty"`
}

// ParticipantCSVRow struktur untuk parse CSV
type ParticipantCSVRow struct {
	Name  string `csv:"name"`  // Kolom "name" di CSV
	Email string `csv:"email"` // Kolom "email" di CSV
	Phone string `csv:"phone"` // Kolom "phone" di CSV
}

// UploadParticipantsRequest request untuk upload CSV
type UploadParticipantsRequest struct {
	EventID string
}

// UploadParticipantsResponse response setelah upload CSV
type UploadParticipantsResponse struct {
	Success       int      `json:"success"`
	Failed        int      `json:"failed"`
	FailedReasons []string `json:"failed_reasons"`
}

// Validate melakukan validasi data participant
func (p *Participant) Validate() error {
	if p.Name == "" {
		return ErrParticipantNameRequired
	}
	if p.Email == "" {
		return ErrParticipantEmailRequired
	}
	if p.Phone == "" {
		return ErrParticipantPhoneRequired
	}

	return nil
}

// IsCheckedIn return true jika participant sudah check-in
func (p *Participant) IsCheckedIn() bool {
	return p.CheckedIn && p.CheckedInAt != nil
}
