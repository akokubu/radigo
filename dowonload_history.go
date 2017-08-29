package main

import (
	"os"
	"bufio"
	"log"
)

func isDone(filename, title string) bool {
	// ファイルオープン
	fp, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)

	for scanner.Scan() {
		if scanner.Text() == title {
			return true
		}
	}
	return false
}

func saveDone(filename, title string) {
	// ファイルオープン
	fp, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	writer := bufio.NewWriter(fp)
	writer.WriteString(title + "\n")
	writer.Flush()
}

