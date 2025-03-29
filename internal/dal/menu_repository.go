package dal

import (
	"encoding/json"
	"hot-coffee/internal/errorHandle"
	"hot-coffee/models"
	"io/ioutil"
)

type MenuRepository interface {
	Create(item models.MenuItem) error
	GetAll() ([]models.MenuItem, error)
	GetItem(id string) (models.MenuItem, error)
	Update(item models.MenuItem, id string) error
	Delete(id string) error
	ExistsByID(id string) bool
	MenuCalcuation(id string, quantity float64) error
	MenuConsumptionOfIngredients(id string, quantity float64, plus bool) error
	SumOfOrder(id string) (float64, error)
}

type MenuRepo struct {
	path           string
	menuMap        []models.MenuItem
	inventoryRepos InventoryRepository
}

func NewMenuRepository(path string, inventoryRepos InventoryRepository) MenuRepository {
	existingData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}

	var menu []models.MenuItem

	if len(existingData) > 0 {
		err = json.Unmarshal(existingData, &menu)
		if err != nil {
			return nil
		}
	}

	return &MenuRepo{path: path, menuMap: menu, inventoryRepos: inventoryRepos}
}

func (m *MenuRepo) Create(item models.MenuItem) error {
	for _, items := range m.menuMap {
		if items.ID == item.ID {
			return errorHandle.ItemIdExists
		}
		if items.Name == item.Name {
			return errorHandle.ItemNameExists
		}
	}

	m.menuMap = append(m.menuMap, item)

	updatedData, err := json.MarshalIndent(m.menuMap, "", "  ")
	if err != nil {
		return errorHandle.ErrorFormatJson
	}

	err = ioutil.WriteFile(m.path, updatedData, 0o644)
	if err != nil {
		return errorHandle.ServerError
	}

	return nil
}

func (m *MenuRepo) GetAll() ([]models.MenuItem, error) {
	if len(m.menuMap) == 0 {
		return nil, errorHandle.EmptyFile
	}
	return m.menuMap, nil
}

func (m *MenuRepo) GetItem(id string) (models.MenuItem, error) {
	var item models.MenuItem
	for _, items := range m.menuMap {
		if items.ID == id {
			item = items
			return item, nil
		}
	}
	return item, errorHandle.NotFoundID
}

func (m *MenuRepo) Update(item models.MenuItem, id string) error {
	for _, items := range m.menuMap {
		if items.Name == item.Name && items.ID != id {
			return errorHandle.ItemNameExists
		}
	}

	for i := range m.menuMap {
		if m.menuMap[i].ID == id {
			m.menuMap[i] = item
		}
	}
	updatedData, err := json.MarshalIndent(m.menuMap, "", "  ")
	if err != nil {
		return errorHandle.ErrorFormatJson
	}
	err = ioutil.WriteFile(m.path, updatedData, 0o644)
	if err != nil {
		return errorHandle.ServerError
	}
	return nil
}

func (m *MenuRepo) Delete(id string) error {
	FoundID := false

	for _, item := range m.menuMap {
		if item.ID == id {
			FoundID = true
		}
	}
	if !FoundID {
		return errorHandle.NotFoundID
	}

	var menu []models.MenuItem

	for _, item := range m.menuMap {
		if item.ID != id {
			menu = append(menu, item)
		}
	}
	m.menuMap = menu

	updatedData, err := json.MarshalIndent(m.menuMap, "", "  ")
	if err != nil {
		return errorHandle.ErrorFormatJson
	}
	err = ioutil.WriteFile(m.path, updatedData, 0o644)
	if err != nil {
		return errorHandle.ServerError
	}
	return nil
}

func (m *MenuRepo) ExistsByID(id string) bool {
	for _, item := range m.menuMap {
		if item.ID == id {
			return true
		}
	}
	return false
}

func (m *MenuRepo) MenuCalcuation(id string, quantity float64) error {
	for _, itemMenu := range m.menuMap {
		if itemMenu.ID == id {
			for _, itemIng := range itemMenu.Ingredients {
				if yes := m.inventoryRepos.Calculation(itemIng.IngredientID, itemIng.Quantity*quantity); !yes {
					return errorHandle.Ingred
				}
			}
		}
	}
	return nil
}

func (m *MenuRepo) MenuConsumptionOfIngredients(id string, quantity float64, plus bool) error {
	for _, itemMenu := range m.menuMap {
		if itemMenu.ID == id {
			for _, itemIng := range itemMenu.Ingredients {
				if err := m.inventoryRepos.ConsumptionOfIngredients(itemIng.IngredientID, itemIng.Quantity*quantity, plus); err != nil {
					return errorHandle.ServerError
				}
			}
		}
	}
	return nil
}

func (m *MenuRepo) SumOfOrder(id string) (float64, error) {
	for _, itemMenu := range m.menuMap {
		if id == itemMenu.ID {
			return itemMenu.Price, nil
		}
	}
	return 0, errorHandle.ItemIdExists
}
