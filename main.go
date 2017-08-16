package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"

	"github.com/grafov/m3u8"
	"golang.org/x/exp/utf8string"
)

type ffmpeg struct {
	*exec.Cmd
}

func newFFMPEG(inputFilePath string) (*ffmpeg, error) {
	cmdPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return nil, err
	}

	return &ffmpeg{exec.Command(cmdPath, "-i", inputFilePath)}, nil
}

func (f *ffmpeg) setArgs(args ...string) {
	f.Args = append(f.Args, args...)
}

func (f *ffmpeg) execute(output string) ([]byte, error) {
	fmt.Println("ffmpeg")
	f.Args = append(f.Args, output)
	fmt.Println(f.Args)
	return f.CombinedOutput()
}

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
		jsonUrl := scanner.Text()
		fmt.Println(jsonUrl)

		radikoData := getRadikoData(jsonUrl)
		doneFilename := fmt.Sprintf("%s.txt", radikoData.Program_name)
		_, err := os.Stat(radikoData.Program_name)
		if err != nil {
			if err := os.Mkdir(radikoData.Program_name, 0777); err != nil {
				log.Fatal(err)
			}
		}

		for _, radikoDetail := range radikoData.Detail_list {
			for i, f := range radikoDetail.File_list {
				title := utf8string.NewString(f.File_title)
				fileName := title.Slice(1, title.RuneCount()-1)

				re := regexp.MustCompile(`(\(\d\))$`)
				titleName := re.ReplaceAllString(fileName, "")
				saveDir := radikoData.Program_name + "/" + titleName
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
				m3u8FilePath := f.File_name
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

func convertM3u8ToMp3(masterM3u8Path, title string) error {
	f, err := newFFMPEG(masterM3u8Path)
	if err != nil {
		return err
	}

	f.setArgs(
		"-protocol_whitelist", "file,crypto,http,https,tcp,tls",
		"-movflags", "faststart",
		"-c", "copy",
		"-y",
		"-bsf:a", "aac_adtstoasc",
	)

	result, err := f.execute("output.mp4")
	log.Println(string(result))
	if err != nil {
		return err
	}

	f, err = newFFMPEG("output.mp4")
	if err != nil {
		return err
	}

	f.setArgs(
		"-y",
		"-acodec", "libmp3lame",
		"-ab", "256k",
	)

	var name = title + ".mp3"
	fmt.Println(name)

	result, err = f.execute(name)
	log.Println(string(result))
	if err != nil {
		return err
	}
	return nil
}

func getM3u8MasterPlaylist(m3u8FilePath string) string {
	resp, err := http.Get(m3u8FilePath)
	if err != nil {
		log.Fatal(err)
	}
	f := resp.Body

	p, t, err := m3u8.DecodeFrom(f, true)
	if err != nil {
		log.Fatal(err)
	}

	if t != m3u8.MASTER {
		log.Fatal("not support file type [%d]", t)
	}

	return p.(*m3u8.MasterPlaylist).Variants[0].URI
}

func getRadikoData(jsonUrl string) RadikoData {
	res, _ := http.Get(jsonUrl)
	defer res.Body.Close()
	byteArr, _ := ioutil.ReadAll(res.Body)

	var jsonData MainData
	err := json.Unmarshal(byteArr, &jsonData)
	if err != nil {
		log.Fatal(err)
	}
	return jsonData.Main
}

type MainData struct {
	Main RadikoData
}

type RadikoData struct {
	Site_id      string
	Program_name string
	Detail_list  []DetailList
}

type DetailList struct {
	File_list []FileList
}

type FileList struct {
	Seq        int
	File_id    string
	File_title string
	File_name  string
}
