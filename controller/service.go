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
