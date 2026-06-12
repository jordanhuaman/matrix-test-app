package testutil

import (
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/jordanhuaman/go-api/src/clients"
	handlers "github.com/jordanhuaman/go-api/src/handler"
	"github.com/jordanhuaman/go-api/src/middleware"
	"github.com/jordanhuaman/go-api/src/models"
	"github.com/jordanhuaman/go-api/src/services"
	"gorm.io/gorm"
)

func NewTestApp(db *gorm.DB, jwtSecret string) *fiber.App {
	os.Setenv("SECRET", jwtSecret)
	app := fiber.New(fiber.Config{})

	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin, Content-Type, Accept, Authorization"},
		AllowMethods: []string{"GET, POST, PATCH, DELETE, OPTIONS"},
	}))

	api := app.Group("/api")
	api.Get("/", handlers.Hello)

	userRepo := models.NewUserRepository(db)
	refreshTokenRepo := models.NewRefreshTokenRepository(db)

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

	userHandler := handlers.NewUserHandler(userRepo, refreshTokenRepo)
	user := api.Group("/users", middleware.Protected())
	user.Get("/:id", userHandler.GetUser)
	user.Patch("/:id", userHandler.UpdateUser)
	user.Delete("/:id", userHandler.DeleteUser)

	matrixInputRepo := models.NewMatrixInputRepository(db)
	matrixResultRepo := models.NewMatrixResultRepository(db)
	nodeClient := clients.NewNodeClient()
	matrixService := services.NewMatrixService(userRepo, matrixInputRepo, matrixResultRepo, nodeClient)
	matrixHandler := handlers.NewMatrixHandler(matrixService)

	matrix := api.Group("/matrix", middleware.Protected())
	matrix.Get("/qr", matrixHandler.ListResults)
	matrix.Get("/qr/:id", matrixHandler.GetResult)
	matrix.Post("/qr", matrixHandler.CreateQR)

	return app
}
