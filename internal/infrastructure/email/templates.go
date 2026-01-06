package email

import (
	"bytes"
	"html/template"
	"time"

	"github.com/fzndps/eventcheck/internal/domain"
)

// QREmailTemplate adalah data untuk email template QR code
type QREmailTemplate struct {
	ParticipantName string
	EventName       string
	EventDate       string
	EventVenue      string
	QRCodeBase64    string // Base64 encoded QR code (for inline)
	UseCID          bool   // Use CID instead of base64 inline
}

// BuildQRCodeEmail membuat HTML email dengan QR code
// useCID=true untuk embedded image (lebih compatible)
// useCID=false untuk base64 inline
func BuildQRCodeEmail(participant *domain.Participant, event *domain.Event, qrCodeBase64 string, useCID bool) string {
	// Format date
	eventDate := event.Date.Format("Monday, 02 January 2006 at 15:04")

	data := QREmailTemplate{
		ParticipantName: participant.Name,
		EventName:       event.Name,
		EventDate:       eventDate,
		EventVenue:      event.Venue,
		QRCodeBase64:    qrCodeBase64,
		UseCID:          useCID,
	}

	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Your Event QR Code</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 30px;
            text-align: center;
            border-radius: 10px 10px 0 0;
        }
        .content {
            background: #f9f9f9;
            padding: 30px;
            border: 1px solid #ddd;
        }
        .qr-container {
            text-align: center;
            margin: 30px 0;
            padding: 20px;
            background: white;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .qr-code {
            max-width: 256px;
            height: auto;
            margin: 20px auto;
            display: block;
        }
        .event-details {
            background: white;
            padding: 20px;
            border-radius: 10px;
            margin: 20px 0;
        }
        .detail-row {
            margin: 10px 0;
            padding: 10px;
            border-left: 4px solid #667eea;
            background: #f5f5f5;
        }
        .detail-label {
            font-weight: bold;
            color: #667eea;
        }
        .instructions {
            background: #fff3cd;
            border: 1px solid #ffc107;
            padding: 15px;
            border-radius: 5px;
            margin: 20px 0;
        }
        .footer {
            text-align: center;
            color: #666;
            font-size: 12px;
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #ddd;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>🎉 Your Event QR Code</h1>
        <p>Welcome to {{.EventName}}</p>
    </div>
    
    <div class="content">
        <p>Hi <strong>{{.ParticipantName}}</strong>,</p>
        
        <p>Thank you for registering! Here's your unique QR code for check-in at the event.</p>
        
        <div class="qr-container">
            <h3>Your QR Code</h3>
            {{if .UseCID}}
            <img src="cid:qrcode" alt="QR Code" class="qr-code">
            {{else}}
            <img src="{{.QRCodeBase64}}" alt="QR Code" class="qr-code">
            {{end}}
            <p style="color: #666; font-size: 14px;">
                Please show this QR code at the registration desk
            </p>
        </div>
        
        <div class="event-details">
            <h3>📅 Event Details</h3>
            
            <div class="detail-row">
                <div class="detail-label">Event Name:</div>
                <div>{{.EventName}}</div>
            </div>
            
            <div class="detail-row">
                <div class="detail-label">Date & Time:</div>
                <div>{{.EventDate}}</div>
            </div>
            
            <div class="detail-row">
                <div class="detail-label">Venue:</div>
                <div>{{.EventVenue}}</div>
            </div>
        </div>
        
        <div class="instructions">
            <h4>📱 How to Check-In:</h4>
            <ol>
                <li>Save this email or take a screenshot of the QR code</li>
                <li>Arrive at the venue on time</li>
                <li>Show your QR code to the registration staff</li>
                <li>Your attendance will be recorded instantly</li>
            </ol>
        </div>
        
        <p><strong>Important Notes:</strong></p>
        <ul>
            <li>This QR code is unique to you - do not share it with others</li>
            <li>You can print this email or show it on your phone</li>
            <li>If you have any issues, please contact the organizer</li>
        </ul>
        
        <p>We look forward to seeing you at the event!</p>
        
        <p>Best regards,<br>
        <strong>EventCheck.in Team</strong></p>
    </div>
    
    <div class="footer">
        <p>This is an automated email from EventCheck.in</p>
        <p>© {{.Year}} EventCheck.in. All rights reserved.</p>
    </div>
</body>
</html>
`

	// Parse and execute template
	t := template.Must(template.New("email").Parse(tmpl))
	var buf bytes.Buffer

	// Add year to template data
	type templateData struct {
		QREmailTemplate
		Year int
	}

	fullData := templateData{
		QREmailTemplate: data,
		Year:            time.Now().Year(),
	}

	t.Execute(&buf, fullData)

	return buf.String()
}

// BuildPlainTextEmail membuat plain text email
func BuildPlainTextEmail(participant *domain.Participant, event *domain.Event) string {
	eventDate := event.Date.Format("Monday, 02 January 2006 at 15:04")

	return `
Hi ` + participant.Name + `,

Thank you for registering for ` + event.Name + `!

Your unique QR code has been generated. Please check the HTML version of this email to view your QR code.

Event Details:
- Event Name: ` + event.Name + `
- Date & Time: ` + eventDate + `
- Venue: ` + event.Venue + `

How to Check-In:
1. Show your QR code at the registration desk
2. Your attendance will be recorded instantly

Important: This QR code is unique to you. Do not share it with others.

We look forward to seeing you!

Best regards,
EventCheck.in Team
`
}

// BuildTestEmail membuat test email
func BuildTestEmail(recipientName string) string {
	return `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Test Email</title>
</head>
<body style="font-family: Arial, sans-serif; padding: 20px;">
    <h2>🎉 Test Email from EventCheck.in</h2>
    <p>Hi ` + recipientName + `,</p>
    <p>This is a test email to verify your SMTP configuration is working correctly.</p>
    <p>If you received this email, your email service is configured properly!</p>
    <p>Best regards,<br>EventCheck.in Team</p>
</body>
</html>
`
}
