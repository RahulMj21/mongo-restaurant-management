package controllers

import (
	"context"
	"time"

	"github.com/RahulMj21/mongo-restaurant-management/database"
	"github.com/RahulMj21/mongo-restaurant-management/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var tableCollection = database.OpenCollection(database.Client, "table")

func GetTables(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	cursor, err := tableCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(500, gin.H{"status": "fail", "message": "cannot get tables"})
		return
	}
	tables := []bson.M{}

	if err := cursor.All(ctx, &tables); err != nil {
		c.JSON(500, gin.H{"status": "fail", "message": "unable to listing tables"})
		return
	}

	c.JSON(200, gin.H{"status": "success", "data": tables})
}

func GetTable(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	tableId := c.Param("id")
	if tableId == "" {
		c.JSON(400, gin.H{"status": "fail", "message": "id cannot be empty"})
		return
	}
	table := models.OrderItem{}
	err := tableCollection.FindOne(ctx, bson.D{{Key: "table_id", Value: tableId}}).Decode(&table)
	if err != nil {
		c.JSON(500, gin.H{"status": "fail", "message": "cannot get the table"})
		return
	}

	c.JSON(200, gin.H{"status": "success", "data": table})
}

func CreateTable(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	table := models.Table{}
	if err := c.BindJSON(&table); err != nil {
		c.JSON(400, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	validationErr := validate.Struct(table)
	if validationErr != nil {
		c.JSON(400, gin.H{"status": "fail", "message": validationErr.Error()})
		return
	}

	table.ID = primitive.NewObjectID()
	table.TableId = table.ID.Hex()
	table.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	table.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	insertedItem, err := tableCollection.InsertOne(ctx, table)
	if err != nil {
		c.JSON(500, gin.H{"status": "fail", "message": "failed to create table"})
		return
	}

	newTable := models.Table{}
	if err := tableCollection.FindOne(ctx, bson.D{{Key: "_id", Value: insertedItem.InsertedID}}).Decode(&newTable); err != nil {
		c.JSON(500, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	c.JSON(201, gin.H{"status": "success", "data": newTable})
}

func UpdateTable(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	tableId := c.Param("id")
	if tableId == "" {
		c.JSON(400, gin.H{"status": "fail", "message": "id cannot be empty"})
		return
	}
	table := models.Table{}
	if err := c.BindJSON(&table); err != nil {
		c.JSON(400, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	tableObj := primitive.D{}

	if table.NumberOfGuests != nil {
		tableObj = append(tableObj, bson.E{Key: "number_of_guests", Value: table.NumberOfGuests})
	}
	if table.TableNumber != nil {
		tableObj = append(tableObj, bson.E{Key: "table_number", Value: table.TableNumber})
	}
	table.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	tableObj = append(tableObj, bson.E{Key: "updated_at", Value: table.UpdatedAt})

	filter := bson.D{{Key: "table_id", Value: tableId}}
	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	result, err := tableCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: tableObj}}, &opt)
	if err != nil {
		c.JSON(500, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "success", "data": result})
}
