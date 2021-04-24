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

// stringInSlice - determine is a string is contained within a slice
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// DatesHaveGap - determine if the is a gap between two dates
// the length of the gap is defined by 'amount' and 'period'
// examples: a:25 p:hours; a:1 p:days; a:5 p:minutes
// there is a 10 second 'grace period' built into this function in case date values are off by just a few seconds
// you can also skip weekends, thus, a 'previous' date of Friday and a 'current' date of Monday would not have a gap
// returns: a bool for having a date gap or not
// if gap, then also return the last allowable date; otherwise the 'epoch time' of "12/31/1969 7:00:00 PM Wednesday"
func DatesHaveGap(previous, current *goment.Goment, amount int, period string, skipWeekends bool, debug int) (bool,*goment.Goment) {
	if previous.ToTime().IsZero() {
		epochTime, _ := goment.New()
		return false, epochTime
	}

	outputFmt := "L LTS dddd"
	previous = previous.Add(5, "seconds")
	current = current.Subtract(5, "seconds")
	gapOnWeekday := stringInSlice(current.Format("dddd"), workWeek)
	if debug > 998 {
		fmt.Println("vvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvv")
		fmt.Println("     previous :", previous.Format(outputFmt))
		fmt.Println("      current :", current.Format(outputFmt))
		fmt.Println(" skipWeekends :", skipWeekends)
		fmt.Println(" gapOnWeekday :", gapOnWeekday)
	}

	// any date after this would be consider a gap
	lastAllowable, err := goment.New(previous)
	if err != nil {
		log.Fatalf("Unable to clone 'lastAllowable' goment object from 'previous'\n%s\n", err)
	}

	// compare previous against what the last allowable date is, before it is considered a date gap
	lastAllowable.Add(amount, period)
	if skipWeekends && lastAllowable.Format("dddd") == "Saturday" {
		lastAllowable.Add(48, "hours")
	}
	if skipWeekends && lastAllowable.Format("dddd") == "Sunday" {
		lastAllowable.Add(24, "hours")
	}

	// compare the last date allowed before a gap occurs against the current date
	hasGap := current.IsSameOrAfter(lastAllowable)
	if debug > 998 {
		fmt.Println(" lastAllowable:", lastAllowable.Format(outputFmt))
		fmt.Println("       hasGap :", hasGap)
	}

	// FIXME: remove after development and testing is initially completed
	if hasGap && !gapOnWeekday && skipWeekends {
		if debug > 998 {
			fmt.Println("hasGap, but occurs on weekend and we want to skip weekends")
		}
		log.Fatalf("Should not reach this code. Curious as to what circumstances this occurs")
		epochTime, _ := goment.New()
		return false, epochTime
	}

	if !hasGap {
		epochTime, _ := goment.New(time.Unix(0,0))
		return false, epochTime
	}

	return hasGap, lastAllowable
}

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

// GetDurationInSeconds - convert from allRootOptions.Amount, allRootOptions.Period
// to seconds
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


func GetTimeDifference(a, b goment.Goment) time.Duration {
	aTime := a.ToTime()
	bTime := b.ToTime()
	d := bTime.Sub(aTime)
	return d
}

// removeIndex - remove an item from a slice
// https://stackoverflow.com/a/67060285/452281
func RemoveIndex(items []goment.Goment, idx int) []goment.Goment{
	ret := make([]goment.Goment, len(items)-1)
	copy(ret[:idx], items[:idx])
	copy(ret[idx:], items[idx+1:])
	return ret
}