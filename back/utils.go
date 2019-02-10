package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

// will return the last timestamp of the file if it exists
// else ""
// if the file exists and is empty it hangs (should look into this)
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
	return line[:30]

}

func ReadLastLines(fname string, lines_number int) string {
	fileHandle, err := os.OpenFile(fname, os.O_APPEND|os.O_RDONLY, 0644)

	counter := 0
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
			counter += 1
			if counter >= lines_number {
				break
			}
		}

		line = fmt.Sprintf("%s%s", string(char), line)
		if cursor == -filesize { // stop if we are at the begining
			break
		}
	}
	return line

}
