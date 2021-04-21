package shared

import (
	"fmt"
	"github.com/nleeper/goment"
	"log"
	"sort"
	"strconv"
	"strings"
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

func DatesHaveGaps(previous, current *goment.Goment, amount int, period string, debug int) (bool,bool, *goment.Goment) {
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
		fmt.Println("------------------------------------------------------")
		fmt.Println()
	}
	return hasGap, gapOnWeekday, previous
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