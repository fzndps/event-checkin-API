package http

import (
	"strconv"

	"github.com/fzndps/eventcheck/internal/delivery/http/middleware"
	"github.com/fzndps/eventcheck/internal/domain"
	"github.com/fzndps/eventcheck/internal/usecase"
	"github.com/fzndps/eventcheck/pkg/validator"
	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	eventUsecase       usecase.EventUsecase
	participantUsecase usecase.ParticipantUsecase
}

func NewEventHandler(eventUsecase usecase.EventUsecase, participantUsecase *usecase.ParticipantUsecase) *EventHandler {
	return &EventHandler{
		eventUsecase:       eventUsecase,
		participantUsecase: *participantUsecase,
	}
}

func (h *EventHandler) CreateEvent(c *gin.Context) {
	// Dapatkan organizer id dari context
	organizerID, exists := middleware.GetOrganizerID(c)
	if !exists {
		validator.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	// Bind JSON request
	var req domain.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validator.BadRequestResponse(c, err.Error())
		return
	}

	event, err := h.eventUsecase.CreateEvent(c.Request.Context(), organizerID, &req)
	if err != nil {
		validator.InternalServerErrorResponse(c, "Failed to create event")
		return
	}

	validator.CreatedResponse(c, "Event created successfully", event)
}

func (h *EventHandler) ListEvents(c *gin.Context) {
	// Dapatkan organizer id dari context
	organizerID, exists := middleware.GetOrganizerID(c)
	if !exists {
		validator.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	// Dapatkan pagination parameter dari query string
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Panggil usecase
	response, err := h.eventUsecase.GetEventByOrganizer(c.Request.Context(), int64(organizerID), page, limit)
	if err != nil {
		validator.NotFoundResponse(c, err.Error())
		return
	}

	validator.SuccessResponse(c, "Event retrieved	sucessfully", response)
}

func (h *EventHandler) GetEventDetail(c *gin.Context) {
	// Dapatkan organizer id dari context
	organizerID, exists := middleware.GetOrganizerID(c)
	if !exists {
		validator.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	// Dapatkan event ID dari URL parameter
	eventID := c.Param("eventID")

	// Panggil usecase
	response, err := h.eventUsecase.GetEventDetail(c.Request.Context(), int64(organizerID), eventID)
	if err != nil {
		validator.NotFoundResponse(c, err.Error())
		return
	}

	validator.SuccessResponse(c, "Event retrieved	sucessfully", response)
}

func (h *EventHandler) UpdateEvent(c *gin.Context) {
	// Dapatkan organizer id dari context
	organizerID, exists := middleware.GetOrganizerID(c)
	if !exists {
		validator.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	eventID := c.Param("eventID")

	// bind JSON request
	var req domain.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validator.BadRequestResponse(c, err.Error())
		return
	}

	// panggil usecase
	event, err := h.eventUsecase.UpdateEvent(c.Request.Context(), int64(organizerID), eventID, &req)
	if err != nil {
		validator.InternalServerErrorResponse(c, err.Error())
		return
	}

	validator.SuccessResponse(c, "Event updated successfully", event)
}

func (h *EventHandler) DeleteEvent(c *gin.Context) {
	// Dapatkan organizer id dari context
	organizerID, exists := middleware.GetOrganizerID(c)
	if !exists {
		validator.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	// Dapatkan parameter event ID dari URL
	eventID := c.Param("eventID")

	// panggil usecase
	err := h.eventUsecase.DeleteEvent(c.Request.Context(), int64(organizerID), eventID)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "Event not found" {
			validator.NotFoundResponse(c, errMsg)
			return
		}
		validator.InternalServerErrorResponse(c, err.Error())
		return
	}

	validator.SuccessResponse(c, "Event deleted successfully", nil)
}

func (h *EventHandler) UploadParticipants(c *gin.Context) {
	// Dapatkan organizer id dari context
	organizerID, exists := middleware.GetOrganizerID(c)
	if !exists {
		validator.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	// Dapatkan parameter event ID dari URL
	eventID := c.Param("eventID")

	// dapatkan file dari multipart form
	file, err := c.FormFile("csv")
	if err != nil {
		validator.BadRequestResponse(c, "The CSV file must be uploaded with the key 'csv'")
		return
	}

	// Validasi file ekstensi
	if file.Header.Get("Content-Type") != "text/csv" && file.Header.Get("Content-Type") != "application/vnd.ms-excel" {
		validator.BadRequestResponse(c, "The file must be in CSV format")
		return
	}

	// Open file
	fileReader, err := file.Open()
	if err != nil {
		validator.InternalServerErrorResponse(c, err.Error())
		return
	}

	defer fileReader.Close()

	// panggil usecase
	response, err := h.participantUsecase.UploadParticipants(c.Request.Context(), int64(organizerID), eventID, fileReader)
	if err != nil {
		validator.InternalServerErrorResponse(c, err.Error())
		return
	}

	validator.SuccessResponse(c, "Upload participants successfully", response)
}

func (h *EventHandler) ListParticipant(c *gin.Context) {
	// Dapatkan organizer id dari context
	organizerID, exists := middleware.GetOrganizerID(c)
	if !exists {
		validator.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	// Dapatkan parameter event ID dari URL
	eventID := c.Param("eventID")

	// panggil usecase
	participants, err := h.participantUsecase.GetParticipantsByEvent(c.Request.Context(), int64(organizerID), eventID)
	if err != nil {
		validator.InternalServerErrorResponse(c, err.Error())
		return
	}

	validator.SuccessResponse(c, "Participants retrieve successfully", participants)
}
