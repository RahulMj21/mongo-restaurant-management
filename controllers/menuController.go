package controllers

import (
	"context"
	"log"
	"time"

	"github.com/RahulMj21/mongo-restaurant-management/database"
	"github.com/RahulMj21/mongo-restaurant-management/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	menu := models.Menu{}

	err := c.BindJSON(&menu)
	if err != nil {
		c.JSON(400, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	validationErr := validate.Struct(menu)
	if validationErr != nil {
		c.JSON(400, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	menu.ID = primitive.NewObjectID()
	menu.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	menu.MenuId = menu.ID.Hex()

	insertedItem, err := menuCollection.InsertOne(ctx, menu)
	if err != nil {
		c.JSON(500, gin.H{
			"status": "fail",
			"data":   err.Error(),
		})
		return
	}

	var newItem primitive.D
	err = menuCollection.FindOne(ctx, bson.D{{Key: "_id", Value: insertedItem.InsertedID}}).Decode(&newItem)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	c.JSON(500, gin.H{
		"status": "success",
		"data":   newItem,
	})
}

func UpdateMenu(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	menu := models.Menu{}

	if err := c.BindJSON(&menu); err != nil {
		c.JSON(400, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	menuId := c.Param("id")
	filter := bson.D{{Key: "menu_id", Value: menuId}}

	var menuObj primitive.D

	if menu.StartDate != nil && menu.EndDate != nil {
		if !inTimeSpan(*menu.StartDate, *menu.EndDate, time.Now()) {
			c.JSON(400, gin.H{"status": "fail", "message": "please retype the time"})
			return
		}

		menuObj = append(menuObj, bson.E{Key: "start_date", Value: menu.StartDate})
		menuObj = append(menuObj, bson.E{Key: "end_date", Value: menu.EndDate})

		if menu.Name != "" {
			menuObj = append(menuObj, bson.E{Key: "name", Value: menu.Name})
		}
		if menu.Category != "" {
			menuObj = append(menuObj, bson.E{Key: "category", Value: menu.Category})
		}
		menu.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menuObj = append(menuObj, bson.E{Key: "updated_at", Value: menu.UpdatedAt})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := menuCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: menuObj}}, &opt)
		if err != nil {
			c.JSON(500, gin.H{
				"status":  "fail",
				"message": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{"status": "success", "data": result})
		return
	}
	c.JSON(400, gin.H{"status": "fail", "message": "cannot update the menu"})
}

func inTimeSpan(start, end, check time.Time) bool {
	return start.After(time.Now()) && end.After(start)
}
