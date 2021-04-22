package shared

import (
	"github.com/nleeper/goment"
	"testing"
)

// TestDatesHaveGaps0 - there are no gaps
func TestDatesHaveGaps0(t *testing.T) {
	debug := 999
	skipWeekends := true
	amount := 25
	period := "hours"
	previous, _ := goment.New("2021-03-31 18:40:01")
	current, _ := goment.New("2021-04-01 18:40:09")

	hasGap, dateSkipped := DatesHaveGap(previous, current, amount, period, skipWeekends, debug)
	t.Logf("      hasGap : %v\n", hasGap )
	outputFmt := "L LTS dddd"
	t.Logf("dateSkipped : %s\n", dateSkipped.Format(outputFmt) )
	if hasGap == true {
		t.Error("'hasGap' should have returned true")
	}

	correctAnswer := "12/31/1969 7:00:00 PM Wednesday"
	if dateSkipped.Format(outputFmt) != correctAnswer {
		t.Errorf("'dateSkipped' returned '%s' but should have returned '%s'\n", dateSkipped.Format(outputFmt), correctAnswer)
	}
}

// TestDatesHaveGaps1 - test for a date gap occurring on a Weekday - Friday
func TestDatesHaveGaps1(t *testing.T) {
	debug := 999
	skipWeekends := true
	amount := 25
	period := "hours"
	previous, _ := goment.New("2021-03-11 18:40:01")
	current, _ := goment.New("2021-03-15 18:40:09")

	hasGap, dateSkipped := DatesHaveGap(previous, current, amount, period, skipWeekends, debug)
	t.Logf("      hasGap : %v\n", hasGap )
	outputFmt := "L LTS dddd"
	t.Logf("dateSkipped : %s\n", dateSkipped.Format(outputFmt) )
	if hasGap == false {
		t.Error("'hasGap' should have returned true")
	}
	correctAnswer := "03/12/2021 7:40:06 PM Friday"
	if dateSkipped.Format(outputFmt) != correctAnswer {
		t.Errorf("'dateSkipped' returned '%s' but should have returned '%s'\n", dateSkipped.Format(outputFmt), correctAnswer)
	}
}

// TestDatesHaveGaps2 - test for a date gap occurring on a Weekday - Monday
func TestDatesHaveGaps2(t *testing.T) {
	debug := 999
	skipWeekends := true
	amount := 25
	period := "hours"
	previous, _ := goment.New("2021-03-12 18:40:01")
	current, _ := goment.New("2021-03-16 18:40:08")

	hasGap, dateSkipped := DatesHaveGap(previous, current, amount, period, skipWeekends, debug)
	t.Logf("      hasGap : %v\n", hasGap )
	outputFmt := "L LTS dddd"
	t.Logf("dateSkipped : %s\n", dateSkipped.Format(outputFmt) )
	if hasGap == false {
		t.Error("'hasGap' should have returned true")
	}
	correctAnswer := "03/15/2021 7:40:06 PM Monday"
	if dateSkipped.Format(outputFmt) != correctAnswer {
		t.Errorf("'dateSkipped' returned '%s' but should have returned '%s'\n", dateSkipped.Format(outputFmt), correctAnswer)
	}
}

