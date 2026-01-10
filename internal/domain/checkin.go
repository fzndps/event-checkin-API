package domain

import "time"

// VerifyPINRequest adalah request untuk verify scanner PIN
type VerifyPINRequest struct {
	// EventSlug  string `json:"event_slug" binding:"required"`
	ScannerPIN string `json:"scanner_pin" binding:"required,len=4"`
}

// VerifyPINResponse adalah response setelah PIN verified
type VerifyPINResponse struct {
	Valid     bool   `json:"valid"`
	EventID   string `json:"event_id,omitempty"`
	EventName string `json:"event_name,omitempty"`
	Message   string `json:"message"`
}

// CheckInRequest adalah request untuk check-in participant
type CheckInRequest struct {
	QRToken string `json:"qr_token" binding:"required"`
}

// CheckInResponse adalah response setelah check-in
type CheckInResponse struct {
	Success          bool         `json:"success"`
	Message          string       `json:"message"`
	Participant      *Participant `json:"participant,omitempty"`
	AlreadyCheckedIn bool         `json:"already_checked_in"`
	CheckedInAt      *time.Time   `json:"checked_in_at,omitempty"`
}

// EventStatsResponse adalah response untuk dashboard stats
type EventStatsResponse struct {
	EventID           string         `json:"event_id"`
	EventName         string         `json:"event_name"`
	TotalParticipants int            `json:"total_participants"`
	CheckedInCount    int            `json:"checked_in_count"`
	NotCheckedInCount int            `json:"not_checked_in_count"`
	CheckInPercentage float64        `json:"checkin_percentage"`
	LastCheckInAt     *time.Time     `json:"last_checkin_at,omitempty"`
	RecentCheckIns    []*Participant `json:"recent_checkins,omitempty"` // Last 10
}

// ScanPageData adalah data untuk scan page (web view)
type ScanPageData struct {
	EventSlug string
	EventName string
	EventDate time.Time
	Venue     string
}
