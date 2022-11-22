package tools

import "log"

func CheckIfError(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
