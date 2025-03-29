package dal

import (
	"encoding/json"
	"hot-coffee/internal/errorHandle"
	"hot-coffee/models"
	"io/ioutil"
)

type InventoryRepository interface {
	Create(item models.InventoryItem) error
	GetAll() ([]models.InventoryItem, error)
	GetItem(id string) (models.InventoryItem, error)
	Update(item models.InventoryItem, id string) error
	Delete(id string) error
	Calculation(id string, quantity float64) bool
	ConsumptionOfIngredients(id string, quantity float64, plus bool) error
}

type inventoryRepo struct {
	path             string
	inventoryMap     []models.InventoryItem
	copyInventoryMap []models.InventoryItem
}

func NewInventoryRepository(path string) InventoryRepository {
	existingData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}

	var inventory []models.InventoryItem
	var copyInventoryMap []models.InventoryItem

	if len(existingData) > 0 {
		err = json.Unmarshal(existingData, &inventory)
		if err != nil {
			return nil
		}
		err = json.Unmarshal(existingData, &copyInventoryMap)
		if err != nil {
			return nil
		}
	}

	return &inventoryRepo{path: path, inventoryMap: inventory, copyInventoryMap: copyInventoryMap}
}

func (i *inventoryRepo) Create(item models.InventoryItem) error {
	for _, itemsInventory := range i.inventoryMap {
		if item.IngredientID == itemsInventory.IngredientID {
			return errorHandle.ItemIdExists
		}
		if item.Name == itemsInventory.Name {
			return errorHandle.ItemNameExists
		}
	}
	i.inventoryMap = append(i.inventoryMap, item)
	i.copyInventoryMap = append(i.copyInventoryMap, item)

	updatedData, err := json.MarshalIndent(i.inventoryMap, "", "  ")
	if err != nil {
		return errorHandle.ErrorFormatJson
	}

	err = ioutil.WriteFile(i.path, updatedData, 0o644)
	if err != nil {
		return errorHandle.ServerError
	}

	return nil
}

func (i *inventoryRepo) GetAll() ([]models.InventoryItem, error) {
	if len(i.inventoryMap) == 0 {
		return nil, errorHandle.EmptyFile
	}
	return i.inventoryMap, nil
}

func (i *inventoryRepo) GetItem(id string) (models.InventoryItem, error) {
	var item models.InventoryItem
	for _, items := range i.inventoryMap {
		if items.IngredientID == id {
			item = items
			return item, nil
		}
	}
	return item, errorHandle.NotFoundID
}

func (i *inventoryRepo) Update(item models.InventoryItem, id string) error {
	for _, items := range i.inventoryMap {
		if items.Name == item.Name && items.IngredientID != id {
			return errorHandle.ItemNameExists
		}
	}

	for items := range i.inventoryMap {
		if i.inventoryMap[items].IngredientID == id {
			i.inventoryMap[items] = item
			i.copyInventoryMap[items] = item
		}
	}

	updatedData, err := json.MarshalIndent(i.inventoryMap, "", "  ")
	if err != nil {
		return errorHandle.ErrorFormatJson
	}

	err = ioutil.WriteFile(i.path, updatedData, 0o644)
	if err != nil {
		return errorHandle.ServerError
	}

	return nil
}

func (i *inventoryRepo) Delete(id string) error {
	FoundID := false
	for _, items := range i.inventoryMap {
		if items.IngredientID == id {
			FoundID = true
		}
	}

	if !FoundID {
		return errorHandle.NotFoundID
	}

	var newInventory []models.InventoryItem

	for _, items := range i.inventoryMap {
		if items.IngredientID != id {
			newInventory = append(newInventory, items)
		}
	}

	i.inventoryMap = newInventory
	i.copyInventoryMap = newInventory

	updatedData, err := json.MarshalIndent(i.inventoryMap, "", "  ")
	if err != nil {
		return errorHandle.ErrorFormatJson
	}

	err = ioutil.WriteFile(i.path, updatedData, 0o644)
	if err != nil {
		return errorHandle.ServerError
	}

	return nil
}

func (i *inventoryRepo) Calculation(id string, quantity float64) bool {
	yes := false
	for _, item := range i.copyInventoryMap {
		if item.IngredientID == id {
			yes = true
		}
	}

	if yes == false {
		return false
	}

	for item := range i.copyInventoryMap {
		if i.copyInventoryMap[item].IngredientID == id {
			if i.copyInventoryMap[item].Quantity-quantity < 0 {
				existingData, err := ioutil.ReadFile(i.path)
				if err != nil {
					return false
				}
				err = json.Unmarshal(existingData, &i.copyInventoryMap)
				if err != nil {
					return false
				}

				return false
			}
			i.copyInventoryMap[item].Quantity -= quantity
		}
	}
	return true
}

func (i *inventoryRepo) ConsumptionOfIngredients(id string, quantity float64, plus bool) error {
	if plus == false {
		for item := range i.inventoryMap {
			if i.inventoryMap[item].IngredientID == id {
				i.inventoryMap[item].Quantity = i.inventoryMap[item].Quantity - quantity
			}
		}
	} else {
		for item := range i.inventoryMap {
			if i.inventoryMap[item].IngredientID == id {
				i.inventoryMap[item].Quantity += quantity
			}
		}
	}

	updatedData, err := json.MarshalIndent(i.inventoryMap, "", "  ")
	if err != nil {
		return errorHandle.ErrorFormatJson
	}

	err = ioutil.WriteFile(i.path, updatedData, 0o644)
	if err != nil {
		return errorHandle.ServerError
	}
	return nil
}
