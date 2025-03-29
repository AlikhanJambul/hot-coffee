package handler

import (
	"encoding/json"
	"hot-coffee/internal/service"
	"hot-coffee/models"
	"log/slog"
	"net/http"
)

type OrderHandler struct {
	service service.OrderService
}

func NewOrderHandler(service service.OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

func (o *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slog.Info("Request CreateOrder")

	var order models.Order

	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		slog.Warn(err.Error())
		JsonWriter(w, 500, "", err)
		return
	}

	err = o.service.Create(order)
	if err != nil {
		slog.Warn(err.Error())
		JsonWriter(w, 500, "", err)
		return
	}
	JsonWriter(w, 200, "Order has been create succesful", nil)
}

func (o *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		slog.Warn("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path != "/orders" {
		slog.Warn("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	slog.Info("Request GetAllOrders")

	orders, err := o.service.GetAll()
	if err != nil {
		slog.Warn(err.Error())
		JsonWriter(w, 500, "", err)
		return
	}
	JsonWriterListInventory(w, 200, nil, nil, orders, nil)
}

func (o *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")
	order, err := o.service.GetItem(id)
	if err != nil {
		slog.Warn(err.Error())
		JsonWriter(w, 500, "", err)
		return
	}
	JsonWriterItemForInventory(w, 200, nil, nil, &order, nil)
}

func (o *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")

	var newOrder models.Order

	err := json.NewDecoder(r.Body).Decode(&newOrder)
	if err != nil {
		slog.Warn(err.Error())
		JsonWriter(w, 500, "", err)
		return
	}

	err = o.service.Update(newOrder, id)
	if err != nil {
		slog.Warn(err.Error())
		JsonWriter(w, 500, "", err)
		return
	}
	JsonWriter(w, 200, "Order has been updated succesful", nil)
}

func (o *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")
	err := o.service.Delete(id)
	if err != nil {
		slog.Warn(err.Error())
		JsonWriter(w, 500, "", err)
		return
	}
	JsonWriter(w, 200, "Order has been deleted succesful", nil)
}

func (o *OrderHandler) StatusClose(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		slog.Warn("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.PathValue("id")
	err := o.service.UpdateStatus(id)
	if err != nil {
		slog.Warn(err.Error())
		JsonWriter(w, 500, "", err)
		return
	}
	JsonWriter(w, 200, "Order has been updated succesful", nil)
}

func (o *OrderHandler) TotalSales(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	totalSum, err := o.service.TotalSum()
	if err != nil {
		slog.Warn(err.Error())
		JsonWriter(w, 500, "", err)
		return
	}
	TotalSalesAndPopularItemResponse(w, 200, totalSum, "")
}

func (o *OrderHandler) TheMostPopularItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	popularItem, err := o.service.MostPopularItem()
	if err != nil {
		slog.Warn(err.Error())
		JsonWriter(w, 500, "", err)
		return
	}
	TotalSalesAndPopularItemResponse(w, 200, 0, popularItem)
}
