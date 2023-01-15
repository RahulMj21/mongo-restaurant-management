package controllers

import (
	"context"
	"time"

	"github.com/RahulMj21/mongo-restaurant-management/database"
	"github.com/RahulMj21/mongo-restaurant-management/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

var validate = validator.New()

func GetFoods(c *gin.Context) {

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
	return 0
}

func toFixed(num float64, precision int) float64 {
	return 1.0
}

func UpdateFood(c *gin.Context) {

}
