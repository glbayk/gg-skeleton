package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/glbayk/gg-skeleton/controllers"
	"github.com/glbayk/gg-skeleton/globals"
	"github.com/glbayk/gg-skeleton/middleware"
	"github.com/glbayk/gg-skeleton/models"
	"github.com/glbayk/gg-skeleton/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/joho/godotenv"
)

var WaitGroup sync.WaitGroup

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		os.Exit(1)
	}

	models.Connect()
	services.MongoClient = services.MongoClientConnect()
}

func main() {
	app := fiber.New()
	fiberUses(app)
	api := app.Group("/api", logger.New())
	v1 := api.Group("/v1")

	app.Get("/metrics", monitor.New(monitor.Config{Title: "MyService Metrics Page"}))

	HealthCheckRouter(v1, "ping")
	UserRouter(v1, "user")
	AuthRouter(v1, "auth")

	gracefulShutdown(app)
}

func fiberUses(app *fiber.App) {
	app.Use(logger.New())
	app.Use(cors.New())
}

func gracefulShutdown(app *fiber.App) {
	go func() {
		if err := app.Listen(":" + os.Getenv("PORT")); err != nil {
			log.Println("Error starting server", "error", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	fmt.Println("Gracefully shutting down...")
	_ = app.Shutdown()
	fmt.Println("Running cleanup tasks...")
	globals.GetWaitGroup().Wait()
	models.Cleanup()
	fmt.Println("Server shutdown")
}

func HealthCheckRouter(parent fiber.Router, param string) {
	parent.Get(param, func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "pong",
		})
	})
}

func AuthRouter(parent fiber.Router, param string) {
	aC := controllers.AuthController{}

	auth := parent.Group(param)

	auth.Post("/register", aC.Register)
	auth.Post("/login", aC.Login)
	auth.Get("/activate/:token", aC.Activate)
	auth.Post("/forgot-password", aC.ForgotPassword)
}

func UserRouter(parent fiber.Router, param string) {
	uC := controllers.UserController{}

	user := parent.Group(param, middleware.Authenticated)

	user.Get("/me", uC.Me)
	user.Get("/refresh-token", uC.RefreshToken)
	user.Post("/change-password", uC.ChangePassword)
}
