package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func newBooking(w http.ResponseWriter, r *http.Request) {
	/*
		Handle /bookings
		Receive post from /bookings/new form
	*/
	user := getUser(w, r)
	if !user.Login {
		permissionDenied(w, r)
		return
	}
	b := Booking{}
	var err error
	b.Devices[student], _ = strconv.Atoi(r.FormValue("student"))
	b.Devices[teacher], _ = strconv.Atoi(r.FormValue("teacher"))
	b.Devices[chromebook], _ = strconv.Atoi(r.FormValue("chromebook"))
	b.Devices[wap], _ = strconv.Atoi(r.FormValue("wap"))
	b.Devices[projector], _ = strconv.Atoi(r.FormValue("wireless-projector"))
	date, err := time.ParseInLocation("2006-01-02", r.FormValue("date"), time.Local)
	checkErr(err, "")
	class, err := strconv.Atoi(r.FormValue("class"))
	checkErr(err, "")
	if class < 0 || class > 10 {
		page := struct {
			User
			Title   string
			Content string
			Target  string
		}{user, "請輸入有效的課堂", "請輸入有效的課堂", `<a href="/bookings/new">回到預約頁面</a>`}
		checkErr(tpl.ExecuteTemplate(w, "msg.html", page), "func newBooking: render fatal: ")
		return
	}
	if msg := b.enough(check(date.Add(classBegin[class]), date.Add(classEnd[class]))); len(msg) != 0 {
		page := struct {
			User
			Title   string
			Content string
			Target  string
		}{user, "數量不足", strings.Join(msg, "<br>"), `<a href="/bookings/new">回到預約頁面</a>`}
		checkErr(tpl.ExecuteTemplate(w, "msg.html", page), "func newBooking: render fatal: ")
		return
	}
	b.From = date.Add(classBegin[class])
	b.Until = date.Add(classEnd[class])
	b.insertBooking(user)
	checkErr(err, "Parse time fatal: ")
	page := struct {
		User
		Booking
		Date      string
		Class     string
		Weekday   string
		ItemsName [5]string
	}{user, b, date.Format("2006-01-02"), className[class], b.From.Weekday().String(), itemsName}
	checkErr(tpl.ExecuteTemplate(w, "examination.html", page), "Execute examination fatal: ")
}

func bookingForm(w http.ResponseWriter, r *http.Request) {
	// Render booking form and error messages
	type bookingPage struct {
		User
		Classes [11]string
		Min     string
		Max     string
	}
	var page bookingPage
	page.User = getUser(w, r)
	now := time.Now()
	max := now.AddDate(0, 1, 0)
	page.Min = now.Format("2006-01-02")
	page.Max = max.Format("2006-01-02")
	page.Classes = className
	err := tpl.ExecuteTemplate(w, "newBooking.html", page)
	if err != nil {
		log.Fatal("Template execute fatal: ", err)
	}
}

func (b *Booking) insertBooking(user User) {
	// insert booking data into database
	log.Println(b.Devices)
	b.returnIfNotEnough()
	result, err := db.Exec(`
	INSERT INTO Bookings (User, LendingTime, ReturnTime)
	VALUES (?, ?, ?);
	`, user.UID, b.From, b.Until)
	checkErr(err, "Insert booking value fatal: ")
	b.ID, err = result.LastInsertId()
	checkErr(err, "Get last insert ID fatal: ")

	stmt, err := db.Prepare(`
	INSERT INTO BookingDevices
	VALUES (?, ?, ?);
	`)
	checkErr(err, "Insert query prepare fatal: ")
	for i, v := range b.Devices {
		if v > 0 {
			_, err := stmt.Exec(b.ID, itemsType[i], v)
			checkErr(err, "Insert booking devices fatal: ")
		}
	}
}

func (b *Booking) enough(r [5]int) []string {
	/*
		Check whether devices are enough for this booking.
		Return slice of error messages. len == 0 if all pass.
	*/
	msg := []string{}
	for i := 0; i < len(itemsName); i++ {
		if r[i] < b.Devices[i] {
			msg = append(msg, fmt.Sprintf("%s 只剩 %d 台", itemsName[i], r[i]))
		}
	}
	return msg
}

func (b *Booking) returnIfNotEnough() (err error) {
	/*
		檢查是否有足夠的館藏量足以借出
		如果否，代表有提早借出之預約
		將最晚的預約還入直到足夠為止
	*/
	stmt, err := db.Prepare(`
	SELECT COUNT(1)
	FROM DevicesStatus
	WHERE Type = ? AND Status = FALSE;
	`)
	if err != nil {
		log.Fatalln(err)
	}
	for i, t := range b.Devices {
		row := stmt.QueryRow(itemsType[i])
		var count int
		if err = row.Scan(&count); err != nil {
			log.Fatalln(err)
		}
		log.Println(count)
		if count <= t {
			err = returnLast(i)
			if err != nil {
				log.Panic(err)
			}
			return b.returnIfNotEnough()
		}
	}
	return nil
}

func returnLast(t int) (err error) {
	/*
		還入已借出預約中，離現在最遠，且包含 t 這個 type 的 records
		僅歸還有包含該 type 的 records
	*/
	row := db.QueryRow(`
	SELECT B.ID
	FROM UnDoneBookings UB
	INNER JOIN Bookings B
	ON UB.ID = B.ID
	INNER JOIN BookingDevices BD
	ON B.ID = BD.BID
	WHERE UB.Amount > 0 AND BD.Type = ?
	ORDER BY LendingTime DESC, ID DESC
	LIMIT 1;
	`, t)
	var bid int64
	if err = row.Scan(&bid); err == sql.ErrNoRows {
		return errors.New("No booking found")
	} else if err != nil {
		log.Fatalln(err)
	}
	log.Printf("%d will be removed\n", bid)
	_, err = db.Exec(`
	DELETE FROM Records R
	INNER JOIN Devices D
	ON R.Device = D.ID
	WHERE Booking = ? AND Type = ?;
	`, bid, t)
	if err != nil {
		log.Fatalln(err)
	}
	return nil
}

func getBookingDevices(id int64) (devices [5]int) {
	rows, err := db.Query(`
	SELECT Type, Amount
	FROM BookingDevices
	WHERE BID = ?
	ORDER BY Type;
	`, id)
	checkErr(err, "Query booking devices fatal: ")
	i := 0
	for rows.Next() {
		var t string
		var amount int
		checkErr(rows.Scan(&t, &amount), "getBookingDevices Scan rows fatal: ")
		for itemsType[i] != t && i < 5 {
			i++
		}
		devices[i] = amount
	}
	return
}

func getBooking(id int64) (Booking, error) {
	b := Booking{}
	row := db.QueryRow(`
		SELECT U.Name, LendingTime, ReturnTime
		FROM Bookings B, Users U
		WHERE B.ID = ? and U.ID = B.User;
	`, id)
	if err := row.Scan(&b.UName, &b.From, &b.Until); err == sql.ErrNoRows {
		return b, ErrBookingNotFound
	} else if err != nil {
		log.Fatalln(err)
	}
	b.ID = id
	b.Devices = getBookingDevices(id)
	return b, nil
}
