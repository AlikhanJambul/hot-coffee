package handler

import (
	"encoding/json"
	"fmt"
	"hot-coffee/internal/errorHandle"
	"hot-coffee/internal/service"
	"hot-coffee/models"
	"log/slog"
	"net/http"
)

type InventoryHandler struct {
	service service.InventoryService
}

func NewInventoryHandler(service service.InventoryService) *InventoryHandler {
	return &InventoryHandler{
		service: service,
	}
}

func (h *InventoryHandler) CreateNewInventory(w http.ResponseWriter, r *http.Request) {
	slog.Info("Request CreateNewInventory")
	var inventory models.InventoryItem
	err := json.NewDecoder(r.Body).Decode(&inventory)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		slog.Warn(err.Error())
		JsonWriter(w, 500, "", err)
		return
	}

	err = h.service.Create(inventory)

	if err != nil {
		slog.Warn(err.Error())
		JsonWriter(w, 500, "", err)
		return
	}

	JsonWriter(w, 201, "Item added successfully", nil)
}

func (h *InventoryHandler) GetAllInventory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		slog.Warn("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path != "/inventory" {
		slog.Warn("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	slog.Info("Request GetAllInventory")
	allInventory, err := h.service.GetAll()
	if err != nil {
		slog.Warn(err.Error())
		JsonWriter(w, 500, "", err)
		return
	}
	JsonWriterListInventory(w, 200, allInventory, nil, nil, nil)
}

func (h *InventoryHandler) GetItemInventory(w http.ResponseWriter, r *http.Request) {
	slog.Info("Request GetItemInventory")
	id := r.PathValue("id")
	w.Header().Set("Content-Type", "application/json")
	item, err := h.service.GetItem(id)
	if err != nil {
		slog.Warn(err.Error())
		JsonWriter(w, 500, "", err)
		return
	}
	JsonWriterItemForInventory(w, 200, &item, nil, nil, nil)
}

func (h *InventoryHandler) UpdateInventory(w http.ResponseWriter, r *http.Request) {
	slog.Info("Request UpdateInventory")
	var newInventory models.InventoryItem
	id := r.PathValue("id")
	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&newInventory)
	if err != nil {
		slog.Warn(err.Error())
		JsonWriter(w, 500, "", err)
		return
	}

	err = h.service.Update(newInventory, id)

	if err != nil {
		slog.Warn(err.Error())
		JsonWriter(w, 500, "", err)
		return
	}

	JsonWriter(w, 201, "Item updated successfully", nil)
}

func (h *InventoryHandler) DeleteInventory(w http.ResponseWriter, r *http.Request) {
	slog.Info("Request DeleteInventory")
	id := r.PathValue("id")
	w.Header().Set("Content-Type", "application/json")

	err := h.service.Delete(id)
	if err != nil {
		slog.Warn(err.Error())
		JsonWriter(w, 500, "", err)
		return
	}

	JsonWriter(w, 200, "Item deleted successfully", nil)
}

func JsonWriter(w http.ResponseWriter, statusCode int, message string, err error) {
	w.Header().Set("Content-Type", "application/json")

	resp := Response{}
	if err != nil {
		statusCode = errorHandle.CheckErrors(err)
		resp.Error = err.Error()
		w.WriteHeader(statusCode)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err), http.StatusInternalServerError)
		}
		return
	}

	resp.Message = message
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err), http.StatusInternalServerError)
	}
}

func JsonWriterListInventory(w http.ResponseWriter, statusCode int, listInventory []models.InventoryItem, listMenu []models.MenuItem, listOrders []models.Order, err error) {
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		statusCode = errorHandle.CheckErrors(err)
		resp := Response{
			Error: err.Error(),
		}

		w.WriteHeader(statusCode)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err), http.StatusInternalServerError)
		}
		return
	}
	resp := GetListItems{}

	if listInventory != nil {
		resp.Inventory = listInventory
	} else if listMenu != nil {
		resp.Menu = listMenu
	} else if listOrders != nil {
		resp.Orders = listOrders
	}

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err), http.StatusInternalServerError)
	}
}

func JsonWriterItemForInventory(w http.ResponseWriter, statusCode int, itemInventory *models.InventoryItem, itemMenu *models.MenuItem, itemOrders *models.Order, err error) {
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		statusCode = errorHandle.CheckErrors(err)
		resp := Response{
			Error: err.Error(),
		}
		w.WriteHeader(statusCode)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err), http.StatusInternalServerError)
		}
		return
	}

	resp := GetItems{}

	if itemInventory != nil {
		resp.ItemInventory = itemInventory
	} else if itemMenu != nil {
		resp.ItemMenu = itemMenu
	} else if itemOrders != nil {
		resp.ItemOrders = itemOrders
	}

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err), http.StatusInternalServerError)
	}
}

func TotalSalesAndPopularItemResponse(w http.ResponseWriter, statusCode int, sum float64, popularItem string) {
	w.Header().Set("Content-Type", "application/json")

	resp := Response{}

	if sum == 0 {
		resp.PopularItem = popularItem
	} else {
		resp.TotalSales = sum
	}

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %v", err), http.StatusInternalServerError)
	}
}
