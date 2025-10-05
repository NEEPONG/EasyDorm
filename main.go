package main

import (
	controller "dormitorymng/controller"
	"net/http"
)

func main() {
	http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir("view"))))
	http.HandleFunc("/", controller.LoginHandler)
	http.HandleFunc("/doLogin", controller.LoginHandler)

	// Rooms management
	http.HandleFunc("/rooms", controller.RoomManagementHandler)
	http.HandleFunc("/rooms/add", controller.AddRoomPageHandler)
	http.HandleFunc("/api/rooms/add", controller.AddRoom)
	http.HandleFunc("/rooms/search", controller.SearchRoomHandler)

	// Mantenance management
	http.HandleFunc("/maintenance", controller.MaintenanceHandler)

	// Tenant management
	http.HandleFunc("/tenants", controller.TenantManagementHandler)

	// Billing management
	http.HandleFunc("/billing", controller.BillingHandler)

	// Dashboard management
	http.HandleFunc("/dashboard", controller.DashboardHandler)
	http.ListenAndServe(":8090", nil)
}
