package dal

import (
	"encoding/json"
	"fmt"
	"hot-coffee/internal/errorHandle"
	"hot-coffee/models"
	"io/ioutil"
	"time"
)

var (
	Directory string
	TotalSum  uint
	counter   int // Глобальный счётчик

)

type OrderRepository interface {
	Create(order models.Order) error
	GetAll() ([]models.Order, error)
	GetItem(id string) (models.Order, error)
	Update(order models.Order, id string) error
	Delete(id string) error
	UpdateStatus(id string) error
}

type OrderRepo struct {
	path     string
	orderMap []models.Order
}

func NewOrderRepository(path string) OrderRepository {
	existingData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}

	var order []models.Order

	if len(existingData) > 0 {
		err = json.Unmarshal(existingData, &order)
		if err != nil {
			return nil
		}
	}

	return &OrderRepo{path: path, orderMap: order}
}

func (o *OrderRepo) Create(order models.Order) error {
	order.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	order.Status = "Open"
	order.ID = GenerateOrderCode()

	o.orderMap = append(o.orderMap, order)

	updatedData, err := json.MarshalIndent(o.orderMap, "", "  ")
	if err != nil {
		return errorHandle.ErrorFormatJson
	}

	err = ioutil.WriteFile(o.path, updatedData, 0o644)
	if err != nil {
		return errorHandle.ServerError
	}

	return nil
}

func (o *OrderRepo) GetAll() ([]models.Order, error) {
	if len(o.orderMap) == 0 {
		return nil, errorHandle.EmptyFile
	}
	return o.orderMap, nil
}

func (o *OrderRepo) GetItem(id string) (models.Order, error) {
	for _, item := range o.orderMap {
		if item.ID == id {
			return item, nil
		}
	}
	return models.Order{}, errorHandle.NotFoundID
}

func (o *OrderRepo) Update(order models.Order, id string) error {
	order.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	order.ID = id

	for item := range o.orderMap {
		if o.orderMap[item].ID == order.ID {
			o.orderMap[item] = order
		}
	}

	updatedData, err := json.MarshalIndent(o.orderMap, "", "  ")
	if err != nil {
		return errorHandle.ErrorFormatJson
	}

	err = ioutil.WriteFile(o.path, updatedData, 0o644)
	if err != nil {
		return errorHandle.ServerError
	}

	return nil
}

func (o *OrderRepo) Delete(id string) error {
	var newOrders []models.Order
	for _, item := range o.orderMap {
		if item.ID != id {
			newOrders = append(newOrders, item)
		}
	}
	o.orderMap = newOrders

	updatedData, err := json.MarshalIndent(o.orderMap, "", "  ")
	if err != nil {
		return errorHandle.ErrorFormatJson
	}

	err = ioutil.WriteFile(o.path, updatedData, 0o644)
	if err != nil {
		return errorHandle.ServerError
	}

	return nil
}

func (o *OrderRepo) UpdateStatus(id string) error {
	if exists := CheckId(id, o.orderMap); !exists {
		return errorHandle.NotFoundID
	}

	for item := range o.orderMap {
		if o.orderMap[item].ID == id {
			if o.orderMap[item].Status == "Close" {
				return errorHandle.StatusExists
			}

			o.orderMap[item].Status = "Close"
		}
	}

	updatedData, err := json.MarshalIndent(o.orderMap, "", "  ")
	if err != nil {
		return errorHandle.ErrorFormatJson
	}

	err = ioutil.WriteFile(o.path, updatedData, 0o644)
	if err != nil {
		return errorHandle.ServerError
	}

	return nil
}

func GenerateOrderCode() string {
	counter++ // Увеличиваем счётчик
	return fmt.Sprintf("order%d", counter)
}

func CheckId(id string, orders []models.Order) bool {
	for _, item := range orders {
		if item.ID == id {
			return true
		}
	}
	return false
}
