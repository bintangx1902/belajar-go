package Repositories

import (
	"TestApp/Models"
	"gorm.io/gorm"
)

type InvoicesRepository interface {
	CreateInvoice(invoice *Models.Invoice) error
	GetInvoicesByUserID(userID uint) ([]Models.Invoice, error)
	GetInvoiceByID(id uint) (Models.Invoice, error)
}

type InvoiceRepository struct {
	DB *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) *InvoiceRepository {
	return &InvoiceRepository{DB: db}
}

func (repo *InvoiceRepository) GetInvoicesByUserID(userID uint) ([]Models.Invoice, error) {
	var invoices []Models.Invoice
	err := repo.DB.Model(&Models.Invoice{}).
		Preload("Items").
		Preload("Users").
		Where("user_id = ?", userID).
		Find(&invoices).
		Error
	return invoices, err
}

func (repo *InvoiceRepository) GetInvoiceByID(id uint) (Models.Invoice, error) {
	var invoice Models.Invoice
	err := repo.DB.Model(&Models.Invoice{}).
		Preload("Items").
		Preload("Users").
		Where("id = ?", id).
		First(&invoice).Error
	return invoice, err
}

func (repo *InvoiceRepository) CreateInvoice(invoice *Models.Invoice) error {
	if err := repo.DB.Create(invoice).Error; err != nil {
		return err
	}

	for _, item := range invoice.Items {
		var existingInvoiceItem Models.InvoiceItems
		err := repo.DB.Where("invoice_id = ? AND items_id = ?", invoice.ID, item.ID).First(&existingInvoiceItem).Error

		if err == nil {
			continue
		} else if err != gorm.ErrRecordNotFound {
			return err
		}

		if err := repo.DB.Create(&Models.InvoiceItems{
			InvoiceID: invoice.ID,
			ItemsID:   item.ID,
		}).Error; err != nil {
			return err
		}
	}

	return nil
}
