package Repositories

import (
	"TestApp/Models"
	"gorm.io/gorm"
)

type ItemRepository struct {
	DB *gorm.DB
}

type ItemsRepository interface {
	CreateInvoice(invoice *Models.Invoice) error
	GetInvoicesByUserID(userID uint) ([]Models.Invoice, error)
}

func NewItemsRepository(db *gorm.DB) ItemRepository {
	return ItemRepository{DB: db}
}

func (repo *ItemRepository) CreateItem(invoice *Models.Invoice) error {
	return repo.DB.Create(invoice).Error
}

func (repo *ItemRepository) GetAllItems(userID uint) ([]Models.Items, error) {
	var items []Models.Items
	err := repo.DB.
		Model(&Models.Items{}).
		Find(&items).
		Error
	return items, err
}
