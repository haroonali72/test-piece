package db

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/test-piece/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DbClient     *mongo.Client
	DbCollection *mongo.Collection
)

func SetupDB() {
	clientOptions := options.Client().ApplyURI("mongodb://mongo:27017")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}

	DbClient = client
	DbCollection = DbClient.Database("test").Collection("users")
}

func GetUserByUsername(username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"username": username}
	var user models.User

	err := DbCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // User not found
		}
		return nil, err // Other error occurred
	}

	return &user, nil
}

func UpdateUser(username string, updateFields bson.M) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"username": username}
	update := bson.M{"$set": updateFields}

	result, err := DbCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetUsers(queryParams map[string][]string) ([]models.User, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}
	limit := int64(50)
	page := int64(1)

	// Parse query parameters
	if params, ok := queryParams["limit"]; ok {
		limit, _ = strconv.ParseInt(params[0], 10, 64)
		delete(queryParams, "limit")
	}
	if params, ok := queryParams["page"]; ok {
		page, _ = strconv.ParseInt(params[0], 10, 64)
		delete(queryParams, "page")
	}

	// Build filter based on query parameters
	for key, values := range queryParams {
		if len(values) == 1 {
			operator := ""
			value := values[0]
			if strings.HasPrefix(value, "eq") {
				operator = "$eq"
				value = strings.TrimPrefix(value, "eq")
			} else if strings.HasPrefix(value, "neq") {
				operator = "$ne"
				value = strings.TrimPrefix(value, "neq")
			} else if strings.HasPrefix(value, "gt") {
				operator = "$gt"
				value = strings.TrimPrefix(value, "gt")
			} else if strings.HasPrefix(value, "gte") {
				operator = "$gte"
				value = strings.TrimPrefix(value, "gte")
			} else if strings.HasPrefix(value, "lt") {
				operator = "$lt"
				value = strings.TrimPrefix(value, "lt")
			} else if strings.HasPrefix(value, "lte") {
				operator = "$lte"
				value = strings.TrimPrefix(value, "lte")
			}
			if operator != "" {
				filter[key] = bson.M{operator: value}
			}
		}
	}

	// Pagination options
	skip := (page - 1) * limit
	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSkip(skip)

	// Sorting options
	sortBy := "expiry_date" // Default sort field
	if params, ok := queryParams["sortBy"]; ok {
		sortBy = params[0]
		delete(queryParams, "sortBy")
	}
	order := 1 // Ascending order by default
	if params, ok := queryParams["order"]; ok {
		if params[0] == "desc" {
			order = -1
		}
		delete(queryParams, "order")
	}
	findOptions.SetSort(bson.D{{Key: sortBy, Value: order}})

	// Perform the query
	cursor, err := DbCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Count total records
	totalRecords, err := DbCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Extract the users from the cursor
	var users []models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, totalRecords, nil
}
