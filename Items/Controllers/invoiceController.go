// controllers/invoice.go
package Controllers

import (
	"TestApp/Items/Services"
	"TestApp/Models"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type InvoiceController struct {
	Service *Services.InvoiceService
}

func NewInvoiceController(service *Services.InvoiceService) *InvoiceController {
	return &InvoiceController{Service: service}
}

func (ctrl *InvoiceController) CreateInvoice(c echo.Context) error {
	var invoice Models.Invoice
	userID := c.Get("user_id").(uint64)
	if err := c.Bind(&invoice); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	invoice.UserID = uint(userID)
	if err := ctrl.Service.CreateInvoice(&invoice); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, invoice)
}

func (ctrl *InvoiceController) GetInvoicesByUserID(c echo.Context) error {
	userID := c.Get("user_id").(uint64)
	invoices, err := ctrl.Service.GetInvoicesByUserID(uint(userID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, invoices)
}

func (ctrl *InvoiceController) GetInvoiceByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid invoice ID"})
	}

	invoice, err := ctrl.Service.GetInvoiceByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Invoice not found"})
	}

	return c.JSON(http.StatusOK, invoice)
}

func (ctrl *InvoiceController) ExportPDF(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid invoice ID"})
	}
	err = ctrl.Service.ExportToPDF(uint(id), c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return nil
}
