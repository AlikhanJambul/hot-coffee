package service

import (
	"hot-coffee/internal/dal"
	"hot-coffee/internal/errorHandle"
	"hot-coffee/models"
)

type InventoryService interface {
	Create(item models.InventoryItem) error
	GetAll() ([]models.InventoryItem, error)
	GetItem(id string) (models.InventoryItem, error)
	Update(item models.InventoryItem, id string) error
	Delete(id string) error
}

type inventoryService struct {
	inventoryRepo dal.InventoryRepository
}

func NewInventoryService(inventoryRepo dal.InventoryRepository) InventoryService {
	return &inventoryService{
		inventoryRepo: inventoryRepo,
	}
}

func (i *inventoryService) Create(item models.InventoryItem) error {
	if item.IngredientID == "" || item.Name == "" || item.Unit == "" {
		return errorHandle.ErrorFormatJson
	}
	if item.Quantity <= 0 {
		return errorHandle.QuantityLessZero
	}

	err := i.inventoryRepo.Create(item)
	return err
}

func (i *inventoryService) GetAll() ([]models.InventoryItem, error) {
	items, err := i.inventoryRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (i *inventoryService) GetItem(id string) (models.InventoryItem, error) {
	item, err := i.inventoryRepo.GetItem(id)
	if err != nil {
		return item, err
	}
	return item, nil
}

func (i *inventoryService) Update(item models.InventoryItem, id string) error {
	if id != item.IngredientID {
		return errorHandle.ChangeID
	}

	if item.Quantity <= 0 {
		return errorHandle.QuantityLessZero
	}

	err := i.inventoryRepo.Update(item, id)
	return err
}

func (i *inventoryService) Delete(id string) error {
	err := i.inventoryRepo.Delete(id)
	return err
}
