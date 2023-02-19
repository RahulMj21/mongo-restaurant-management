package controllers

import (
	"log"
	"strconv"
	"time"

	"github.com/RahulMj21/mongo-restaurant-management/database"
	"github.com/RahulMj21/mongo-restaurant-management/helpers"
	"github.com/RahulMj21/mongo-restaurant-management/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
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
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	user := models.User{}
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	validationErr := validate.Struct(user)
	if validationErr != nil {
		c.JSON(400, gin.H{"status": "fail", "message": validationErr.Error()})
		return
	}

	count, err := userCollection.CountDocuments(ctx, bson.D{{Key: "email", Value: user.Email}})
	if err != nil {
		c.JSON(500, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	if count > 0 {
		c.JSON(500, gin.H{"status": "fail", "message": "email already taken"})
		return
	}

	password := HashPassword(*user.Password)
	user.Password = &password
	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	user.UserId = user.ID.Hex()

	accessToken, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.FirstName, *user.LastName, user.UserId)

	user.AccessToken = &accessToken
	user.RefreshToken = &refreshToken

	insertedItem, err := userCollection.InsertOne(ctx, user)
	if err != nil || insertedItem.InsertedID == nil {
		c.JSON(500, gin.H{"status": "fail", "message": "user creation failed"})
		return
	}

	newUser := models.User{
		ID:           user.ID,
		UserId:       user.UserId,
		Email:        user.Email,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
		Avatar:       user.Avatar,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}

	c.JSON(201, gin.H{"status": "success", "data": newUser})
}

type LoginBody struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	body := LoginBody{}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	validationErr := validate.Struct(body)
	if validationErr != nil {
		c.JSON(400, gin.H{"status": "fail", "message": validationErr.Error()})
		return
	}

	user := models.User{}
	err := userCollection.FindOne(ctx, bson.D{{Key: "email", Value: &body.Email}}).Decode(&user)
	if err != nil {
		c.JSON(400, gin.H{"status": "fail", "message": "wrong email"})
		return
	}

	isVerified := VerifyPassword(*user.Password, body.Password)
	if isVerified == false {
	}
	c.JSON(400, gin.H{"status": "fail", "message": "wrong password"})
	return

	// need to add session model for storing refreshToken
	// we need to createTokens and store them and send to user

	newUser := models.User{
		ID:           user.ID,
		UserId:       user.UserId,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		Avatar:       user.Avatar,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		AccessToken:  user.AccessToken,
		RefreshToken: user.RefreshToken,
	}

	c.JSON(200, gin.H{"status": "success", "data": newUser})
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false
	}
	return true
}
