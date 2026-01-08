package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/fzndps/eventcheck/internal/domain"
	"github.com/fzndps/eventcheck/internal/domain/repository"
	"github.com/fzndps/eventcheck/internal/infrastructure/email"
	"github.com/fzndps/eventcheck/internal/infrastructure/qrcode"
)

type QREmailUsecae struct {
	eventRepo       repository.EventRepository
	participantRepo repository.ParticipantRepository
	qrGenerator     *qrcode.Generator
	emailService    *email.EmailService
}

func NewQREmailUsecase(
	eventRepo repository.EventRepository,
	participantRepo repository.ParticipantRepository,
	qrGenerator *qrcode.Generator,
	emailService *email.EmailService,
) *QREmailUsecae {
	return &QREmailUsecae{
		eventRepo:       eventRepo,
		participantRepo: participantRepo,
		qrGenerator:     qrGenerator,
		emailService:    emailService,
	}
}

// SendQRCodes untuk mengirim qr code ke semua participants yang belum menerima
func (u *QREmailUsecae) SendQRCodes(ctx context.Context, organizerID int64, eventID string) (*domain.SendQRCodesResponse, error) {
	// 1. Check authorization
	isOwned, err := u.eventRepo.IsOwnedBy(ctx, eventID, organizerID)
	if err != nil {
		return nil, err
	}
	if !isOwned {
		return nil, domain.ErrUnauthorizedAccess
	}

	// 2. Get event details
	event, err := u.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	// 3. Get participants yang belum dikirim QR
	participants, err := u.participantRepo.GetPendingQR(ctx, eventID)
	if err != nil {
		return nil, err
	}

	if len(participants) == 0 {
		return &domain.SendQRCodesResponse{
			TotalParticipants: 0,
			EmailsSent:        0,
			EmailsFailed:      0,
			FailedEmails:      []string{},
		}, nil
	}

	// 4. Generate & send emails
	var emailsSent, emailsFailed int
	var failedEmails []string

	for _, participant := range participants {
		// Generate QR code as PNG bytes (for CID embedding)
		qrBytes, err := u.qrGenerator.GenerateQRCode(participant.QRToken, 256)
		if err != nil {
			log.Printf("Failed to generate QR for participant %d: %v", participant.ID, err)
			emailsFailed++
			failedEmails = append(failedEmails, participant.Email)
			continue
		}

		// Build email HTML (use CID method)
		emailBody := email.BuildQRCodeEmail(participant, event, "", true)

		// Send email with embedded image (CID method)
		subject := fmt.Sprintf("Your QR Code for %s", event.Name)
		err = u.emailService.SendEmailWithEmbeddedImage(
			participant.Email,
			subject,
			emailBody,
			qrBytes,
		)

		if err != nil {
			log.Printf("Failed to send email to %s: %v", participant.Email, err)
			emailsFailed++
			failedEmails = append(failedEmails, participant.Email)
			continue
		}

		// Mark QR as sent
		err = u.participantRepo.MarkQRSent(ctx, participant.ID)
		if err != nil {
			log.Printf("Failed to mark QR sent for participant %d: %v", participant.ID, err)
		}

		emailsSent++
		log.Printf("QR code sent to %s (%s)", participant.Name, participant.Email)
	}

	res := &domain.SendQRCodesResponse{
		TotalParticipants: len(participants),
		EmailsSent:        emailsSent,
		EmailsFailed:      emailsFailed,
		FailedEmails:      failedEmails,
	}

	return res, nil
}

func (u QREmailUsecae) ResendQRCode(ctx context.Context, organizerID int64, eventID string, participanID int64) error {
	// cek authorization untuk event
	isOwned, err := u.eventRepo.IsOwnedBy(ctx, eventID, organizerID)
	if err != nil {
		return fmt.Errorf("This event does not belong to the organizer :%v", err)
	}

	if !isOwned {
		return domain.ErrUnauthorizedAccess
	}

	// get participan
	participant, err := u.participantRepo.GetByID(ctx, participanID)
	if err != nil {
		return fmt.Errorf("failed to get participant: %v", err)
	}

	if participant.EventID != eventID {
		return domain.ErrUnauthorizedAccess
	}

	// get detail event
	event, err := u.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return fmt.Errorf("failed to get detail event: %v", err)
	}

	// generate qr code
	qrByte, err := u.qrGenerator.GenerateQRCode(participant.QRToken, 256)
	if err != nil {
		return fmt.Errorf("failed to generate QR code: %w", err)
	}

	// build email HTML
	emailBody := email.BuildQRCodeEmail(participant, event, "", true)

	// send email
	subject := fmt.Sprintf("Your QR code for %s (Resent)", event.Name)
	err = u.emailService.SendEmailWithEmbeddedImage(
		participant.Email,
		subject,
		emailBody,
		qrByte,
	)

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if !participant.QRSent {
		err := u.participantRepo.MarkQRSent(ctx, participant.ID)
		if err != nil {
			return fmt.Errorf("failed to mark QR sent for participant %d: %v", participant.ID, err)
		}
	}

	log.Printf("QR code resent to %s (%s)", participant.Name, participant.Email)

	return nil
}

func (u *QREmailUsecae) SendTestEmail(ctx context.Context, toEmail, recipientName string) error {
	emailBody := email.BuildTestEmail(recipientName)

	emailData := &email.EmailData{
		To:      toEmail,
		Subject: "Test email from EventCheck.In",
		Body:    emailBody,
		IsHTML:  true,
	}

	return u.emailService.SendEmail(emailData)
}
