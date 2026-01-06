package main

import (
	"fmt"
	"log"

	"github.com/fzndps/eventcheck/config"
	"github.com/fzndps/eventcheck/internal/delivery/http"
	"github.com/fzndps/eventcheck/internal/delivery/http/middleware"
	"github.com/fzndps/eventcheck/internal/infrastructure/database"
	"github.com/fzndps/eventcheck/internal/infrastructure/email"
	"github.com/fzndps/eventcheck/internal/infrastructure/qrcode"
	"github.com/fzndps/eventcheck/internal/repository/mysql"
	"github.com/fzndps/eventcheck/internal/usecase"
	"github.com/fzndps/eventcheck/pkg/jwt"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	defer db.Close()

	jwtManager := jwt.NewJWTManager(cfg.JWT.Secret)

	qrGenerator := qrcode.NewGenerator()

	emailService := email.NewEmailService(&cfg.SMTP)

	// initialize repo layer
	organizerRepo := mysql.NewOrganizerRepositoryImpl(db)
	eventRepo := mysql.NewEventRepository(db)
	participantRepo := mysql.NewParticipantRepository(db)

	// Initialize service/usecase layer
	authUsecase := usecase.NewAuthUsecase(organizerRepo, jwtManager, cfg)
	eventUsecase := usecase.NewEventUsecase(eventRepo, participantRepo)
	participantUsecase := usecase.NewParticipantUsecase(eventRepo, participantRepo)
	qrEmailUsecase := usecase.NewQREmailUsecase(eventRepo, participantRepo, qrGenerator, emailService)

	// initialize handler layer
	authHandler := http.NewAutHandler(authUsecase)
	eventHandler := http.NewEventHandler(*eventUsecase, participantUsecase)
	qrEmailHandler := http.NewQREmailHandler(qrEmailUsecase)

	authMiddleware := middleware.NewAuthMiddleware(jwtManager)

	router := http.SetupRouter(&http.RouterConfig{
		AuthHandler:    authHandler,
		EventHandler:   eventHandler,
		QREmailHandler: qrEmailHandler,
		AuthMiddleware: authMiddleware,
	})

	addr := fmt.Sprintf(":%s", cfg.App.Port)
	fmt.Printf("Server starting on http://localhost:%s\n", addr)

	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
