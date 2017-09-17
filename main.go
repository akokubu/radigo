package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
		radikoData := getRadikoData(radikoIndex)
		doneFilename := fmt.Sprintf("%s.txt", radikoIndex.ProgramName)

		makeSaveDir(radikoIndex.ProgramName)
		fmt.Println(radikoIndex.ProgramName)

		for i, fileInfo := range radikoData.fileInfoList {
			title := fileInfo.title

			fmt.Print(fileInfo.title)
			saveDir := radikoIndex.ProgramName + "/" + radikoData.programName

			if i == 0 {
				_, err := os.Stat(saveDir)
				if err != nil {
					if err := os.Mkdir(saveDir, 0777); err != nil {
						log.Fatal(err)
					}
				}
			}

			if isDone(doneFilename, title) {
				fmt.Println(" already downloaded")
				continue
			}

			// MP3保存
			m3u8FilePath := fileInfo.fileName
			masterM3u8Path := getM3u8MasterPlaylist(m3u8FilePath)
			err := convertM3u8ToMp3(masterM3u8Path, saveDir+"/"+fileInfo.title)
			if err != nil {
				log.Fatal(err)
			}

			saveDone(doneFilename, title)
			fmt.Println("done")
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
