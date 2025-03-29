package main

import (
	"fmt"
	"hot-coffee/internal/dal"
	"hot-coffee/internal/handler"
	"hot-coffee/internal/service"
	"hot-coffee/internal/start"
	"log/slog"
	"net/http"
	"strconv"
)

func main() {
	portFlag, dirFlag := start.AllFlags()
	port := strconv.Itoa(portFlag)
	start.CreateDir(dirFlag)
	dal.Directory = dirFlag

	inventoryRepo := dal.NewInventoryRepository(dirFlag + "/inventory.json")
	inventoryService := service.NewInventoryService(inventoryRepo)
	inventoryHandler := handler.NewInventoryHandler(inventoryService)

	menuRepo := dal.NewMenuRepository(dirFlag+"/menu_items.json", inventoryRepo)
	menuService := service.NewMenuService(menuRepo)
	menuHandler := handler.NewMenuHandler(menuService)

	orderRepo := dal.NewOrderRepository(dirFlag + "/orders.json")
	orderService := service.NewOrderService(orderRepo, menuRepo, inventoryRepo)
	orderHandler := handler.NewOrderHandler(orderService)

	http.HandleFunc("POST /menu", menuHandler.CreateNewMenu)
	http.HandleFunc("GET /menu", menuHandler.GetAllMenu)
	http.HandleFunc("GET /menu/{id}", menuHandler.GetItemMenu)
	http.HandleFunc("PUT /menu/{id}", menuHandler.UpdateMenu)
	http.HandleFunc("DELETE /menu/{id}", menuHandler.DeleteItemFromMenu)

	http.HandleFunc("POST /inventory", inventoryHandler.CreateNewInventory)
	http.HandleFunc("GET /inventory", inventoryHandler.GetAllInventory)
	http.HandleFunc("GET /inventory/{id}", inventoryHandler.GetItemInventory)
	http.HandleFunc("PUT /inventory/{id}", inventoryHandler.UpdateInventory)
	http.HandleFunc("DELETE /inventory/{id}", inventoryHandler.DeleteInventory)

	http.HandleFunc("POST /orders", orderHandler.CreateOrder)
	http.HandleFunc("GET /orders", orderHandler.GetAllOrders)
	http.HandleFunc("GET /orders/{id}", orderHandler.GetOrder)
	http.HandleFunc("PUT /orders/{id}", orderHandler.UpdateOrder)
	http.HandleFunc("DELETE /orders/{id}", orderHandler.DeleteOrder)
	http.HandleFunc("POST /orders/{id}/close", orderHandler.StatusClose)

	http.HandleFunc("GET /reports/total-sales", orderHandler.TotalSales)
	http.HandleFunc("GET /reports/popular-items", orderHandler.TheMostPopularItem)

	slog.Info("Server is running", slog.Int("port", portFlag))
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
