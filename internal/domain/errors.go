package domain

import "errors"

var (
	// Participant errors
	ErrParticipantNameRequired  = errors.New("nama participant harus diisi")
	ErrParticipantEmailRequired = errors.New("email participant harus diisi")
	ErrParticipantPhoneRequired = errors.New("nomor HP participant harus diisi")
	ErrInvalidCSVFormat         = errors.New("format CSV tidak valid")
	ErrEmptyCSV                 = errors.New("file CSV kosong")
)
