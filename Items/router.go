package Items

import (
	"TestApp/Items/Controllers"
	"TestApp/Items/Repositories"
	"TestApp/Items/Services"
	"TestApp/Middlewares"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func ItemsRouter(route *echo.Group, db *gorm.DB) {
	repo := Repositories.NewInvoiceRepository(db)
	service := Services.NewInvoiceService(repo)
	controller := Controllers.NewInvoiceController(service)

	InvoiceGroup := route.Group("/invoices")
	// GET
	InvoiceGroup.GET("", controller.GetInvoicesByUserID, Middlewares.JWTMiddleware)
	InvoiceGroup.GET("/:id", controller.GetInvoiceByID, Middlewares.JWTMiddleware)
	InvoiceGroup.GET("/:id/download", controller.ExportPDF)

	// POST
	InvoiceGroup.POST("", controller.CreateInvoice, Middlewares.JWTMiddleware)
}
