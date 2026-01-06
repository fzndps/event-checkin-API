package http

import (
	"github.com/fzndps/eventcheck/internal/delivery/http/middleware"
	"github.com/fzndps/eventcheck/internal/domain"
	"github.com/fzndps/eventcheck/internal/usecase"
	"github.com/fzndps/eventcheck/pkg/validator"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase *usecase.AuthUsecase
}

func NewAutHandler(authUsecase *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		validator.BadRequestResponse(c, err.Error())
		return
	}

	organizer, err := h.authUsecase.Register(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "email already registered" {
			validator.BadRequestResponse(c, err.Error())
			return
		}

		validator.InternalServerErrorResponse(c, err.Error())
		return
	}

	responseData := map[string]any{
		"organizer": organizer.ToResponse(),
	}

	validator.CreatedResponse(c, "User registered successfully", responseData)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validator.BadRequestResponse(c, err.Error())
		return
	}

	organizer, err := h.authUsecase.Login(c.Request.Context(), &req)
	if err != nil {
		validator.UnauthorizedResponse(c, err.Error())
		return
	}

	responseData := map[string]any{
		"organizer": organizer.Organizer.ToResponse(),
		"token":     organizer.Token,
	}

	validator.SuccessResponse(c, "Login successfully", responseData)
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	organizerID, ok := middleware.GetOrganizerID(c)
	if !ok {
		validator.UnauthorizedResponse(c, "Organizer not authenticated")
		return
	}

	organizer, err := h.authUsecase.GetProfileByID(c.Request.Context(), organizerID)
	if err != nil {
		validator.NotFoundResponse(c, "Organizer not found")
		return
	}

	validator.SuccessResponse(c, "Profile retrieved successfully", organizer.ToResponse())
}
