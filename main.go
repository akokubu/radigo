package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"golang.org/x/exp/utf8string"
)

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
		jsonURL := radikoIndex.IndexURL

		radikoData := getRadikoData(jsonURL)
		doneFilename := fmt.Sprintf("%s.txt", radikoIndex.ProgramName)

		makeSaveDir(radikoIndex.ProgramName)

		for _, radikoDetail := range radikoData.DetailList {
			for i, f := range radikoDetail.FileList {
				title := utf8string.NewString(f.FileTitle)
				fileName := title.Slice(1, title.RuneCount()-1)

				re := regexp.MustCompile(`(\(\d\))$`)
				titleName := re.ReplaceAllString(fileName, "")
				saveDir := radikoIndex.ProgramName + "/" + titleName
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

type radikoIndexArray struct {
	Programs []radikoIndex `json:"programs"`
}

type radikoIndex struct {
	ProgramName string `json:"program_name"`
	IndexURL    string `json:"url"`
}

func getRadikoIndexes(indexPath string) []radikoIndex {
	raw, err := ioutil.ReadFile(indexPath)
	if err != nil {
		log.Fatal(err)
	}

	var ri radikoIndexArray
	err = json.Unmarshal(raw, &ri)
	if err != nil {
		log.Fatal(err)
	}

	return ri.Programs
}
