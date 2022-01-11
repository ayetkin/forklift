package functions

import (
	log "github.com/sirupsen/logrus"
	"time"
)

func NewWait(stop func(chan bool), message string) {
	done := make(chan bool)
	go func() {
		for {
			if IsWaitClosed(done) {
				return
			}
			<-time.After(5 * time.Second)
			stop(done)
			log.Infof("Waiting for %s...", message)
		}
	}()

	select {
	case err := <-done:
		log.Warnf("Wait stopped, %s. %v", message, err)
	}
}

func IsWaitClosed(ch <-chan bool) bool {
	select {
	case <-ch:
		return true
	default:
	}
	return false
}
