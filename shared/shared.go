package shared

import (
	"fmt"
	"github.com/nleeper/goment"
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
	debug := false
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
