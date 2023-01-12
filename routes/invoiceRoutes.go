package routes

import (
	"github.com/RahulMj21/mongo-restaurant-management/controllers"
	"github.com/gin-gonic/gin"
)

func InvoiceRoutes(app *gin.Engine) {
	app.GET("/invoices",controllers.GetInvoices)
	app.GET("/invoice/:id",controllers.GetInvoice)
	app.POST("/invoices",controllers.CreateInvoice)
	app.PATCH("/invoices/:id",controllers.UpdateInvoice)

}