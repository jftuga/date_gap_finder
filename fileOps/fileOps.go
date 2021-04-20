package fileOps

import (
	"encoding/csv"
	"log"
	"os"
)

func CsvOpen(fname string) *csv.Reader {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatalf("Unable to open file: '%s'\n%s\n", fname, err)
	}
	//defer file.Close()
	return csv.NewReader(file)
}
