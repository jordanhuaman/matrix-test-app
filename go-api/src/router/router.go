package router

import (
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/jordanhuaman/go-api/src/clients"
	"github.com/jordanhuaman/go-api/src/database"
	handlers "github.com/jordanhuaman/go-api/src/handler"
	"github.com/jordanhuaman/go-api/src/middleware"
	"github.com/jordanhuaman/go-api/src/models"
	"github.com/jordanhuaman/go-api/src/services"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {
	// Middleware
	api := app.Group("/api", logger.New())
	api.Get("/", handlers.Hello)

	// Auth
	userRepo := models.NewUserRepository(database.Database.Db)
	refreshTokenRepo := models.NewRefreshTokenRepository(database.Database.Db)
	jwtSecret := os.Getenv("SECRET")
	if jwtSecret == "" {
		panic("SECRET environment variable is required")
	}
	accessTTL := 15 * time.Minute
	if ttlEnv := os.Getenv("ACCESS_TOKEN_TTL_MINUTES"); ttlEnv != "" {
		if ttlMinutes, err := strconv.Atoi(ttlEnv); err == nil && ttlMinutes > 0 {
			accessTTL = time.Duration(ttlMinutes) * time.Minute
		}
	}
	authService := services.NewAuthService(userRepo, refreshTokenRepo, jwtSecret, accessTTL)
	authHandler := handlers.NewAuthHandler(authService)

	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/register", authHandler.Register)
	auth.Post("/logout", middleware.Protected(), authHandler.Logout)
	auth.Post("/refresh-token", authHandler.RefreshToken)

	// User
	userHandler := handlers.NewUserHandler(userRepo, refreshTokenRepo)
	user := api.Group("/users", middleware.Protected())
	user.Get("/:id", userHandler.GetUser)
	user.Patch("/:id", userHandler.UpdateUser)
	user.Delete("/:id", userHandler.DeleteUser)

	// Matrix
	matrixInputRepo := models.NewMatrixInputRepository(database.Database.Db)
	matrixResultRepo := models.NewMatrixResultRepository(database.Database.Db)
	nodeClient := clients.NewNodeClient()
	matrixService := services.NewMatrixService(userRepo, matrixInputRepo, matrixResultRepo, nodeClient)
	matrixHandler := handlers.NewMatrixHandler(matrixService)

	matrix := api.Group("/matrix", middleware.Protected())
	matrix.Get("/qr", matrixHandler.ListResults)
	matrix.Get("/qr/:id", matrixHandler.GetResult)
	matrix.Post("/qr", matrixHandler.CreateQR)

}
