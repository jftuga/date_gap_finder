package cmd

// date-time library
// a group of functions related to time

import (
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/nleeper/goment"
	"log"
	"time"
)

// ConvertDate - convert a date to a different layout
func ConvertDate(t time.Time, layoutAny string) string {
	layout, err := dateparse.ParseFormat(layoutAny)
	if err != nil {
		log.Fatalf("Error #80050: Can parse '%s'; %s\n", layoutAny, err)
	}
	return t.Format(layout)
}

// GetDuration - convert from allRootOptions.Amount, allRootOptions.Period
// to a time.Duration
func GetDuration(a int, period string) time.Duration {
	duration := ""
	if period == "hours" {
		duration = "h"
	} else if period == "minutes" {
		duration = "m"
	} else if period == "seconds" {
		duration = "s"
	} else if period == "days" {
		duration="h"
		a *= 24
	} else {
		log.Fatalf("Error #80620: unable to convert to time.Duration: '%d, %s'\n", a, period)
	}

	parsed, err := time.ParseDuration(fmt.Sprintf("%d%s", a, duration))
	if err != nil {
		log.Fatalf("Error #80625: unable to convert to time.Duration: '%d, %s'; %s\n", a, period, err)
	}
	return parsed
}

// GetDurationInSeconds - convert from allRootOptions.Amount, allRootOptions.Period to seconds
func GetDurationInSeconds(a int, period string) int {
	duration := ""
	if period == "hours" {
		duration = "h"
	} else if period == "minutes" {
		duration = "m"
	} else if period == "seconds" {
		duration = "s"
	} else if period == "days" {
		duration="h"
		a *= 24
	} else {
		log.Fatalf("Error #80620: unable to convert to time.Duration: '%d, %s'\n", a, period)
	}

	parsed, err := time.ParseDuration(fmt.Sprintf("%d%s", a, duration))
	if err != nil {
		log.Fatalf("Error #80625: unable to convert to time.Duration: '%d, %s'; %s\n", a, period, err)
	}
	return int(parsed.Seconds())
}

// GetTimeDifference - return the difference between to Goment objects
func GetTimeDifference(a, b goment.Goment) time.Duration {
	aTime := a.ToTime()
	bTime := b.ToTime()
	d := bTime.Sub(aTime)
	return d
}
