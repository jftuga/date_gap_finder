package fileOps

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"github.com/jftuga/date_gap_finder/filecopy"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func CsvOpenRead(fname string) (*csv.Reader, *os.File) {
	file, err := os.Open(fname)
	if err != nil {
		log.Fatalf("Unable to open file: '%s'\n%s\n", fname, err)
	}
	//defer file.Close()
	return csv.NewReader(file), file
}

func CsvOpenWriteFile(fname string) (*csv.Writer, *os.File) {
	file, err := os.Create(fname)
	if err != nil {
		log.Fatalf("Unable to open file: '%s'\n%s\n", fname, err)
	}
	//defer file.Close()
	return csv.NewWriter(file), file
}

func CsvOpenWriteBuf() *csv.Writer {
	var buf bytes.Buffer
	return csv.NewWriter(&buf)
}

func SaveToCsv(fname string, data [][]string) {
	fname = "new--" + fname
	w := csv.NewWriter(os.Stdout)
	w.WriteAll(data)
}

// SplitFilename return the filename without extension and the extension (with leading dot)
func SplitFilename(fileName string) (string, string) {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName)), filepath.Ext(fileName)
}

// RenameFile rename a file; expects full path name for both args
func RenameFile(oldPath, newPath string) bool {
	err := os.Rename(oldPath, newPath)
	if err != nil {
		log.Fatalf("Error #20050: Unable to rename file from '%s' to '%s'; %s\n", oldPath, newPath, err)
		return false
	}
	return true
}

// OverwriteCsv a CSV file with new data; also create a backup file with date and .bak extension
func OverwriteCsv(fname string, data []string) bool {
	base, _ :=  SplitFilename(fname)
	t := time.Now()
	now := t.Format("20060102.150405")
	newFilename := base + "--" + now + ".bak"

	err := filecopy.Copy(newFilename, fname)
	if err != nil {
		log.Fatalf("Error #20053: Unable to copy file from '%s' to '%s'; %s\n", fname, newFilename, err)
		return false
	}

	file, err := os.Create(fname)
	if err != nil {
		log.Fatalf("Error #20055: Unable to open file for writing: '%s'; %s\n", fname, err)
		return false
	}
	w := bufio.NewWriter(file)
	for _, row := range data {
		_, err = w.WriteString(row + "\n")
		if err != nil {
			log.Fatalf("Error #20060: Unable to write CSV data to file: '%s'; %s\n", fname, err)
			return false
		}
	}
	w.Flush()
	err = file.Close()
	if err != nil {
		log.Fatalf("Error #20065: Unable to close CSV file: '%s'; %s\n", fname, err)
		return false
	}
	return true
}
