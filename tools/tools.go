package tools

import (
	"fmt"
	"log"
	"strings"
)

func CheckIfError(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

func CheckTheEndline(secretData string) {
	if strings.Contains(secretData, "\n") {
		log.Printf(secretData + " Contains n ")

	}
	if strings.Contains(secretData, "\r") {
		fmt.Printf(secretData + " Contains r ")
	}
}
