package fusex

import (
	"log"
	"time"
)

// flog logs the action name and time taken to execute a function.
// It is used to measure the performance of the filesystem
// and provide some insight into the filesystem's behavior.
//
// Usage: defer flog(time.Now(), "action")
func flog(now time.Time, name string) {
	log.Printf("FUSEX: %s: %s", name, time.Since(now))
}
