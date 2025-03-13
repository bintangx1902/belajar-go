package Services

import (
	"TestApp/Items/Repositories"
	"TestApp/Models"
	"TestApp/Utils"
	"fmt"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/orientation"
	"github.com/johnfercher/maroto/v2/pkg/consts/pagesize"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
)

type InvoiceService struct {
	Repo *Repositories.InvoiceRepository
}

func NewInvoiceService(repo *Repositories.InvoiceRepository) *InvoiceService {
	return &InvoiceService{Repo: repo}
}

func (ser *InvoiceService) CreateInvoice(invoice *Models.Invoice) error {
	for _, item := range invoice.Items {
		var existingItem Models.Items
		if err := ser.Repo.DB.First(&existingItem, item.ID).Error; err != nil {
			return fmt.Errorf("item with ID %d does not exist", item.ID)
		}
	}

	return ser.Repo.CreateInvoice(invoice)
}

func (ser *InvoiceService) GetInvoicesByUserID(userID uint) ([]Models.Invoice, error) {
	return ser.Repo.GetInvoicesByUserID(userID)
}

func (ser *InvoiceService) GetInvoiceByID(id uint) (Models.Invoice, error) {
	return ser.Repo.GetInvoiceByID(id)
}

func (ser *InvoiceService) ExportToPDF(id uint, ctx echo.Context) error {
	invoice, err := ser.Repo.GetInvoiceByID(id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": "invoice not found"})
	}

	cfg := config.NewBuilder().
		WithOrientation(orientation.Vertical).
		WithPageSize(pagesize.A4).
		WithLeftMargin(15).
		WithRightMargin(15).
		WithBottomMargin(15).
		WithTopMargin(15).
		Build()

	m := maroto.New(cfg)

	data := Utils.ConvertInvoiceItemsToContent(invoice.Items)
	Utils.ExportPDF(m,
		fmt.Sprintf("Invoice id %d", invoice.ID),
		fmt.Sprintf("This report generate at "),
		data,
	)

	document, err := m.Generate()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate PDF"})
	}

	filePath := fmt.Sprintf("pdf/invoice_%d.pdf", invoice.ID)
	err = document.Save(filePath)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"msg":    "Failed to save PDF",
			"error:": err.Error(),
		})
	}

	defer os.Remove(filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open PDF"})
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get file size"})
	}

	ctx.Response().Header().Set("Content-Type", "application/pdf")
	ctx.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=invoice_%d.pdf", invoice.ID))
	ctx.Response().Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))

	_, err = io.Copy(ctx.Response().Writer, file)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to stream PDF"})
	}
	return nil
}
