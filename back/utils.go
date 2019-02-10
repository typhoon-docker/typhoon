package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

// Will return the last timestamp of the file if it exists, else ""
// If the file exists and is empty it hangs (TODO should look into this)
func readLastLineTimestamp(fname string) string {
	fileHandle, err := os.OpenFile(fname, os.O_APPEND|os.O_RDONLY, 0644)

	if os.IsNotExist(err) {
		return ""
	} else if err != nil {
		log.Fatal("Cannot open file")
	}
	defer fileHandle.Close()

	line := ""
	var cursor int64 = 0
	stat, _ := fileHandle.Stat()
	filesize := stat.Size()
	if filesize == 0 {
		return ""
	}
	for {
		cursor -= 1
		fileHandle.Seek(cursor, io.SeekEnd)

		char := make([]byte, 1)
		fileHandle.Read(char)

		if cursor != -1 && (char[0] == 10 || char[0] == 13) { // stop if we find a line
			break
		}

		line = fmt.Sprintf("%s%s", string(char), line)

		if cursor == -filesize { // stop if we are at the begining
			break
		}
	}

	return increaseTimestamp(line[:30])
}

// Add 1 millisecond to the given timestamp
func increaseTimestamp(timestamp string) string {
	layout := "2006-01-02T15:04:05.000Z"
	t, err := time.Parse(layout, timestamp)

	if err != nil {
		fmt.Println(err)
	}

	t.Add(time.Millisecond * 1)
	return t.Format(layout)
}

// Read the last lines of the file
func ReadLastLines(fname string, lines_number int) string {
	fileHandle, err := os.OpenFile(fname, os.O_APPEND|os.O_RDONLY, 0644)

	counter := 0
	if os.IsNotExist(err) {
		return ""
	} else if err != nil {
		log.Fatal("Cannot open file")
	}
	defer fileHandle.Close()

	var sb strings.Builder
	var cursor int64 = 0
	stat, _ := fileHandle.Stat()
	filesize := stat.Size()
	if filesize == 0 {
		return ""
	}
	for {
		cursor -= 1
		fileHandle.Seek(cursor, io.SeekEnd)

		char := make([]byte, 1)
		fileHandle.Read(char)

		if cursor != -1 && (char[0] == 10 || char[0] == 13) { // stop if we find a line
			counter += 1
			if counter >= lines_number {
				break
			}
		}

		sb.WriteString(string(char))
		if cursor == -filesize { // stop if we are at the begining
			break
		}
	}
	return Reverse(sb.String())
}

// Reverse the string
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
