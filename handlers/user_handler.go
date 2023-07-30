package handlers

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/test-piece/db"
	"github.com/test-piece/models"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateUser(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return err
	}

	// Validate the user model
	if err := user.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Set the user's creation timestamp (optional)
	//userCreated := time.Now().Unix()

	// Add additional data if needed before inserting into the database
	//user.Created = userCreated

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := db.DbCollection.InsertOne(ctx, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Add the generated user ID to the response
	userID := result.InsertedID
	userIDMap := map[string]interface{}{"user_id": userID}

	return c.Status(fiber.StatusCreated).JSON(userIDMap)
}

func GetUser(c *fiber.Ctx) error {
	username := c.Params("username")
	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "username not found"})
	}

	user, err := db.GetUserByUsername(username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func UpdateUser(c *fiber.Ctx) error {
	username := c.Params("username")
	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "username not found"})
	}

	var updateFields bson.M
	if err := c.BodyParser(&updateFields); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	result, err := db.UpdateUser(username, updateFields)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	response := fiber.Map{
		"status":  true,
		"message": fmt.Sprintf("The user %s has been successfully updated!", username),
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func GetUsers(c *fiber.Ctx) error {
	limitStr := c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50 // Default limit value if not provided or invalid
	}

	queryParams := make(map[string][]string)
	err = c.QueryParser(queryParams)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": false, "message": "request payload is not valid"})
	}

	users, totalRecords, err := db.GetUsers(queryParams)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": false, "message": "Internal Server Error"})
	}

	totalPages := int64(math.Ceil(float64(totalRecords) / float64(limit)))

	response := fiber.Map{
		"total_records": totalRecords,
		"total_pages":   totalPages,
		"records":       users,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
