package controller

import (
	"html/template"
	"net/http"

	model "dormitorymng/model"
)

func renderError(w http.ResponseWriter, message string) {
	tmpl, err := template.ParseFiles("view/html/error.html")
	if err != nil {
		http.Error(w, "เกิดข้อผิดพลาด", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, map[string]string{"Message": message})
}

func LoginHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		// ...ตรวจสอบ login และ set cookie...
		http.Redirect(w, req, "/dashboard", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("view/html/login.html")
	if err != nil {
		renderError(w, "เกิดข้อผิดพลาดในการโหลดหน้าเข้าสู่ระบบ")
		return
	}
	tmpl.Execute(w, nil)
}

func DashboardHandler(w http.ResponseWriter, req *http.Request) {
	Data := model.DashBoardData{
		EmptyRoomCount:      GetEmptyRoomCount(),
		TotalOwed:           GetOutstandingPayments(),
		MaintenanceReqCount: GetMaintenanceReqCount(),
		TotalMembers:        GetToTalMembers(),
		RoomCount:           GetAllRooms(),
	}

	tmpl, err := template.ParseFiles("view/html/dashboard.html")
	if err != nil {
		renderError(w, "เกิดข้อผิดพลาดในการโหลดหน้าแดชบอร์ด")
		return
	}
	tmpl.Execute(w, Data)
}

func RoomManagementHandler(w http.ResponseWriter, req *http.Request) {
	var listRoom map[string][]model.RoomData = GetMapRoomData()
	tmpl, err := template.ParseFiles("view/html/rooms.html")
	if err != nil {
		renderError(w, "เกิดข้อผิดพลาดในการโหลดจัดการห้องพัก")
		return
	}
	tmpl.Execute(w, listRoom)
}

func TenantManagementHandler(w http.ResponseWriter, req *http.Request) {
	var memberList []model.MemberData = GetAllMembers()
	tmpl, err := template.ParseFiles("view/html/tenants.html")
	if err != nil {
		renderError(w, "เกิดข้อผิดพลาดในการโหลดหน้าจัดการผู้เช่า")
		return
	}
	tmpl.Execute(w, memberList)
}

func BillingHandler(w http.ResponseWriter, req *http.Request) {
	var listPayment map[string][]model.PaymentData = GetPaymentData()
	tmpl, err := template.ParseFiles("view/html/billing.html")
	if err != nil {
		renderError(w, "เกิดข้อผิดพลาดในการโหลดหน้าบิลค่าเช่า")
		return
	}
	tmpl.Execute(w, listPayment)
}
