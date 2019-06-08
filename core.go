package main

import (
	"database/sql"
	"log"
	"time"
)

func (b *Booking) getStatus() string {
	if b.alreadyLendout() {
		if b.alreadyReturned() {
			return StatusReturned
		}
		return StatusLending
	}
	if b.ableLendout() {
		return StatusAbleLendout
	} else if b.Until.Before(time.Now()) {
		return StatusOverdue
	} else {
		return StatusNotAbleLendout
	}
}

func (b *Booking) alreadyLendout() bool {
	/*
		Return whether booking is all lending out
	*/
	stmt, err := db.Prepare(`
	SELECT COUNT(1)
	FROM Records R, Devices D
	WHERE R.Booking = ? and D.ID = R.Device and D.Type = ?;
	`)
	checkErr(err, "Prepare count booking record fatal: ")
	for i, t := range itemsType {
		row := stmt.QueryRow(b.ID, t)
		var count int
		if err := row.Scan(&count); err == sql.ErrNoRows {
			return false
		} else if err != nil {
			log.Fatalln(err)
		}
		if count < b.Devices[i] {
			return false
		}
	}
	return true
}

func (b *Booking) alreadyReturned() bool {
	if !b.alreadyLendout() {
		return false
	}
	var v int
	row := db.QueryRow(`
		SELECT Amount FROM UnDoneBookings WHERE Booking = ?;
	`, b.ID)

	if err := row.Scan(&v); err == sql.ErrNoRows {
		return true
	} else if err != nil {
		log.Fatalln("func alreadyReturned: scan error: ", err)
	}
	return false
}

func (b *Booking) ableLendout() bool {
	/*
		Return this booking is able to lendout or not
		Check whether every booking before it after now are already lended
	*/
	if b.Until.Before(time.Now()) {
		// 已經到了歸還時間，不須借
		return false
	}
	if b.alreadyLendout() {
		// 已經全部借出，不需借
		return false
	}
	// 檢查在此事件之前的所有預約是否已經借出，如果否，則不能借出。
	rows, err := db.Query(`
			SELECT ID
			FROM Bookings
			WHERE ReturnTime > ? and ReturnTime < ?;
			`, time.Now(), b.Until)
	checkErr(err, "func ableLendout: Query fatal: ")
	for rows.Next() {
		var id int64
		checkErr(rows.Scan(&id), "func ableLendout: scan fatal: ")
		booking := Booking{
			ID:      id,
			Devices: getBookingDevices(id),
		}
		if !booking.alreadyLendout() {
			return false
		}
	}
	return true
}
