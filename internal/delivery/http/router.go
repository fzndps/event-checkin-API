package http

import (
	"net/http"

	"github.com/fzndps/eventcheck/internal/delivery/http/middleware"
	"github.com/gin-gonic/gin"
)

type RouterConfig struct {
	AuthHandler    *AuthHandler
	EventHandler   *EventHandler
	QREmailHandler *QREmailHandler
	AuthMiddleware *middleware.AuthMiddleware
}

func SetupRouter(cfg *RouterConfig) *gin.Engine {
	router := gin.Default()

	router.Use(corsMiddleware())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "EventCheck.in API is running",
		})
	})

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", cfg.AuthHandler.Register)
			auth.POST("/login", cfg.AuthHandler.Login)

			auth.GET("/profile", cfg.AuthMiddleware.AuthRequired(), cfg.AuthHandler.GetProfile)
		}

		events := v1.Group("/events")
		events.Use(cfg.AuthMiddleware.AuthRequired())
		{
			events.POST("", cfg.EventHandler.CreateEvent)
			events.GET("", cfg.EventHandler.ListEvents)
			events.GET("/:eventID", cfg.EventHandler.GetEventDetail)
			events.PUT("/:eventID", cfg.EventHandler.UpdateEvent)
			events.DELETE("/:eventID", cfg.EventHandler.DeleteEvent)
			events.POST("/:eventID/participants/upload", cfg.EventHandler.UploadParticipants)
			events.GET("/:eventID/participants", cfg.EventHandler.ListParticipant)

			events.POST("/:eventID/send-qr", cfg.QREmailHandler.SendQRCodes)
			events.POST("/:eventID/participants/:participantID/resend-qr", cfg.QREmailHandler.ResendQRCode)

		}

		email := v1.Group("/email")
		{
			email.POST("/test", cfg.QREmailHandler.SendTestEmail)
		}
	}

	return router
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
