package main

import (
	"testing"
	"time"
)

func TestAlreadyLendOut(t *testing.T) {
	b := Booking{
		Devices: [5]int{15, 2, 30, 0, 1},
		From:    time.Now(),
		Until:   time.Now().Add(time.Minute * 45),
	}
	user := User{
		ID: 1,
	}
	b.insertBooking(user)
	if alreadyLendout(b) {
		t.Error("Wrong result: should get false but get true")
	}
}
