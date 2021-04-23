package shared

import (
	"github.com/matryer/is"
	"github.com/nleeper/goment"
	"os"
	"testing"
)

const defaultConfigNum int = 1
const epochDate string = "12/31/1969 7:00:00 PM Wednesday"

func config(config int) (int, bool, int, string, string) {
	var debug, amount int
	var skipWeekends bool
	var period string
	outputFmt := "L LTS dddd"

	// defaults
	debug = 994 + len(os.Args)
	skipWeekends = true
	amount = 25
	period = "hours"

	if config == 2 {
		skipWeekends = false
	}

	if config == 3 {
		amount = 1
		period = "days"
	}

	if config == 4 {
		amount = 5
		period = "minutes"
	}

	return debug, skipWeekends, amount, period, outputFmt
}

// TestNoDateGaps1 - there are no date gaps, 25 hour period
func TestNoDateGaps1(t *testing.T) {
	debug, skipWeekends, amount, period, outputFmt := config(defaultConfigNum)
	previous, _ := goment.New("2021-03-31 18:40:01") // Wed
	current, _ := goment.New("2021-04-01 18:40:09") // Thur

	hasGap, dateSkipped := DatesHaveGap(previous, current, amount, period, skipWeekends, debug)
	iss := is.New(t)
	iss.True(!hasGap)
	iss.Equal(dateSkipped.Format(outputFmt), epochDate )
}

// TestNoDateGaps2 - there are no date gaps, 5 minute period (test the 10-second granularity)
func TestNoDateGaps2(t *testing.T) {
	debug, skipWeekends, amount, period, outputFmt := config(4)
	previous, _ := goment.New("2021-03-31 18:40:30") // Wed
	current, _ := goment.New("2021-03-31 18:45:39") // Wed

	hasGap, dateSkipped := DatesHaveGap(previous, current, amount, period, skipWeekends, debug)
	iss := is.New(t)
	iss.True(!hasGap)
	iss.Equal(dateSkipped.Format(outputFmt), epochDate )
}

// TestGapOnFriday1 - test for a date gap occurring on a Friday, 25 hour period
func TestGapOnFriday1(t *testing.T) {
	debug, skipWeekends, amount, period, outputFmt := config(defaultConfigNum)
	previous, _ := goment.New("2021-03-11 18:40:01") // Thur
	current, _ := goment.New("2021-03-15 18:40:09") // Mon

	hasGap, dateSkipped := DatesHaveGap(previous, current, amount, period, skipWeekends, debug)
	iss := is.New(t)
	iss.True(hasGap)
	iss.Equal(dateSkipped.Format(outputFmt), "03/12/2021 7:40:06 PM Friday")
}

// TestGapOnFriday2 - test for a date gap occurring on a Friday, one day period
func TestGapOnFriday2(t *testing.T) {
	debug, skipWeekends, amount, period, outputFmt := config(3)
	previous, _ := goment.New("2021-03-11 18:40:21") // Thur
	current, _ := goment.New("2021-03-15 18:40:09") // Mon

	hasGap, dateSkipped := DatesHaveGap(previous, current, amount, period, skipWeekends, debug)
	iss := is.New(t)
	iss.True(hasGap)
	iss.Equal(dateSkipped.Format(outputFmt), "03/12/2021 6:40:26 PM Friday")
}

// TestGapOnMonday - test for a date gap occurring on a Monday
func TestGapOnMonday(t *testing.T) {
	debug, skipWeekends, amount, period, outputFmt := config(defaultConfigNum)
	previous, _ := goment.New("2021-03-12 18:40:01") // Fri
	current, _ := goment.New("2021-03-16 18:40:08") // Tue

	hasGap, dateSkipped := DatesHaveGap(previous, current, amount, period, skipWeekends, debug)
	iss := is.New(t)
	iss.True(hasGap)
	iss.Equal(dateSkipped.Format(outputFmt), "03/15/2021 7:40:06 PM Monday")
}

// TestGapOnSaturday1 - test for a date gap occurring on a Saturday, skipWeekends=true (therefore, no gaps)
func TestGapOnSaturday1(t *testing.T) {
	debug, skipWeekends, amount, period, outputFmt := config(defaultConfigNum)
	previous, _ := goment.New("2021-03-12 18:40:01") // Fri
	current, _ := goment.New("2021-03-14 18:40:08")  // Sun

	hasGap, dateSkipped := DatesHaveGap(previous, current, amount, period, skipWeekends, debug)
	iss := is.New(t)
	iss.True(!hasGap)
	iss.Equal(dateSkipped.Format(outputFmt), epochDate)
}

// TestGapOnSaturday2 - test for a date gap occurring on a Saturday, skipWeekends=false (therefore, gap occurs)
func TestGapOnSaturday2(t *testing.T) {
	debug, skipWeekends, amount, period, outputFmt := config(2)
	previous, _ := goment.New("2021-03-12 18:40:01") // Fri
	current, _ := goment.New("2021-03-14 18:40:08")  // Sun

	hasGap, dateSkipped := DatesHaveGap(previous, current, amount, period, skipWeekends, debug)
	iss := is.New(t)
	iss.True(hasGap)
	iss.Equal(dateSkipped.Format(outputFmt), "03/13/2021 7:40:06 PM Saturday")
}

// TestGapOnSunday1 - test for a date gap occurring on a Sunday, skipWeekends=true (therefore, no gaps)
func TestGapOnSunday1(t *testing.T) {
	debug, skipWeekends, amount, period, outputFmt := config(defaultConfigNum)
	previous, _ := goment.New("2021-04-03 11:25:26") // Sat
	current, _ := goment.New("2021-04-05 11:25:00")  // Mon

	hasGap, dateSkipped := DatesHaveGap(previous, current, amount, period, skipWeekends, debug)
	iss := is.New(t)
	iss.True(!hasGap)
	iss.Equal(dateSkipped.Format(outputFmt), epochDate)
}

// TestGapOnSunday2 - test for a date gap occurring on a Sunday, skipWeekends=false (therefore, gap occurs)
func TestGapOnSunday2(t *testing.T) {
	debug, skipWeekends, amount, period, outputFmt := config(2)
	previous, _ := goment.New("2021-04-03 11:25:26") // Sat
	current, _ := goment.New("2021-04-05 11:25:00")  // Mon

	hasGap, dateSkipped := DatesHaveGap(previous, current, amount, period, skipWeekends, debug)
	iss := is.New(t)
	iss.True(hasGap)
	iss.Equal(dateSkipped.Format(outputFmt), "04/04/2021 12:25:31 PM Sunday")
}
