package email

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"mime/quotedprintable"
	"net/smtp"
	"net/textproto"
	"strings"

	"github.com/fzndps/eventcheck/config"
)

// EmailService adalah service untuk mengirim email
type EmailService struct {
	config *config.SMTPConfig
}

// NewEmailService constructor untuk email service
func NewEmailService(config *config.SMTPConfig) *EmailService {
	return &EmailService{
		config: config,
	}
}

// EmailData untuk compose email
type EmailData struct {
	To          string
	Subject     string
	Body        string
	Attachments map[string][]byte
	IsHTML      bool
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

	} else {
		// Email without attachments (simple)
		if data.IsHTML {
			buf.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		} else {
			buf.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		}
		buf.WriteString("\r\n")
		buf.WriteString(data.Body)
	}

	if len(data.Attachments) > 0 {
		writer.Close()
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

	// Write attachment dengan base64 encoding
	encoded := s.encodeBase64(content)
	attachPart.Write([]byte(encoded))

	return nil
}

// encodeBase64 encode bytes ke base64 dengan line breaks setiap 76 karakter
func (s *EmailService) encodeBase64(data []byte) string {
	const lineLength = 76

	// Encode to base64
	encoded := make([]byte, len(data)*2)
	n := len(encoded)
	for i := 0; i < len(data); i += 3 {
		end := i + 3
		if end > len(data) {
			end = len(data)
		}

		chunk := data[i:end]
		encodedChunk := []byte(fmt.Sprintf("%x", chunk))
		copy(encoded[i*2:], encodedChunk)
	}
	encoded = encoded[:n]

	// Add line breaks
	var result strings.Builder
	for i := 0; i < len(encoded); i += lineLength {
		end := i + lineLength
		if end > len(encoded) {
			end = len(encoded)
		}
		result.Write(encoded[i:end])
		result.WriteString("\r\n")
	}

	return result.String()
}

// SendBulkEmails mengirim email ke multiple recipients
// Return: success count, failed count, errors
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
