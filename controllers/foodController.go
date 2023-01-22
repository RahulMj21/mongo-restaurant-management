package controllers

import (
	"context"
	"math"
	"strconv"
	"time"

	"github.com/RahulMj21/mongo-restaurant-management/database"
	"github.com/RahulMj21/mongo-restaurant-management/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

var validate = validator.New()

func GetFoods(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	cancel()

	resultPerPage, err := strconv.Atoi(c.Query("resultPerPage"))
	if err != nil || resultPerPage < 1 {
		resultPerPage = 10
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	startIndex, err := strconv.Atoi(c.Query("startIndex"))
	if err != nil || startIndex < 0 {
		startIndex = (page - 1) * resultPerPage
	}

	matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
	groupStage := bson.D{{Key: "$group", Value: bson.D{
		{Key: "id", Value: bson.D{{Key: "_id", Value: "null"}}},
		{Key: "total_count", Value: bson.D{{Key: "$sum", Value: "1"}}},
		{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
	}}}
	projectState := bson.D{{Key: "$project", Value: bson.D{
		{Key: "_id", Value: 0},
		{Key: "total_count", Value: 1},
		{Key: "food_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, resultPerPage}}}},
	}}}

	result, err := foodCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectState})
	if err != nil {
		c.JSON(500, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	foods := []primitive.D{}

	if err = result.All(ctx, &foods); err != nil {
		c.JSON(500, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "success", "data": foods})
}

func GetFood(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	food_id := c.Param("id")
	food := models.Food{}
	err := foodCollection.FindOne(ctx, bson.M{"food_id": food_id}).Decode(&food)
	defer cancel()
	if err != nil {
		c.JSON(404, gin.H{"error": "food not found"})
	}
	c.JSON(200, food)
}

func CreateFood(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	menu := models.Menu{}
	food := models.Food{}

	if err := c.BindJSON(&food); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	validationErr := validate.Struct(food)
	if validationErr != nil {
		c.JSON(400, gin.H{"error": validationErr.Error()})
		return
	}

	err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.MenuId}).Decode(&menu)
	defer cancel()
	if err != nil {
		c.JSON(404, gin.H{"error": "menu not found"})
		return
	}

	food.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	food.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	food.ID = primitive.NewObjectID()
	food.FoodId = food.ID.Hex()
	num := toFixed(*food.Price, 2)
	food.Price = &num

	result, err := foodCollection.InsertOne(ctx, food)
	defer cancel()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, result)
}

func round(num float64) int {
	return int(math.Round(num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(output)) / output
}

func UpdateFood(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	menu := models.Menu{}
	food := models.Food{}

	foodId := c.Param("id")

	if err := c.BindJSON(&food); err != nil {
		c.JSON(400, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	var foodObj primitive.D

	if food.Name != nil {
		foodObj = append(foodObj, bson.E{Key: "name", Value: food.Name})
	}
	if food.Price != nil {
		foodObj = append(foodObj, bson.E{Key: "price", Value: food.Price})
	}
	if food.FoodImage != nil {
		foodObj = append(foodObj, bson.E{Key: "food_image", Value: food.FoodImage})
	}
	if food.MenuId != nil {
		err := menuCollection.FindOne(ctx, bson.D{{Key: "menu_id", Value: food.MenuId}}).Decode(&menu)
		if err != nil {
			c.JSON(400, gin.H{
				"status":  "fail",
				"message": err.Error(),
			})
			return
		}
		foodObj = append(foodObj, bson.E{Key: "menu_id", Value: menu.MenuId})
	}

	food.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	foodObj = append(foodObj, bson.E{Key: "updated_at", Value: food.UpdatedAt})

	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	filter := bson.D{{Key: "food_id", Value: foodId}}

	result, err := foodCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: foodObj}}, &opt)
	if err != nil {
		c.JSON(500, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   result,
	})

}
