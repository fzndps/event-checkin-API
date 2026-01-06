package email

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"mime/quotedprintable"
	"net/smtp"
	"net/textproto"

	"github.com/fzndps/eventcheck/config"
)

// SMTPConfig adalah konfigurasi untuk SMTP email
// type SMTPConfig struct {
// 	Host     string
// 	Port     string
// 	Username string
// 	Password string
// 	From     string
// }

// EmailService adalah service untuk mengirim email
type EmailService struct {
	config *config.SMTPConfig
}

// NewEmailService constructor
func NewEmailService(config *config.SMTPConfig) *EmailService {
	return &EmailService{
		config: config,
	}
}

// EmailData adalah data untuk compose email
type EmailData struct {
	To          string            // Recipient email
	Subject     string            // Email subject
	Body        string            // Email body (HTML or plain text)
	Attachments map[string][]byte // Attachments: filename -> content
	IsHTML      bool              // True jika body adalah HTML
}

// SendEmail mengirim email dengan optional attachments
func (s *EmailService) SendEmail(data *EmailData) error {
	// Setup authentication
	auth := smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPHost)

	// Compose email message
	message, err := s.composeMessage(data)
	if err != nil {
		return fmt.Errorf("failed to compose message: %w", err)
	}

	// Send email
	addr := fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort)
	err = smtp.SendMail(addr, auth, s.config.SMTPFrom, []string{data.To}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// composeMessage membuat email message dengan format MIME
func (s *EmailService) composeMessage(data *EmailData) ([]byte, error) {
	var buf bytes.Buffer

	// Headers
	buf.WriteString(fmt.Sprintf("From: %s\r\n", s.config.SMTPFrom))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", data.To))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", data.Subject))
	buf.WriteString("MIME-Version: 1.0\r\n")

	// Create multipart writer
	writer := multipart.NewWriter(&buf)
	boundary := writer.Boundary()

	if len(data.Attachments) > 0 {
		// Email with attachments
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", boundary))
		buf.WriteString("\r\n")

		// Write body part
		if err := s.writeBodyPart(writer, data); err != nil {
			return nil, err
		}

		// Write attachment parts
		for filename, content := range data.Attachments {
			if err := s.writeAttachmentPart(writer, filename, content); err != nil {
				return nil, err
			}
		}

		writer.Close()

	} else {
		// Email without attachments (simple)
		if data.IsHTML {
			buf.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		} else {
			buf.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		}
		buf.WriteString("Content-Transfer-Encoding: quoted-printable\r\n")
		buf.WriteString("\r\n")

		// Write body dengan quoted-printable encoding
		qpWriter := quotedprintable.NewWriter(&buf)
		qpWriter.Write([]byte(data.Body))
		qpWriter.Close()
	}

	return buf.Bytes(), nil
}

// writeBodyPart menulis body part dari email
func (s *EmailService) writeBodyPart(writer *multipart.Writer, data *EmailData) error {
	// Create body part headers
	bodyHeader := make(textproto.MIMEHeader)
	if data.IsHTML {
		bodyHeader.Set("Content-Type", "text/html; charset=UTF-8")
	} else {
		bodyHeader.Set("Content-Type", "text/plain; charset=UTF-8")
	}
	bodyHeader.Set("Content-Transfer-Encoding", "quoted-printable")

	// Create part
	bodyPart, err := writer.CreatePart(bodyHeader)
	if err != nil {
		return err
	}

	// Write body dengan quoted-printable encoding
	qpWriter := quotedprintable.NewWriter(bodyPart)
	qpWriter.Write([]byte(data.Body))
	qpWriter.Close()

	return nil
}

// writeAttachmentPart menulis attachment part dari email
func (s *EmailService) writeAttachmentPart(writer *multipart.Writer, filename string, content []byte) error {
	// Create attachment headers
	attachHeader := make(textproto.MIMEHeader)
	attachHeader.Set("Content-Type", "application/octet-stream")
	attachHeader.Set("Content-Transfer-Encoding", "base64")
	attachHeader.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	// Create part
	attachPart, err := writer.CreatePart(attachHeader)
	if err != nil {
		return err
	}

	// Write attachment dengan base64 encoding (proper format)
	encoded := base64.StdEncoding.EncodeToString(content)

	// Add line breaks every 76 characters (RFC 2045)
	for i := 0; i < len(encoded); i += 76 {
		end := i + 76
		if end > len(encoded) {
			end = len(encoded)
		}
		attachPart.Write([]byte(encoded[i:end]))
		attachPart.Write([]byte("\r\n"))
	}

	return nil
}

// SendBulkEmails mengirim email ke multiple recipients
func (s *EmailService) SendBulkEmails(emails []*EmailData) (int, int, []error) {
	var successCount, failedCount int
	var errors []error

	for i, emailData := range emails {
		err := s.SendEmail(emailData)
		if err != nil {
			failedCount++
			errors = append(errors, fmt.Errorf("email %d (%s): %w", i+1, emailData.To, err))
		} else {
			successCount++
		}
	}

	return successCount, failedCount, errors
}
