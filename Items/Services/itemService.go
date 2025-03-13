package Services

import (
	"TestApp/Items/Repositories"
	"TestApp/Models"
)

type ItemsService interface {
	CreateInvoice(invoice *Models.Invoice) error
	GetInvoicesByUserID(userID uint) ([]Models.Invoice, error)
}

type ItemService struct {
	repo *Repositories.ItemsRepository
}

func NewItemsService(repo Repositories.ItemsRepository) ItemService {
	return ItemService{repo: &repo}
}
