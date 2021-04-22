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

var workWeek = []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

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
		fmt.Println("     previous :", previous.Format(outputFmt))
		fmt.Println("      current :", current.Format(outputFmt))
		fmt.Println(" gapOnWeekday :", gapOnWeekday)
	}

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
	hasGap := current.IsSameOrAfter(lastAllowable)
	if debug > 998 {
		fmt.Println(" lastAllowable:", lastAllowable.Format(outputFmt))
		fmt.Println("       hasGap :", hasGap)
	}

	if hasGap && !gapOnWeekday && skipWeekends {
		fmt.Println("hasGap, but occurs on weekend and we want to skip weekends")
		epochTime, _ := goment.New()
		return false, epochTime
	}

	if !hasGap {
		epochTime, _ := goment.New(time.Unix(0,0))
		return false, epochTime
	}

	return hasGap, lastAllowable
}


func DatesHaveGaps2(previous, current *goment.Goment, amount int, period string, debug int) (bool,bool, *goment.Goment) {
	if previous.ToTime().IsZero() {
		return false, false, nil
	}

	outputFmt := "L LTS dddd"
	previous = previous.Add(5, "seconds")
	current = current.Subtract(5, "seconds")

	if debug > 998 {
		fmt.Println("     previous :", previous.Format(outputFmt))
		fmt.Println("      current :", current.Format(outputFmt))
	}
	checked, err := goment.New(previous)
	if err != nil {
		log.Fatalf("Unable to instantiate new goment object from 'previous'\n%s\n", err)
	}
	checked.Add(amount, period)
	//checked := previous.Add(amount, period)
	hasGap := !checked.IsSameOrAfter(current)
	if debug > 998 {
		fmt.Printf("IsSameOrAfter : %s : %v\n", checked.Format(outputFmt), hasGap)
	}
	gapOnWeekday := stringInSlice(checked.Format("dddd"), workWeek)
	if !gapOnWeekday {
		if debug > 998 {
			fmt.Println("previous      :", previous.Format("dddd"))
			fmt.Println("current       :", current.Format("dddd"))
		}

		if previous.Format("dddd") == "Friday" && current.Format("dddd") == "Tuesday" {
			if debug > 998 {
				fmt.Println("missed Monday : true")
			}
			previous.Add(72, "hours")
			gapOnWeekday = true
		}
	}
	if hasGap && debug > 998 {
		fmt.Printf(" gapOnWeekday : %v - %s\n", gapOnWeekday, checked.Format("dddd"))
		diff := current.ToTime().Sub(previous.ToTime())
		fmt.Printf("diff in hours : %s\n", diff)
	}
	if debug > 998 {
		fmt.Println("------------------------------------------------------ ",hasGap, gapOnWeekday, checked.Format(outputFmt))
		fmt.Println()
	}
	return hasGap, gapOnWeekday, checked
}

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