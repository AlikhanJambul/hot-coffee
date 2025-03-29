package service

import (
	"hot-coffee/internal/dal"
	"hot-coffee/internal/errorHandle"
	"hot-coffee/models"
)

type OrderService interface {
	Create(order models.Order) error
	GetAll() ([]models.Order, error)
	GetItem(id string) (models.Order, error)
	Update(item models.Order, id string) error
	Delete(id string) error
	UpdateStatus(id string) error
	TotalSum() (float64, error)
	MostPopularItem() (string, error)
}

type orderService struct {
	orderRepo     dal.OrderRepository
	menuRepo      dal.MenuRepository
	inventoryRepo dal.InventoryRepository
}

type GetPopularItem struct {
	name     string
	quantity int
}

func NewOrderService(orderRepo dal.OrderRepository, menuRepo dal.MenuRepository, inventoryRepo dal.InventoryRepository) OrderService {
	return &orderService{
		orderRepo:     orderRepo,
		menuRepo:      menuRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (o *orderService) Create(order models.Order) error {
	if order.ID != "" || order.CustomerName == "" || len(order.Items) == 0 {
		return errorHandle.ErrorFormatJson
	}
	for _, item := range order.Items {
		if item.Quantity <= 0 {
			return errorHandle.QuantityLessZero
		}
		if item.ProductID == "" {
			return errorHandle.ErrorFormatJson
		}
	}

	for _, items := range order.Items {
		if exists := o.menuRepo.ExistsByID(items.ProductID); !exists {
			return errorHandle.NotFoundID
		}
	}

	for _, item := range order.Items {
		err := o.menuRepo.MenuCalcuation(item.ProductID, float64(item.Quantity))
		if err != nil {
			return err
		}
	}

	for _, item := range order.Items {
		err := o.menuRepo.MenuConsumptionOfIngredients(item.ProductID, float64(item.Quantity), false)
		if err != nil {
			return err
		}
	}

	err := o.orderRepo.Create(order)
	return err
}

func (o *orderService) GetAll() ([]models.Order, error) {
	result, err := o.orderRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return result, err
}

func (o *orderService) GetItem(id string) (models.Order, error) {
	if id == "" {
		return models.Order{}, errorHandle.ErrorFormatJson
	}
	result, err := o.orderRepo.GetItem(id)
	if err != nil {
		return models.Order{}, err
	}
	return result, err
}

func (o *orderService) Update(order models.Order, id string) error {
	if id != order.ID {
		return errorHandle.ChangeID
	}

	if order.Status != "" {
		return errorHandle.ErrorFormatJson
	}

	if order.CreatedAt != "" {
		return errorHandle.ErrorFormatJson
	}

	oldOrder, err := o.orderRepo.GetItem(id)
	if err != nil {
		return err
	}

	if oldOrder.Status == "Close" {
		return errorHandle.StatusExists
	}

	if oldOrder.CustomerName != order.CustomerName {
		return errorHandle.ChangeName
	}

	for _, items := range oldOrder.Items {
		err := o.menuRepo.MenuConsumptionOfIngredients(items.ProductID, float64(items.Quantity), true)
		if err != nil {
			return err
		}
	}

	if order.CustomerName == "" || len(order.Items) == 0 {
		return errorHandle.ErrorFormatJson
	}
	for _, item := range order.Items {
		if item.Quantity <= 0 {
			return errorHandle.QuantityLessZero
		}
		if item.ProductID == "" {
			return errorHandle.ErrorFormatJson
		}
	}

	for _, items := range order.Items {
		if exists := o.menuRepo.ExistsByID(items.ProductID); !exists {
			return errorHandle.NotFoundID
		}
	}

	for _, item := range order.Items {
		err := o.menuRepo.MenuCalcuation(item.ProductID, float64(item.Quantity))
		if err != nil {
			return err
		}
	}

	for _, item := range order.Items {
		err := o.menuRepo.MenuConsumptionOfIngredients(item.ProductID, float64(item.Quantity), false)
		if err != nil {
			return err
		}
	}

	err = o.orderRepo.Update(order, id)

	return err
}

func (o *orderService) Delete(id string) error {
	order, err := o.orderRepo.GetItem(id)
	if err != nil {
		return err
	}

	if order.Status == "Close" {
		return errorHandle.DeleteOrder
	}

	for _, items := range order.Items {
		err := o.menuRepo.MenuConsumptionOfIngredients(items.ProductID, float64(items.Quantity), true)
		if err != nil {
			return err
		}
	}

	err = o.orderRepo.Delete(id)
	return err
}

func (o *orderService) UpdateStatus(id string) error {
	err := o.orderRepo.UpdateStatus(id)
	return err
}

func (o *orderService) TotalSum() (float64, error) {
	var sum float64
	orders, err := o.orderRepo.GetAll()
	if err != nil {
		return 0, err
	}
	for _, orderItem := range orders {
		for _, item := range orderItem.Items {
			priceItem, err := o.menuRepo.SumOfOrder(item.ProductID)
			if err != nil {
				return 0, err
			}
			sum += priceItem * float64(item.Quantity)
		}
	}
	return sum, nil
}

func (o *orderService) MostPopularItem() (string, error) {
	popularItem := GetPopularItem{
		name:     "",
		quantity: 0,
	}
	orders, err := o.GetAll()
	if err != nil {
		return "", err
	}
	for _, order := range orders {
		for _, item := range order.Items {
			if item.Quantity > popularItem.quantity {
				popularItem.name = item.ProductID
				popularItem.quantity = item.Quantity
			}
		}
	}
	return popularItem.name, nil
}
