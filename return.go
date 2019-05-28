package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

func returnDevice(dID string) error {
	result, err := db.Exec(`
		UPDATE Records
		SET LentUntil = CURRENT_TIMESTAMP()
		WHERE LentUntil IS NULL and Device = ?;
	`)
	if err != nil {
		log.Fatalln(err)
	}
	if rowNum, err := result.RowsAffected(); err != nil {
		log.Fatalln(err)
	} else if rowNum == 0 {
		return errors.New("Record Not Found")
	} else if rowNum > 1 {
		return fmt.Errorf("Return device fatal: RowsAffected is %d", rowNum)
	}
	return nil
}

func handleReturnDevice(w http.ResponseWriter, r *http.Request) {
	user := getUser(w, r)
	if user.Type != "Admin" {
		w.WriteHeader(401)
		return
	}
	dID := r.FormValue("dID")
	if err := returnDevice(dID); err.Error() == "Record Not Found" {
		w.WriteHeader(404)
		return
	} else if err != nil {
		w.WriteHeader(403)
		return
	}
}
