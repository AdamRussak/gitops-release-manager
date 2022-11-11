package core

import (
	log "github.com/sirupsen/logrus"
)

func OnErrorFail(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %s\n", message, err)
	}
}
