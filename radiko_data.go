package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/exp/utf8string"
)

func getCount(fileTitle string) string {
	fileTitle = strings.Replace(fileTitle, "(最終回)", "", 1)
	fileTitle = strings.TrimSpace(fileTitle)
	title := utf8string.NewString(fileTitle)
	countStr := title.Slice(1, title.RuneCount()-1)
	count, err := strconv.Atoi(countStr)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%02d", count)
}

func getFileInfo(main radikoData, programName string, detail detailList, file fileList) fileInfo {
	var title string
	var fileTitle string

	switch programName {
	case "青春アドベンチャー":
		title = detail.Headline
		fmt.Println(title)
		if title == "" {
			titleRegexp := regexp.MustCompile("「(.*)」(.*)")
			group := titleRegexp.FindSubmatch([]byte(file.FileTitle))
			title = string(group[1])
			fileTitle = title + "_" + getCount(string(group[2]))
			fmt.Println(fileTitle)
		} else {
			fileTitle = title + "_" + getCount(file.FileTitle)
		}

	case "FMシアター":
		ft := utf8string.NewString(file.FileTitle)
		title = ft.Slice(1, ft.RuneCount()-1)
		fileTitle = title

	case "新日曜名作座":
		ft := utf8string.NewString(file.FileTitle)
		title = ft.Slice(1, ft.RuneCount()-1)

		re := regexp.MustCompile(`\((\d*)\)$`)
		countStr := re.FindString(title)

		title = strings.Replace(title, countStr, "", 1)
		fileTitle = title + "_" + getCount(countStr)

	case "特集オーディオドラマ":
		ft := utf8string.NewString(file.FileTitle)
		title = ft.Slice(1, ft.RuneCount()-1)
		fileTitle = title

	default:
		log.Fatal(programName + " is not support")
	}
	return fileInfo{
		title:     title,
		fileTitle: fileTitle,
		fileName:  file.FileName,
	}
}

func getRadikoData(radikoIndex radikoIndex) programInfo {
	jsonURL := radikoIndex.IndexURL
	res, err := http.Get(jsonURL)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err != nil {
			err = res.Body.Close()
		}
	}()
	byteArr, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var jsonData root
	err = json.Unmarshal(byteArr, &jsonData)
	if err != nil {
		log.Fatal(err)
	}

	var fInfoList []fileInfo

	detailList := jsonData.Main.DetailList
	for _, detail := range detailList {
		for _, file := range detail.FileList {
			fInfoList = append(fInfoList, getFileInfo(jsonData.Main, radikoIndex.ProgramName, detail, file))
		}
	}

	pInfo := programInfo{
		programName:  radikoIndex.ProgramName,
		fileInfoList: fInfoList,
	}

	return pInfo
}

type programInfo struct {
	programName  string
	fileInfoList []fileInfo
}

type fileInfo struct {
	title     string
	fileTitle string
	fileName  string
}

type root struct {
	Main radikoData
}

type radikoData struct {
	SiteID      string       `json:"site_id"`
	ProgramName string       `json:"program_name"`
	DetailList  []detailList `json:"detail_list"`
}

type detailList struct {
	Headline string     `json:"headline"`
	FileList []fileList `json:"file_list"`
}

type fileList struct {
	Seq       int    `json:"seq"`
	FileID    string `json:"file_id"`
	FileTitle string `json:"file_title"`
	FileName  string `json:"file_name"`
}
