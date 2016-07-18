package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
)

// MeasureTime can be used to measue time from start of a function until it ends.
// Use like:
//
// func factorial(n *big.Int) (result *big.Int) {
//     defer MeasureTime(time.Now(), fields, "factorial")
//     // ... do some things, maybe even return under some condition
//     return n
// }
//
func MeasureTime(start time.Time, fields log.Fields, message string) {
	elapsed := time.Since(start)

	log.WithFields(fields).WithField("time_used", elapsed.Nanoseconds()).Info(message)
}
