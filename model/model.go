package model

type DashBoardData struct {
	RoomCount           int
	EmptyRoomCount      int
	TotalOwed           float64
	MaintenanceReqCount int
	TotalMembers        int
	CurrenctRevenue     float64
	BillDataSummary     BillDataSummary
}

type PaymentData struct {
	RoomId       int
	BillDate     string
	RoomPrice    float64
	ElectricBill float64
	WaterBill    float64
	Status       int
	PayAmt       float64
}

type RoomData struct {
	RoomId    int
	RoomFloor int
	Price     float64
	RoomType  string
}

type MemberData struct {
	MemberId   int
	MemberName string
	MemberTel  string
	MemberRoom int
}

type MaintenanceData struct {
	MemberName  string
	RoomId      int
	RequestDate string
	Detail      string
	Status      string
}

type BillDataSummary struct {
	PayPercentage float64
	PaidAmt       int
	BillAmt       int
}

type RoomStatus struct {
	RoomId     int
	RoomFloor  int
	Price      float64
	RoomType   string
	IsOccupied string
}
