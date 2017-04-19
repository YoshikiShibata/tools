package trace

import (
	"log"
	"time"
)

// Trace can be used to log  entering and exiting of a func
func Trace(msg string) func() {
	start := time.Now()
	log.Printf("enter %s", msg)
	return func() { log.Printf("exit %s (%s)", msg, time.Since(start)) }
}
