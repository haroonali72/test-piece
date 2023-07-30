package main

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/test-piece/db"
	"github.com/test-piece/handlers"
)

func main() {
	app := fiber.New()

	// Connect to MongoDB
	db.SetupDB()
	defer db.DbClient.Disconnect(context.Background())

	app.Post("/user", handlers.CreateUser)
	app.Get("/user/:username", handlers.GetUser)
	app.Patch("/user/:username", handlers.UpdateUser)
	app.Get("/users", handlers.GetUsers)

	app.Listen(":3000")
}
