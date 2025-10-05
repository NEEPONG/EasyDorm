package controller

import (
	"database/sql"
	model "dormitorymng/model"

	_ "github.com/go-sql-driver/mysql"
)

func getDb() *sql.DB {
	db, err := sql.Open("mysql", "root:1234@tcp(localhost:3307)/goproject")
	if err != nil {
		panic(err.Error()) // หรือจะ handle error แบบอื่นก็ได้
	}
	return db
}

func GetEmptyRoomCount() int {
	db := getDb()

	var count int
	err := db.QueryRow(`
        SELECT COUNT(r.roomId)
        FROM rooms r
        LEFT JOIN member m ON r.roomId = m.memberRoom
        WHERE m.memberRoom IS NULL
    `).Scan(&count)
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()
	return count
}

func GetOutstandingPayments() float64 {
	db := getDb()

	var totalOwed float64
	err := db.QueryRow(`
		SELECT SUM(r.price+p.electricBill+p.waterBill)
		FROM rooms r LEFT JOIN member m ON r.roomId = m.memberRoom
		LEFT JOIN payment p ON r.roomId = p.roomId
		AND DATE_FORMAT(p.payDate, '%Y-%m') = '2025-10'
		WHERE m.memberId IS NOT NULL
		AND p.billStatus = 0;`).Scan(&totalOwed)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	return totalOwed
}

func GetToTalMembers() int {
	db := getDb()

	var count int
	err := db.QueryRow(`SELECT COUNT(memberId) FROM member`).Scan(&count)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	return count
}

func GetMaintenanceReqCount() int {
	db := getDb()
	var count int
	err := db.QueryRow(`SELECT COUNT(roomId) FROM maintenance WHERE status = 'Pending'`).Scan(&count)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	return count
}

func GetAllRooms() int {
	db := getDb()

	var count int
	err := db.QueryRow(`SELECT COUNT(roomId) FROM rooms`).Scan(&count)

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()
	return count
}

func GetMapRoomData() map[string][]model.RoomData {
	db := getDb()
	defer db.Close()

	rowsEmpty, errEmpty := db.Query(`SELECT r.roomId, r.roomFloor, r.price, r.roomType
        FROM rooms r LEFT JOIN member m ON r.roomId = m.memberRoom
        WHERE m.memberRoom IS NULL;`)

	if errEmpty != nil {
		panic(errEmpty.Error())
	}
	defer rowsEmpty.Close()

	var unOccupied []model.RoomData
	for rowsEmpty.Next() {
		var room model.RoomData
		err := rowsEmpty.Scan(&room.RoomId, &room.RoomFloor, &room.Price, &room.RoomType)
		if err != nil {
			panic(err.Error())
		}
		unOccupied = append(unOccupied, room)
	}
	if err := rowsEmpty.Err(); err != nil {
		panic(err.Error())
	}

	rowsOccupied, errOccupied := db.Query(`SELECT r.roomId, r.roomFloor, r.price, r.roomType
        FROM rooms r LEFT JOIN member m ON r.roomId = m.memberRoom
        WHERE m.memberRoom IS NOT NULL;`)

	if errOccupied != nil {
		panic(errOccupied.Error())
	}
	defer rowsOccupied.Close()

	var Occupied []model.RoomData
	for rowsOccupied.Next() {
		var room model.RoomData
		err := rowsOccupied.Scan(&room.RoomId, &room.RoomFloor, &room.Price, &room.RoomType)
		if err != nil {
			panic(err.Error())
		}
		Occupied = append(Occupied, room)
	}
	if err := rowsOccupied.Err(); err != nil {
		panic(err.Error())
	}

	listRoom := map[string][]model.RoomData{
		"true":  Occupied,
		"false": unOccupied,
	}
	return listRoom
}

func GetAllMembers() []model.MemberData {
	db := getDb()
	rows, err := db.Query(`SELECT * FROM member ORDER BY memberRoom`)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	defer rows.Close()

	var members []model.MemberData
	for rows.Next() {
		var member model.MemberData
		err := rows.Scan(&member.MemberId, &member.MemberName, &member.MemberTel, &member.MemberRoom)
		if err != nil {
			panic(err.Error())
		}
		members = append(members, member)
	}
	return members
}

