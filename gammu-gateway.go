// gammu-gateway
package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"unicode/utf8"
)

func main() {
	http.HandleFunc("/send", sendSmsHandler)
	log.Fatal(http.ListenAndServe(":5137", nil))
}

var sendSmsHandler = func(w http.ResponseWriter, r *http.Request) {
	phoneNumber := r.FormValue("phone_number")
	message := r.FormValue("message")
	if utf8.RuneCountInString(phoneNumber) != 11 {
		fmt.Fprintf(w, "Phone number with invalid length.")
		return
	}
	if !strings.HasPrefix(phoneNumber, "+") {
		phoneNumber = fmt.Sprint("+86", phoneNumber)
	}
	fmt.Println(phoneNumber, message)
	go sendSms(phoneNumber, message)
	fmt.Fprintf(w, "OK.")
}

var sendSms = func(phoneNumber string, message string) {
	command := fmt.Sprint("/usr/bin/gammu sendsms TEXT ", phoneNumber, " -unicode -text '", message, "'")
	out, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", out)
}
