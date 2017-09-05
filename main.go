package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"

	"golang.org/x/exp/utf8string"
)

type radikoIndex struct {
	jsonURL string
}

func makeSaveDir(programName string) {
	_, err := os.Stat(programName)
	if err != nil {
		if err := os.Mkdir(programName, 0777); err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	var indexPath string
	flag.StringVar(&indexPath, "i", "index.txt", "json list file")
	flag.Parse()

	radikoIndexes := getRadikoIndexes(indexPath)
	for _, radikoIndex := range radikoIndexes {
		jsonURL := radikoIndex.jsonURL
		fmt.Println(jsonURL)

		radikoData := getRadikoData(jsonURL)
		doneFilename := fmt.Sprintf("%s.txt", radikoData.ProgramName)

		makeSaveDir(radikoData.ProgramName)

		for _, radikoDetail := range radikoData.DetailList {
			for i, f := range radikoDetail.FileList {
				title := utf8string.NewString(f.FileTitle)
				fileName := title.Slice(1, title.RuneCount()-1)

				re := regexp.MustCompile(`(\(\d\))$`)
				titleName := re.ReplaceAllString(fileName, "")
				saveDir := radikoData.ProgramName + "/" + titleName
				if i == 0 {
					_, err := os.Stat(saveDir)
					if err != nil {
						if err := os.Mkdir(saveDir, 0777); err != nil {
							log.Fatal(err)
						}
					}
				}

				fmt.Print(fileName + " ")
				if isDone(doneFilename, fileName) {
					fmt.Println("already downloaded")
					continue
				}
				m3u8FilePath := f.FileName
				masterM3u8Path := getM3u8MasterPlaylist(m3u8FilePath)

				err := convertM3u8ToMp3(masterM3u8Path, saveDir+"/"+fileName)
				if err != nil {
					log.Fatal(err)
				}
				saveDone(doneFilename, fileName)
				fmt.Println("done")
			}
		}
	}
}

func getRadikoIndexes(indexPath string) []radikoIndex {
	fp, err := os.Open(indexPath)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err != nil {
			err = fp.Close()
		}
	}()

	scanner := bufio.NewScanner(fp)
	var indexes []radikoIndex
	for scanner.Scan() {
		jsonURL := scanner.Text()
		indexes = append(indexes, radikoIndex{jsonURL: jsonURL})
	}
	return indexes
}
