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

type InvoiceViewFormat struct {
	Invoice_id       string
	Payment_method   string
	Order_id         string
	Payment_status   *string
	Payment_due      interface{}
	Table_number     interface{}
	Payment_due_date time.Time
	Order_details    interface{}
}

var invoiceCollection = database.OpenCollection(database.Client, "invoice")

func GetInvoices(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	result, err := invoiceCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	var invoices bson.M

	if err := result.All(ctx, &invoices); err != nil {
		c.JSON(500, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   invoices,
	})
}

func GetInvoice(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	invoiceId := c.Param("id")

	invoice := models.Invoice{}
	err := invoiceCollection.FindOne(ctx, bson.D{{Key: "invoice_id", Value: invoiceId}}).Decode(&invoice)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	var invoiceView InvoiceViewFormat

	allOrderItems, err := ItemsByOrderId(invoiceId)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	invoiceView.Order_id = invoice.OrderId
	invoiceView.Payment_due_date = invoice.PaymentDueDate

	invoiceView.Payment_method = "null"
	if invoice.PaymentMethod != nil {
		invoiceView.Payment_method = *invoice.PaymentMethod
	}

	invoiceView.Invoice_id = invoice.InvoiceId
	invoiceView.Payment_status = invoice.PaymentStatus
	invoiceView.Payment_due = allOrderItems[0]["payment_due"]
	invoiceView.Table_number = allOrderItems[0]["table_number"]
	invoiceView.Order_details = allOrderItems[0]["order_items"]

	c.JSON(200, gin.H{
		"status": "success",
		"data":   invoiceView,
	})
}

func CreateInvoice(c *gin.Context) {
	_, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
}

func UpdateInvoice(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	invoice := models.Invoice{}
	invoiceId := c.Param("id")

	if err := c.BindJSON(&invoice); err != nil {
		c.JSON(400, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	filter := bson.D{{Key: "invoice_id", Value: invoiceId}}

	var invoiceObj primitive.D

	if invoice.PaymentMethod != nil {
		invoiceObj = append(invoiceObj, bson.E{Key: "payment_method", Value: invoice.PaymentMethod})
	}
	if invoice.PaymentStatus != nil {
		invoiceObj = append(invoiceObj, bson.E{Key: "payment_status", Value: invoice.PaymentStatus})
	}

	invoice.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	invoiceObj = append(invoiceObj, bson.E{Key: "updated_at", Value: invoice.UpdatedAt})

	upsert := true

	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	status := "PENDING"
	if invoice.PaymentStatus == nil {
		invoice.PaymentStatus = &status
	}

	result, err := invoiceCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: invoiceObj}}, &opt)
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
