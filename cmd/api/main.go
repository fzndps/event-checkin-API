package main

import (
	"fmt"
	"log"

	"github.com/fzndps/eventcheck/config"
	"github.com/fzndps/eventcheck/internal/delivery/http"
	"github.com/fzndps/eventcheck/internal/delivery/http/middleware"
	"github.com/fzndps/eventcheck/internal/infrastructure/database"
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

	// initialize repo layer
	organizerRepo := mysql.NewOrganizerRepositoryImpl(db)
	eventRepo := mysql.NewEventRepository(db)
	participantRepo := mysql.NewParticipantRepository(db)

	// Initialize service/usecase layer
	authUsecase := usecase.NewAuthUsecase(organizerRepo, jwtManager, cfg)
	eventUsecase := usecase.NewEventUsecase(eventRepo, participantRepo)
	participantUsecase := usecase.NewParticipantUsecase(eventRepo, participantRepo)

	// initialize handler layer
	authHandler := http.NewAutHandler(authUsecase)
	eventHandler := http.NewEventHandler(*eventUsecase, participantUsecase)

	authMiddleware := middleware.NewAuthMiddleware(jwtManager)

	router := http.SetupRouter(&http.RouterConfig{
		AuthHandler:    authHandler,
		EventHandler:   eventHandler,
		AuthMiddleware: authMiddleware,
	})

	addr := fmt.Sprintf(":%s", cfg.App.Port)
	fmt.Printf("Server starting on http://localhost:%s\n", addr)

	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
