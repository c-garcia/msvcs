package tools

import (
	"log"
)

func WarnOnError(msg string, err error) {
	if err != nil {
		log.Println(msg)
	}
}

func FailOnError(msg string, err error) {
	if err != nil {
		log.Fatal(msg)
	}
}
