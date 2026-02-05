package http

import (
	"errors"
	"net/http"

	"github.com/fzndps/eventcheck/internal/domain"
	"github.com/fzndps/eventcheck/internal/usecase"
	"github.com/fzndps/eventcheck/pkg/validator"
	"github.com/gin-gonic/gin"
)

type CheckInHandler struct {
	checkInUsecase *usecase.CheckInUsecase
}

func NewCheckInHandler(checkInUsecase *usecase.CheckInUsecase) *CheckInHandler {
	return &CheckInHandler{
		checkInUsecase: checkInUsecase,
	}
}

// menampilkan data
func (h *CheckInHandler) GetScanPage(c *gin.Context) {
	eventSlug := c.Param("event_slug")

	// get event data
	scanData, err := h.checkInUsecase.GetScanPageData(c.Request.Context(), eventSlug)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.HTML(http.StatusNotFound, "404.html", gin.H{
				"error": "Event not found",
			})
			return
		}
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{
			"error": "Internal server error",
		})
		return
	}

	c.HTML(http.StatusOK, "scan.html", gin.H{
		"EventSlug": scanData.EventSlug,
		"EventName": scanData.EventName,
		"EventDate": scanData.EventDate.Format("Monday, 02 January 2006"),
		"Venue":     scanData.Venue,
	})
}

func (h *CheckInHandler) VerifyPIN(c *gin.Context) {
	var req domain.VerifyPINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validator.BadRequestResponse(c, err.Error())
		return
	}

	eventSlug := c.Param("event_slug")

	// call usecase
	res, err := h.checkInUsecase.VerifyPIN(c.Request.Context(), eventSlug, req.ScannerPIN)
	if err != nil {
		validator.InternalServerErrorResponse(c, err.Error())
		return
	}

	if res.Valid {
		validator.SuccessResponse(c, "", res)
	} else {
		c.JSON(http.StatusUnauthorized, res)
	}
}

func (h *CheckInHandler) CheckedIn(c *gin.Context) {
	var req domain.CheckInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validator.BadRequestResponse(c, err.Error())
		return
	}

	// panggil usecase
	res, err := h.checkInUsecase.CheckIn(c.Request.Context(), req.QRToken)
	if err != nil {
		validator.InternalServerErrorResponse(c, err.Error())
		return
	}

	if res.Success {
		validator.SuccessResponse(c, "", res)
	} else {
		// mengembalikan response 200 dengan success=false
		validator.SuccessResponse(c, "", res)
	}
}

func (h *CheckInHandler) GetEventStats(c *gin.Context) {
	eventSlug := c.Param("event_slug")

	stats, err := h.checkInUsecase.GetEventStats(c.Request.Context(), eventSlug)
	if err != nil {
		statusCode, message := h.handleError(err)
		c.JSON(statusCode, gin.H{"error": message})
		return
	}

	validator.SuccessResponse(c, "", stats)
}

// handleError convert domain error ke HTTP status code
func (h *CheckInHandler) handleError(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrEventNotFound):
		return http.StatusNotFound, err.Error()

	case errors.Is(err, domain.ErrInvalidPIN):
		return http.StatusUnauthorized, err.Error()

	case errors.Is(err, domain.ErrQRTokenNotFound):
		return http.StatusNotFound, err.Error()

	default:
		return http.StatusInternalServerError, "Terjadi kesalahan server"
	}
}
