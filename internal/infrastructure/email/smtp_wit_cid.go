package email

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"mime/quotedprintable"
	"net/smtp"
	"net/textproto"
)

// SendEmailWithEmbeddedImage mengirim email dengan QR code sebagai embedded image (CID)
// Ini lebih reliable daripada base64 inline untuk compatibility dengan email clients
func (s *EmailService) SendEmailWithEmbeddedImage(to, subject, htmlBody string, qrCodeImage []byte) error {
	var buf bytes.Buffer

	// Headers
	buf.WriteString(fmt.Sprintf("From: %s\r\n", s.config.SMTPFrom))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", to))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	buf.WriteString("MIME-Version: 1.0\r\n")

	// Create multipart/related writer (for embedded images)
	writer := multipart.NewWriter(&buf)
	boundary := writer.Boundary()

	buf.WriteString(fmt.Sprintf("Content-Type: multipart/related; boundary=%s\r\n", boundary))
	buf.WriteString("\r\n")

	// Part 1: HTML body
	htmlHeader := make(textproto.MIMEHeader)
	htmlHeader.Set("Content-Type", "text/html; charset=UTF-8")
	htmlHeader.Set("Content-Transfer-Encoding", "quoted-printable")

	htmlPart, err := writer.CreatePart(htmlHeader)
	if err != nil {
		return fmt.Errorf("failed to create HTML part: %w", err)
	}

	qpWriter := quotedprintable.NewWriter(htmlPart)
	qpWriter.Write([]byte(htmlBody))
	qpWriter.Close()

	// Part 2: Embedded QR code image
	imageHeader := make(textproto.MIMEHeader)
	imageHeader.Set("Content-Type", "image/png")
	imageHeader.Set("Content-Transfer-Encoding", "base64")
	imageHeader.Set("Content-ID", "<qrcode>") // CID reference
	imageHeader.Set("Content-Disposition", "inline; filename=\"qrcode.png\"")

	imagePart, err := writer.CreatePart(imageHeader)
	if err != nil {
		return fmt.Errorf("failed to create image part: %w", err)
	}

	// Encode image to base64
	encoded := base64.StdEncoding.EncodeToString(qrCodeImage)

	// Write with line breaks every 76 characters (RFC 2045)
	for i := 0; i < len(encoded); i += 76 {
		end := i + 76
		if end > len(encoded) {
			end = len(encoded)
		}
		imagePart.Write([]byte(encoded[i:end]))
		imagePart.Write([]byte("\r\n"))
	}

	writer.Close()

	// Send email
	auth := smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPHost)
	addr := fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort)

	err = smtp.SendMail(addr, auth, s.config.SMTPFrom, []string{to}, buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
