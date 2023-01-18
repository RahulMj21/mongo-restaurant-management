package controllers

import (
	"context"
	"log"
	"time"

	"github.com/RahulMj21/mongo-restaurant-management/database"
	"github.com/RahulMj21/mongo-restaurant-management/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

func GetMenus(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	cursor, err := menuCollection.Find(ctx, bson.M{})
	defer cancel()
	if err != nil {
		c.JSON(500, gin.H{"error": "error while fetching menus"})
	}

	var allMenus []bson.M
	if err = cursor.All(ctx, &allMenus); err != nil {
		log.Fatal(err)
	}
	c.JSON(200, allMenus)
}

func GetMenu(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	menu_id := c.Param("id")
	menu := models.Menu{}

	err := menuCollection.FindOne(ctx, bson.D{{Key: "menu_id", Value: menu_id}}).Decode(&menu)
	defer cancel()
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
	}

	c.JSON(200, gin.H{"status": "success", "data": menu})
}

func CreateMenu(c *gin.Context) {

}

func UpdateMenu(c *gin.Context) {

}
