package shared

import (
	"fmt"
	"github.com/nleeper/goment"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

// used to skip weekends, do not change this
var workWeek = []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}

// GetKeyVal - split a string of "N,S" into a number and a string
// 1,-1 => 1 (int); "-1" (string)
func GetKeyVal(combined string) (int,string) {
	slots := strings.Split(combined,",")
	if len(slots) != 2 {
		log.Fatalf("Error: can not split in to column number and column value: '%s'", combined)
	}

	i, err := strconv.Atoi(slots[0])
	if err != nil {
		log.Fatalf("Error: can't convert to integer: '%s'", slots[0])
	}
	return i, strings.TrimSpace(slots[1])
}

// SortIntMapByKey return a map sorted by key and also return the largest key value
func SortIntMapByKey(m map[int]string) ([]int,int) {
	keys := make([]int, len(m))
	i := 0
	largest := 0
	for k := range m {
		if k > largest {
			largest = k
		}
		keys[i] = k
		i++
	}
	sort.Ints(keys)
	return keys, largest
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
		a += 24
	} else {
		log.Fatalf("Error #80620: unable to convert to time.Duration: '%d, %s'\n", a, period)
	}

	parsed, err := time.ParseDuration(fmt.Sprintf("%d%s", a, duration))
	if err != nil {
		log.Fatalf("Error #80625: unable to convert to time.Duration: '%d, %s'; %s\n", a, period, err)
	}
	return int(parsed.Seconds())
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
		a += 24
	} else {
		log.Fatalf("Error #80620: unable to convert to time.Duration: '%d, %s'\n", a, period)
	}

	parsed, err := time.ParseDuration(fmt.Sprintf("%d%s", a, duration))
	if err != nil {
		log.Fatalf("Error #80625: unable to convert to time.Duration: '%d, %s'; %s\n", a, period, err)
	}
	return parsed
}

// GetTimeDifference - return the difference between to Goment objects
func GetTimeDifference(a, b goment.Goment) time.Duration {
	aTime := a.ToTime()
	bTime := b.ToTime()
	d := bTime.Sub(aTime)
	return d
}