func GetPaymentData() map[string][]model.PaymentData {
	db := getDb()
	defer db.Close()
	paidRows, err := db.Query(`
		SELECT r.roomId, p.payDate, r.price, p.electricBill, p.waterBill, 
		p.billStatus, p.electricBill+p.waterBill+r.price
		FROM rooms r LEFT JOIN member m ON r.roomId = m.memberRoom
		LEFT JOIN payment p ON r.roomId = p.roomId
		AND DATE_FORMAT(p.payDate, '%Y-%m') = '2025-10'
		WHERE m.memberId IS NOT NULL
		AND p.billStatus = 1;
	`)
	if err != nil {
		panic(err.Error())
	}
	defer paidRows.Close()

	var paidPayments []model.PaymentData
	for paidRows.Next() {
		var payment model.PaymentData
		err := paidRows.Scan(&payment.RoomId, &payment.BillDate, &payment.RoomPrice, &payment.ElectricBill,
			&payment.WaterBill, &payment.Status, &payment.PayAmt)
		if err != nil {
			panic(err.Error())
		}
		paidPayments = append(paidPayments, payment)
	}

	unPaidRows, err := db.Query(`
		SELECT r.roomId, p.payDate, r.price, p.electricBill, p.waterBill, 
		p.billStatus, p.electricBill+p.waterBill+r.price
		FROM rooms r LEFT JOIN member m ON r.roomId = m.memberRoom
		LEFT JOIN payment p ON r.roomId = p.roomId
		AND DATE_FORMAT(p.payDate, '%Y-%m') = '2025-10'
		WHERE m.memberId IS NOT NULL
		AND p.billStatus = 1;
	`)
	if err != nil {
		panic(err.Error())
	}
	defer unPaidRows.Close()

	var unPaidPayment []model.PaymentData
	for unPaidRows.Next() {
		var payment model.PaymentData
		err := unPaidRows.Scan(&payment.RoomId, &payment.BillDate, &payment.RoomPrice, &payment.ElectricBill,
			&payment.WaterBill, &payment.Status, &payment.PayAmt)
		if err != nil {
			panic(err.Error())
		}
		unPaidPayment = append(unPaidPayment, payment)
	}

	listPayment := map[string][]model.PaymentData{
		"true":  paidPayments,
		"false": unPaidPayment,
	}
	return listPayment
}

func GetAllMaintenanceRequests() []model.MaintenanceData {
	db := getDb()
	defer db.Close()
	rows, err := db.Query(`
		SELECT m.memberName, ma.roomId, ma.date, ma.text, ma.status FROM maintenance ma
		INNER JOIN member m ON ma.memberId = m.memberId
		ORDER BY ma.status, ma.roomId
	`)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var maintenanceRequests []model.MaintenanceData
	for rows.Next() {
		var request model.MaintenanceData
		err := rows.Scan(&request.MemberName, &request.RoomId, &request.RequestDate, &request.Detail, &request.Status)
		if err != nil {
			panic(err.Error())
		}
		maintenanceRequests = append(maintenanceRequests, request)
	}
	return maintenanceRequests
}

func GetCurrentMonthRevenue() float64 {
	db := getDb()
	defer db.Close()
	var revenue float64
	err := db.QueryRow(`
		SELECT SUM(r.price + p.electricBill + p.waterBill)
		FROM payment p JOIN rooms r ON p.roomId = r.roomId
		WHERE p.billStatus = 1 AND
		YEAR(p.payDate) = YEAR(CURRENT_DATE()) AND 
		MONTH(p.payDate) = MONTH(CURRENT_DATE());
	`).Scan(&revenue)
	if err != nil {
		panic(err.Error())
	}
	return revenue
}

func GetBillingSummary() model.BillDataSummary {
	db := getDb()
	defer db.Close()
	var summary model.BillDataSummary
	err := db.QueryRow(`
		SELECT (SUM(billStatus) * 100.0) / COUNT(*) AS OverallPaymentSuccessRate_Percent,
		SUM(billStatus) AS TotalBillsPaid,
		COUNT(*) AS TotalBillsIssued
		FROM payment
		WHERE
		YEAR(payDate) = YEAR(CURRENT_DATE()) AND 
		MONTH(payDate) = MONTH(CURRENT_DATE());
	`).Scan(&summary.PayPercentage, &summary.PaidAmt, &summary.BillAmt)
	if err != nil {
		panic(err.Error())
	}
	return summary
}

func InsertRoom(roomNumber int, roomFloor int, roomType string, price float64) error {
	db := getDb()
	defer db.Close()

	_, err := db.Exec(`
		INSERT INTO rooms (roomId, roomFloor, roomType, price)
		VALUES (?, ?, ?, ?);
	`, roomNumber, roomFloor, roomType, price)
	if err != nil {
		return err
	}
	return nil
}

