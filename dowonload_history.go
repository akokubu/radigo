package main

import (
	"bufio"
	"log"
	"os"
)

func isDone(filename, title string) bool {
	// ファイルオープン
	fp, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer func() {
		if err != nil {
			err = fp.Close()
		}
	}()

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
	defer func() {
		if err != nil {
			err = fp.Close()
		}
	}()

	writer := bufio.NewWriter(fp)
	_, err = writer.WriteString(title + "\n")
	if err != nil {
		log.Fatal(err)
	}
	err = writer.Flush()
	if err != nil {
		log.Fatal(err)
	}
}
