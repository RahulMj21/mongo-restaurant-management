package controllers

import (
	"strconv"
	"time"

	"github.com/RahulMj21/mongo-restaurant-management/database"
	"github.com/RahulMj21/mongo-restaurant-management/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

var userCollection = database.OpenCollection(database.Client, "user")

func GetUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	recordsPerPage, recordsErr := strconv.Atoi(c.Query("recordsPerPage"))
	if recordsErr != nil || recordsPerPage <= 0 {
		recordsPerPage = 10
	}

	page, pageErr := strconv.Atoi(c.Query("page"))
	if pageErr != nil || page <= 0 {
		page = 1
	}

	startIndex := (page - 1) * recordsPerPage

	matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "_id", Value: 1},
		{Key: "first_name", Value: 1},
		{Key: "last_name", Value: 1},
		{Key: "password", Value: 0},
		{Key: "email", Value: 1},
		{Key: "avatar", Value: 1},
		{Key: "phone", Value: 1},
		{Key: "access_token", Value: 0},
		{Key: "refresh_token", Value: 0},
		{Key: "created_at", Value: 1},
		{Key: "updated_at", Value: 1},
		{Key: "user_id", Value: 1},
	}}}
	groupStage := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}},
		{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
		{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
	}}}
	projectStage2 := bson.D{{Key: "$project", Value: bson.D{
		{Key: "_id", Value: 0},
		{Key: "total_count", Value: 1},
		{Key: "users", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordsPerPage}}}},
	}}}

	cursor, err := userCollection.Aggregate(ctx, mongo.Pipeline{matchStage, projectStage, groupStage, projectStage2})
	if err != nil {
		c.JSON(500, gin.H{"status": "fail", "message": "cannot get users"})
		return
	}

	users := []primitive.M{}
	if err := cursor.All(ctx, &users); err != nil {
		c.JSON(500, gin.H{"status": "fail", "message": "error while listing users"})
		return
	}

	c.JSON(200, gin.H{"status": "success", "data": users})
}

func GetUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	userId := c.Param("id")
	if userId == "" {
		c.JSON(400, gin.H{"status": "fail", "message": "user_id cannot be empty"})
		return
	}

	user := models.User{}
	opt := options.FindOne().SetProjection(bson.D{
		{Key: "password", Value: 0},
		{Key: "access_token", Value: 0},
		{Key: "refresh_token", Value: 0},
	})

	err := userCollection.FindOne(ctx, bson.D{{Key: "user_id", Value: userId}}, opt).Decode(&user)
	if err != nil {
		c.JSON(400, gin.H{"status": "fail", "message": "user not found"})
		return
	}

	c.JSON(200, gin.H{"status": "success"})
}

func SignUp(c *gin.Context) {
	c.JSON(200, "hello")
}

func Login(c *gin.Context) {
	c.JSON(200, "hello")
}

func HashPassword(password string) string {
	return password
}

func VerifyPassword(hashedPassword string, password string) (bool, string) {
	return password == hashedPassword, "hello"
}
