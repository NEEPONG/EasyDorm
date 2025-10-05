package controller

import (
	"html/template"
	"net/http"
	"strconv"

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
		CurrenctRevenue:     GetCurrentMonthRevenue(),
		BillDataSummary:     GetBillingSummary(),
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

func TenantSearch(w http.ResponseWriter, req *http.Request) {
	query := req.FormValue("q")

	var members []model.MemberData = SearchTenants(query)
	tmpl, err := template.ParseFiles("view/html/tenants.html")
	if err != nil {
		renderError(w, "เกิดข้อผิดพลาดในการโหลดหน้าจัดการผู้เช่า")
		return
	}
	tmpl.Execute(w, members)
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

func MaintenanceHandler(w http.ResponseWriter, req *http.Request) {
	var listMaintenance []model.MaintenanceData = GetAllMaintenanceRequests()
	tmpl, err := template.ParseFiles("view/html/maintenance.html")
	if err != nil {
		renderError(w, "เกิดข้อผิดพลาดในการโหลดหน้าบริการซ่อมบำรุง")
		return
	}
	tmpl.Execute(w, listMaintenance)
}

func AddRoomPageHandler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles("view/html/api/addRooms.html")
	if err != nil {
		renderError(w, "เกิดข้อผิดพลาดในการโหลดหน้าบริการซ่อมบำรุง")
		return
	}
	tmpl.Execute(w, nil)
}

func AddRoom(w http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	roomNumber, _ := strconv.Atoi(req.FormValue("room_number"))
	roomFloor, _ := strconv.Atoi(req.FormValue("room_floor"))
	roomType := req.FormValue("room_type")
	price, _ := strconv.ParseFloat(req.FormValue("room_price"), 64)

	err := InsertRoom(roomNumber, roomFloor, roomType, price)
	if err != nil {
		renderError(w, "เกิดข้อผิดพลาดในการเพิ่มห้องพัก: "+err.Error())
		return
	}
	http.Redirect(w, req, "/rooms", http.StatusSeeOther)
}

func SearchRoomHandler(w http.ResponseWriter, req *http.Request) {
	searchId, _ := strconv.Atoi(req.FormValue("q"))
	mapRoomData := SearchRooms(strconv.Itoa(searchId))
	tmpl, err := template.ParseFiles("view/html/rooms.html")
	if err != nil {
		renderError(w, "เกิดข้อผิดพลาดในการโหลดหน้าจัดการห้องพัก")
		return
	}
	tmpl.Execute(w, mapRoomData)
}

func RoomEditPage(w http.ResponseWriter, req *http.Request) {
	roomId, _ := strconv.Atoi(req.FormValue("roomId"))

	room := GetRoomById(roomId)
	tmpl, err := template.ParseFiles("view/html/api/editRooms.html")
	if err != nil {
		renderError(w, "เกิดข้อผิดพลาดในการโหลดหน้าจัดการห้องพัก")
		return
	}

	tmpl.Execute(w, room)
}

func RoomEditHandler(w http.ResponseWriter, req *http.Request) {
	roomId, _ := strconv.Atoi(req.FormValue("room_number"))
	roomFloor, _ := strconv.Atoi(req.FormValue("room_floor"))
	roomType := req.FormValue("room_type")
	roomPrice, _ := strconv.ParseFloat(req.FormValue("room_price"), 64)

	Data := model.RoomData{
		RoomId:    roomId,
		RoomFloor: roomFloor,
		RoomType:  roomType,
		Price:     roomPrice,
	}

	UpdateRoomById(Data)

	http.Redirect(w, req, "/rooms", http.StatusSeeOther)
}

func AddMemberPage(w http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles("view/html/api/addMember.html")
	if err != nil {
		renderError(w, "เกิดข้อผิดพลาดในการโหลดหน้าจัดการห้องพัก")
		return
	}

	tmpl.Execute(w, nil)
}

func AddMemberHandler(w http.ResponseWriter, req *http.Request) {
	memName := req.FormValue("member_name")
	memTel := req.FormValue("member_tel")
	memRoom, _ := strconv.Atoi(req.FormValue("member_room"))

	member := model.MemberData{
		MemberId:   0,
		MemberName: memName,
		MemberTel:  memTel,
		MemberRoom: memRoom,
	}

	InsertMember(member)

	http.Redirect(w, req, "/tenants", http.StatusSeeOther)
}
