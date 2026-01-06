package http

import (
	"log"
	"net/http"

	"github.com/fzndps/eventcheck/internal/usecase"
	"github.com/fzndps/eventcheck/pkg/validator"
	"github.com/gin-gonic/gin"
)

type QREmailHandler struct {
	qrEmailUsecase *usecase.QREmailUsecae
}

func NewQREmailHandler(qrEmailUsecase *usecase.QREmailUsecae) *QREmailHandler {
	return &QREmailHandler{
		qrEmailUsecase: qrEmailUsecase,
	}
}

func (h *QREmailHandler) SendQRCodes(c *gin.Context) {
	// Get organizerID dari context
	organizerID, exists := c.Get("organizer_id")
	if !exists {
		validator.UnauthorizedResponse(c, "Unauthorized")
		return
	}

	var organizerID64 int64

	// FIX untuk Error 1 & Robustness:
	// Gunakan Type Switch untuk menangani berbagai kemungkinan tipe angka dari Context
	switch v := organizerID.(type) {
	case int:
		organizerID64 = int64(v)
	case int64:
		organizerID64 = v
	case float64: // Jaga-jaga jika parser JWT mengembalikan float64
		organizerID64 = int64(v)
	default:
		// FIX untuk Error 2: Tambahkan return agar kode berhenti di sini jika error
		validator.BadRequestResponse(c, "organizerID is not a valid number")
		return
	}

	// Get eventID dari parameter
	eventID := c.Param("eventID")

	// Call usecase
	response, err := h.qrEmailUsecase.SendQRCodes(
		c.Request.Context(),
		organizerID64,
		eventID,
	)
	if err != nil {
		log.Print("error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	statusCode := http.StatusOK
	message := "The QR code has been successfully sent."

	if response.EmailsFailed > 0 {
		statusCode = http.StatusMultiStatus
		message = "QR codes were sent with several failures"
	}

	if response.EmailsSent == 0 && response.EmailsFailed > 0 {
		statusCode = http.StatusInternalServerError
		message = "Failed to send all QR codes"
	}

	if response.TotalParticipants == 0 {
		message = "All QR codes have already been sent"
	}

	c.JSON(statusCode, gin.H{
		"message": message,
		"data":    response,
	})
}

// sendTesEmail testing email pada endpoint POST /api/email/test
func (h *QREmailHandler) SendTestEmail(c *gin.Context) {
	// reqeust body
	var req struct {
		Email string `json:"email" binding:"required,email"`
		Name  string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		validator.BadRequestResponse(c, err.Error())
		return
	}

	// send test email
	err := h.qrEmailUsecase.SendTestEmail(c.Request.Context(), req.Email, req.Name)
	if err != nil {
		validator.InternalServerErrorResponse(c, "Failed to send test email: "+err.Error())
		return
	}

	validator.SuccessResponse(c, "Test email sent successfully"+req.Email, nil)

}

// func (h *QREmailHandler) handleError(err error) (int, string) {
// 	switch {
// 	case errors.Is(err, domain.ErrEventNotFound):
// 		return http.StatusNotFound, err.Error()
// 	case errors.Is(err, domain.ErrUnauthorizedAccess):
// 		return http.StatusForbidden, err.Error()

// 	default:
// 		return http.StatusInternalServerError, "Internal server error"
// 	}

// }
