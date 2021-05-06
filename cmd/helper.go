package cmd

// small functions that perform miscellaneous tasks


import (
	"github.com/nleeper/goment"
	"log"
	"sort"
	"strconv"
	"strings"
)

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

// SortIntMapByKey return a slice sorted by key and also return the largest key value
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

// RemoveSliceItem - Remove an item from a slice
func RemoveSliceItem(s []goment.Goment, rm int) []goment.Goment {
	var result []goment.Goment
	for i := range s {
		if i == rm {
			continue
		}
		result = append(result,s[i])
	}
	return result
}

// BeheadSlice - remove items from the beginning of a slice
func BeheadSlice(s []goment.Goment, rm int) []goment.Goment {
	var result []goment.Goment
	for i := range s {
		if i <= rm {
			continue
		}
		result = append(result,s[i])
	}
	return result
}

// FindInSlice - return the position of item in a slice, or -1 if not found
func FindInSlice(s []goment.Goment, item goment.Goment) int {
	for i, val := range s {
		if val.Format(dateOutputFmt) == item.Format(dateOutputFmt) {
			return i
		}
	}
	return -1
}
