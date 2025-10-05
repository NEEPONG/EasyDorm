package main

import (
	controller "dormitorymng/controller"
	"net/http"
)

func main() {
	http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir("view"))))
	http.HandleFunc("/", controller.LoginHandler)
	http.HandleFunc("/doLogin", controller.LoginHandler)
	http.HandleFunc("/rooms", controller.RoomManagementHandler)
	http.HandleFunc("/tenants", controller.TenantManagementHandler)
	http.HandleFunc("/billing", controller.BillingHandler)
	http.HandleFunc("/dashboard", controller.DashboardHandler)
	http.ListenAndServe(":8090", nil)
}
