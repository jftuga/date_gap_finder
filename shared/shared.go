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

func DatesHaveGaps(previous, current *goment.Goment, amount int, period string) (bool,bool) {
	debug := true
	if previous.ToTime().IsZero() {
		return false, false
	}

	outputFmt := "L LTS dddd"
	previous = previous.Add(5, "seconds")
	current = current.Subtract(5, "seconds")

	if debug {
		fmt.Println("     previous :", previous.Format(outputFmt))
		fmt.Println("      current :", current.Format(outputFmt))
	}
	checked := previous.Add(amount, period)
	hasGap := !checked.IsSameOrAfter(current)
	if debug {
		fmt.Printf("IsSameOrAfter : %s : %v\n", checked.Format(outputFmt), hasGap)
	}
	gapOnWeekday := stringInSlice(checked.Format("dddd"), workWeek)
	if hasGap && debug {
		fmt.Printf(" gapOnWeekday : %v - %s\n", gapOnWeekday, checked.Format("dddd"))
	}
	if debug {
		fmt.Println("------------------------------------------------------")
		fmt.Println()
	}
	return hasGap, gapOnWeekday
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