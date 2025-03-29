package handler

import (
	"encoding/json"
	"hot-coffee/internal/service"
	"hot-coffee/models"
	"log/slog"
	"net/http"
)

type Response struct {
	Error       string  `json:"error,omitempty"`
	Message     string  `json:"message,omitempty"`
	TotalSales  float64 `json:"total_sales,omitempty"`
	PopularItem string  `json:"popular_item,omitempty"`
}

type GetListItems struct {
	Menu      []models.MenuItem      `json:"menu,omitempty"`
	Inventory []models.InventoryItem `json:"inventory,omitempty"`
	Orders    []models.Order         `json:"orders,omitempty"`
}

type GetItems struct {
	ItemMenu      *models.MenuItem      `json:"Item_Menu,omitempty"`
	ItemInventory *models.InventoryItem `json:"Item_Inventory,omitempty"`
	ItemOrders    *models.Order         `json:"Item_Orders,omitempty"`
}

type MenuHandler struct {
	service service.MenuService
}

func NewMenuHandler(service service.MenuService) *MenuHandler {
	return &MenuHandler{
		service: service,
	}
}

func (h *MenuHandler) CreateNewMenu(w http.ResponseWriter, r *http.Request) {
	var menu models.MenuItem
	slog.Info("Request CreateNewMenu")
	err := json.NewDecoder(r.Body).Decode(&menu)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		slog.Warn(err.Error())
		JsonWriter(w, 500, "", err)
		return
	}

	err = h.service.Create(menu)

	if err != nil {
		slog.Warn(err.Error())
		JsonWriter(w, 500, "", err)
		return
	}

	JsonWriter(w, 201, "Item added successfully", nil)
}

func (h *MenuHandler) GetAllMenu(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		slog.Warn("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path != "/menu" {
		slog.Warn("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	slog.Info("Request GetAllMenu")
	allMenuItems, err := h.service.GetAll()
	if err != nil {
		slog.Warn(err.Error())
		JsonWriter(w, 500, "", err)
		return
	}

	JsonWriterListInventory(w, 200, nil, allMenuItems, nil, nil)
}

func (h *MenuHandler) GetItemMenu(w http.ResponseWriter, r *http.Request) {
	slog.Info("Request GetItemMenu")
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")
	result, err := h.service.GetItem(id)
	if err != nil {
		slog.Warn(err.Error())
		JsonWriter(w, 500, "", err)
		return
	}
	JsonWriterItemForInventory(w, 200, nil, &result, nil, nil)
}

func (h *MenuHandler) UpdateMenu(w http.ResponseWriter, r *http.Request) {
	slog.Info("Request UpdateMenu")
	var newMenu models.MenuItem
	id := r.PathValue("id")
	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&newMenu)
	if err != nil {
		slog.Warn(err.Error())
		JsonWriter(w, 500, "", err)
		return
	}

	err = h.service.Update(newMenu, id)

	if err != nil {
		slog.Warn(err.Error())
		JsonWriter(w, 500, "", err)
		return
	}

	JsonWriter(w, 201, "Item updated successfully", nil)
}

func (h *MenuHandler) DeleteItemFromMenu(w http.ResponseWriter, r *http.Request) {
	slog.Info("Request DeleteItemFromMenu")
	w.Header().Set("Content-Type", "application/json")
	id := r.PathValue("id")
	err := h.service.Delete(id)
	if err != nil {
		slog.Warn(err.Error())
		JsonWriter(w, 500, "", err)
		return
	}
	JsonWriter(w, 200, "Item deleted successfully", nil)
}
