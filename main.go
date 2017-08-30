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

func main() {
	var indexPath string
	flag.StringVar(&indexPath, "i", "index.txt", "json list file")
	flag.Parse()

	fp, err := os.Open(indexPath)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)

	for scanner.Scan() {
		jsonURL := scanner.Text()
		fmt.Println(jsonURL)

		radikoData := getRadikoData(jsonURL)
		doneFilename := fmt.Sprintf("%s.txt", radikoData.ProgramName)
		_, err := os.Stat(radikoData.ProgramName)
		if err != nil {
			if err := os.Mkdir(radikoData.ProgramName, 0777); err != nil {
				log.Fatal(err)
			}
		}

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
