package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func getRadikoData(jsonURL string) radikoData {
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
	return jsonData.Main
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
	FileList []fileList `json:"file_list"`
}

type fileList struct {
	Seq       int    `json:"seq"`
	FileID    string `json:"file_id"`
	FileTitle string `json:"file_title"`
	FileName  string `json:"file_name"`
}
