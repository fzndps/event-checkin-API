package domain

import "errors"

var (
	// Participant errors
	ErrParticipantNameRequired  = errors.New("participant name is required")
	ErrParticipantEmailRequired = errors.New("participant email is required")
	ErrParticipantPhoneRequired = errors.New("participant phone number is required")
	ErrInvalidCSVFormat         = errors.New("invalid CSV format")
	ErrEmptyCSV                 = errors.New("CSV file is empty")

	// Event errors
	ErrEventNotFound      = errors.New("event not found")
	ErrInvalidEventDate   = errors.New("invalid event date (cannot be in the past)")
	ErrSlugAlreadyExists  = errors.New("event slug already in use")
	ErrUnauthorizedAccess = errors.New("you do not have access to this event")
)
