package service

import (
	"hot-coffee/internal/dal"
	"hot-coffee/internal/errorHandle"
	"hot-coffee/models"
)

type MenuService interface {
	Create(item models.MenuItem) error
	GetAll() ([]models.MenuItem, error)
	GetItem(id string) (models.MenuItem, error)
	Update(item models.MenuItem, id string) error
	Delete(id string) error
}

type menuService struct {
	menuRepo dal.MenuRepository
}

func NewMenuService(menuRepo dal.MenuRepository) MenuService {
	return &menuService{
		menuRepo: menuRepo,
	}
}

func (m *menuService) Create(item models.MenuItem) error {
	if item.ID == "" || item.Name == "" || len(item.Ingredients) == 0 || item.Description == "" {
		return errorHandle.ErrorFormatJson
	}

	if item.Price <= 0 {
		return errorHandle.PriceLessZero
	}

	err := m.menuRepo.Create(item)
	return err
}

func (m *menuService) GetAll() ([]models.MenuItem, error) {
	result, err := m.menuRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return result, err
}

func (m *menuService) GetItem(id string) (models.MenuItem, error) {
	item, err := m.menuRepo.GetItem(id)
	if err != nil {
		return models.MenuItem{}, err
	}
	return item, err
}

func (m *menuService) Update(item models.MenuItem, id string) error {
	if item.ID == "" || item.Name == "" || len(item.Ingredients) == 0 || item.Description == "" {
		return errorHandle.ErrorFormatJson
	}
	if item.Price <= 0 {
		return errorHandle.PriceLessZero
	}
	if item.ID != id {
		return errorHandle.ChangeID
	}
	err := m.menuRepo.Update(item, id)
	return err
}

func (m *menuService) Delete(id string) error {
	err := m.menuRepo.Delete(id)
	return err
}
