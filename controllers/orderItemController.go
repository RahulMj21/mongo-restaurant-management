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

type OrderItemPack struct {
	TableId    *string
	OrderItems []models.OrderItem
}

var orderItemCollection = database.OpenCollection(database.Client, "order_item")

func GetOrderItems(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	cursor, err := orderItemCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(500, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	allOrderItems := []primitive.M{}
	if err := cursor.All(ctx, &allOrderItems); err != nil {
		c.JSON(500, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "success", "data": allOrderItems})
}

func GetOrderItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	orderItemId := c.Param("id")
	if orderItemId == "" {
		c.JSON(400, gin.H{"status": "fail", "message": "order_item id cannot be empty"})
		return
	}

	orderItem := models.OrderItem{}

	err := orderItemCollection.FindOne(ctx, bson.D{{Key: "order_item_id", Value: orderItemId}}).Decode(&orderItem)
	if err != nil {
		c.JSON(500, gin.H{"status": "fail", "message": "cannot find the order_item"})
		return
	}

	c.JSON(200, gin.H{"status": "success", "data": orderItem})
}

func GetOrderItemsByOrderId(c *gin.Context) {
	orderID := c.Param("order_id")
	if orderID == "" {
		c.JSON(400, gin.H{"status": "fail", "message": "order id cannot be empty"})
		return
	}
	allOrderItems, err := ItemsByOrderId(orderID)
	if err != nil {
		c.JSON(500, gin.H{"status": "fail", "message": "cannot get order items"})
		return
	}

	c.JSON(200, gin.H{"status": "success", "data": allOrderItems})
}

func ItemsByOrderId(id string) (OrderItems []primitive.M, err error) {
	return
}

func CreateOrderItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	orderItemPack := OrderItemPack{}
	order := models.Order{}
	if err := c.BindJSON(&orderItemPack); err != nil {
		c.JSON(400, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	order.OrderDate, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.TableId = orderItemPack.TableId
	order_id := OrderItemOrderCreator(order)

	orderItemsToBeInserted := []interface{}{}

	for _, orderItem := range orderItemPack.OrderItems {
		orderItem.OrderId = order_id
		validationErr := validate.Struct(orderItem)
		if validationErr != nil {
			c.JSON(400, gin.H{"status": "fail", "message": validationErr.Error()})
			return
		}
		orderItem.ID = primitive.NewObjectID()
		orderItem.OrderItemId = orderItem.ID.Hex()
		orderItem.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		orderItem.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		num := toFixed(*orderItem.UnitPrice, 2)
		orderItem.UnitPrice = &num

		orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)
	}
	insertedItems, err := orderItemCollection.InsertMany(ctx, orderItemsToBeInserted)
	if err != nil {
		c.JSON(500, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	c.JSON(201, gin.H{"status": "success", "data": insertedItems})
}

func UpdateOrderItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	orderItemId := c.Param("id")
	orderItem := models.OrderItem{}

	if orderItemId == "" {
		c.JSON(400, gin.H{"status": "fail", "message": "id cannot be empty"})
		return
	}
	if err := c.BindJSON(&orderItem); err != nil {
		c.JSON(400, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var orderItemObj primitive.D

	if orderItem.UnitPrice != nil {
		unitPrice := toFixed(*orderItem.UnitPrice, 2)
		orderItemObj = append(orderItemObj, bson.E{Key: "unit_price", Value: &unitPrice})
	}
	if orderItem.FoodId != nil {
		orderItemObj = append(orderItemObj, bson.E{Key: "food_id", Value: orderItem.FoodId})
	}
	if orderItem.Quantity != nil {
		orderItemObj = append(orderItemObj, bson.E{Key: "quantity", Value: orderItem.Quantity})
	}

	orderItem.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	orderItemObj = append(orderItemObj, bson.E{Key: "updated_at", Value: orderItem.UpdatedAt})

	filter := bson.D{{Key: "order_item_id", Value: orderItemId}}
	upsert := true
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	result, err := orderItemCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: orderItemObj}}, &opt)
	if err != nil {
		c.JSON(500, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "fail", "data": result})
}
