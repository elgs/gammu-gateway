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
	log.Fatal(http.ListenAndServe("192.168.7.169:5137", nil))
}

var sendSmsHandler = func(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Failed to handle request:", err)
		}
	}()
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
	sendSms(phoneNumber, message)
	fmt.Fprintf(w, "OK.")
}

var sendSms = func(phoneNumber string, message string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Failed to send:", err)
		}
	}()
	command := fmt.Sprint("/usr/bin/gammu sendsms TEXT ", phoneNumber, " -unicode -text '", message, "'")
	out, err := exec.Command("sh", "-c", command).CombinedOutput()
	if err != nil {
		log.Fatal("Failed to execute:", err)
	}
	fmt.Printf("gammu output %s\n", out)
}
