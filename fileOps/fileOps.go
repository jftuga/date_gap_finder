package fileOps

import (
	"bytes"
	"encoding/csv"
	"log"
	"os"
)

func CsvOpenRead(fname string) *csv.Reader {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatalf("Unable to open file: '%s'\n%s\n", fname, err)
	}
	//defer file.Close()
	return csv.NewReader(file)
}

// FIXME: how should this be closed, since 'defer' can't be used here
func CsvOpenWriteFile(fname string) *csv.Writer {
	file, err := os.Create(fname)
	if err != nil {
		log.Fatalf("Unable to open file: '%s'\n%s\n", fname, err)
	}
	//defer file.Close()
	return csv.NewWriter(file)
}

// FIXME: how should writer.Flush be called
func CsvOpenWriteBuf() *csv.Writer {
	var buf bytes.Buffer
	return csv.NewWriter(&buf)
}

func SaveToCsv(fname string, data [][]string) {
	fname = "new--" + fname
	w := csv.NewWriter(os.Stdout)
	w.WriteAll(data)
}