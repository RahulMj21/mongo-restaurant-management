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

var ordersCollection = database.OpenCollection(database.Client, "order")

func GetOrders(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	cursor, err := ordersCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	var orders []primitive.M

	if err := cursor.All(ctx, &orders); err != nil {
		c.JSON(500, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   orders,
	})
}

func GetOrder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	orderId := c.Param("id")
	order := models.Order{}

	filter := bson.D{{Key: "order_id", Value: orderId}}
	if err := ordersCollection.FindOne(ctx, filter).Decode(&order); err != nil {
		c.JSON(500, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   order,
	})
}

func CreateOrder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	order := models.Order{}
	table := models.Table{}

	if err := c.BindJSON(&order); err != nil {
		c.JSON(400, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	err := validate.Struct(order)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	if order.TableId != nil {
		if err := tableCollection.FindOne(ctx, bson.D{{Key: "table_id", Value: order.TableId}}).Decode(&table); err != nil {
			c.JSON(400, gin.H{
				"status":  "fail",
				"message": err.Error(),
			})
			return
		}
	}

	order.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.OrderId = order.ID.Hex()

	insertedItem, err := ordersCollection.InsertOne(ctx, order)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	newItem := models.Order{}

	err = ordersCollection.FindOne(ctx, bson.D{{Key: "_id", Value: insertedItem.InsertedID}}).Decode(&newItem)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
	}
}

func UpdateOrder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	table := models.Table{}
	order := models.Order{}

	orderId := c.Param("id")

	if err := c.BindJSON(&order); err != nil {
		c.JSON(400, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	var orderObj primitive.D

	if order.TableId != nil {
		err := tableCollection.FindOne(ctx, bson.D{{Key: "table_id", Value: order.TableId}}).Decode(&table)
		if err != nil {
			c.JSON(400, gin.H{
				"status":  "fail",
				"message": err.Error(),
			})
			return
		}
		orderObj = append(orderObj, bson.E{Key: "table_id", Value: table.TableId})
	}

	order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	orderObj = append(orderObj, bson.E{Key: "updated_at", Value: order.UpdatedAt})

	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	filter := bson.D{{Key: "order_id", Value: orderId}}

	result, err := ordersCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: orderObj}}, &opt)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   result,
	})
}

func OrderItemOrderCreator(order models.Order) string {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	order.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID()
	order.OrderId = order.ID.Hex()

	ordersCollection.InsertOne(ctx, order)

	return order.OrderId
}
