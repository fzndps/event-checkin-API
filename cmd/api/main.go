package main

import (
	"fmt"
	"log"
	"time"

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

	db.Close()

	jwtManager := jwt.NewJWTManager(cfg.JWT.Secret, time.Duration(cfg.JWT.Expiry))

	organizerRepo := mysql.NewOrganizerRepositoryImpl(db)

	authUsecase := usecase.NewAuthUsecase(organizerRepo, jwtManager)

	authHandler := http.NewAutHandler(authUsecase)

	authMiddleware := middleware.NewAuthMiddleware(jwtManager)

	router := http.SetupRouter(&http.RouterConfig{
		AuthHandler:    authHandler,
		AuthMiddleware: authMiddleware,
	})

	addr := fmt.Sprintf(":%s", cfg.App.Port)
	fmt.Printf("Server starting on http://localhost:%s\n", addr)

	if err := router.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