func SearchRooms(keyword string) map[string][]model.RoomStatus {
	db := getDb()
	defer db.Close()
	searchTerm := "%" + keyword + "%"

	rows, err := db.Query(`
        SELECT 
		r.roomId, 
		r.roomFloor, 
		r.price, 
		r.roomType,
		IF(COUNT(m.memberRoom) > 0, 'มีผู้เช่า', 'ว่าง') AS RoomStatus 
		FROM 
			rooms r
		LEFT JOIN 
			member m ON r.roomId = m.memberRoom
		WHERE 
			r.roomId LIKE ? OR r.price LIKE ? OR r.roomFloor LIKE ?
		GROUP BY
			r.roomId, r.roomFloor, r.price, r.roomType
		ORDER BY
			r.roomId;
    `, searchTerm, searchTerm, searchTerm)

	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	occupiedRooms := []model.RoomStatus{}
	unoccupiedRooms := []model.RoomStatus{}

	// 4. ลูปเพื่อสแกนและแยกประเภทห้อง
	for rows.Next() {
		var room model.RoomStatus
		var status string // ประกาศตัวแปรเพื่อรับค่าสถานะห้อง

		// ต้องสแกนค่า status ที่ดึงมาจาก SQL Query ด้วย!
		err := rows.Scan(&room.RoomId, &room.RoomFloor, &room.Price, &room.RoomType, &status)
		if err != nil {
			panic(err.Error())
		}

		// กำหนดสถานะลงในโครงสร้าง (ถ้ามี)
		room.IsOccupied = status

		// แยกห้องตามสถานะ
		if status == "มีผู้เช่า" {
			occupiedRooms = append(occupiedRooms, room)
		} else {
			unoccupiedRooms = append(unoccupiedRooms, room)
		}
	}

	if err := rows.Err(); err != nil {
		panic(err.Error())
	}

	// 5. ส่งค่ากลับเป็น map[string][]model.RoomData ตามรูปแบบที่ต้องการ
	listRoom := map[string][]model.RoomStatus{
		"true":  occupiedRooms,
		"false": unoccupiedRooms,
	}

	return listRoom
}

func GetRoomById(id int) model.RoomData {
	db := getDb()
	defer db.Close()
	var room model.RoomData
	err := db.QueryRow(`
		SELECT * FROM rooms WHERE roomId = ?;
	`, id).Scan(&room.RoomId, &room.RoomFloor, &room.Price, &room.RoomType)
	if err != nil {
		panic(err.Error())
	}
	return room
}

func UpdateRoomById(newRoom model.RoomData) error {
	db := getDb()
	defer db.Close()

	query := `
		UPDATE rooms
		SET
			roomFloor = ?,
			roomType  = ?,
			price     = ?
		WHERE roomId = ?
	`

	_, err := db.Exec(query,
		newRoom.RoomFloor,
		newRoom.RoomType,
		newRoom.Price,
		newRoom.RoomId,
	)
	if err != nil {
		panic(err.Error())
	}
	return nil
}

func GetMaxMemberId() int {
	db := getDb()
	defer db.Close()

	var maxId int
	err := db.QueryRow(`
		SELECT MAX(memberId) + 1 FROM member;
	`).Scan(&maxId)
	if err != nil {
		panic(err.Error())
	}
	return maxId

}

func InsertMember(newMember model.MemberData) {
	db := getDb()
	defer db.Close()

	query := `
		INSERT INTO member (memberId, memberName, memberTel, memberRoom) VALUES(?, ?, ?, ?);
	`

	_, err := db.Exec(query,
		GetMaxMemberId(),
		newMember.MemberName,
		newMember.MemberTel,
		newMember.MemberRoom,
	)
	if err != nil {
		panic(err.Error())
	}
}

func SearchTenants(keyword string) []model.MemberData {
	db := getDb()
	defer db.Close()
	searchTerm := "%" + keyword + "%"

	rows, err := db.Query(`
        SELECT * FROM member WHERE memberName LIKE ? OR memberTel LIKE ? OR memberRoom LIKE ?;
    `, searchTerm, searchTerm, searchTerm)

	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	members := []model.MemberData{}

	for rows.Next() {
		var member model.MemberData

		err := rows.Scan(&member.MemberId, &member.MemberName, &member.MemberTel, &member.MemberRoom)
		if err != nil {
			panic(err.Error())
		}

		members = append(members, member)
	}

	return members
}
