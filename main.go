package main

import (
    "fmt"
    "log"
	"flag"
    "io"
    "os"
    "sync"
    "path"
    "errors"
    "sort"
    "unicode"
    "bufio"
    "regexp"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "os/exec"
    "strings"
    "golang.org/x/exp/utf8string"
    "github.com/grafov/m3u8"
)

const (
	workDirPath   = "./tmp"
	maxAttempts   = 4
	maxGoroutines = 16
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

const tmpConcatAACFileName = "concat.aac"

type sliceFileInfo []os.FileInfo

func (f sliceFileInfo) Len() int      { return len(f) }
func (f sliceFileInfo) Swap(i, j int) { f[i], f[j] = f[j], f[i] }

func (f sliceFileInfo) Less(i, j int) bool { return naturalComp(f[i].Name(), f[j].Name(), false) < 0 }

/*
func isExist(filename string) bool {
    _, err := os.Stat(filename)
    if err == nil {
        return true
    } else {
        return os.IsExist(err)
    }
}

func createFile(filename string) {
    if !isExist(filename) {
        fmt.Println("file is not exists")
        _, err := os.Create(filename)
        if err != nil {
            log.Fatal(err)
        }
    }
}
*/

func isDone(filename, title string) bool {
    // ファイルオープン
    fp, err := os.Open(filename)
    if err != nil {
        return false
        //log.Fatal(err)
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
    //var jsonUrl string
    var indexPath  string
	flag.StringVar(&indexPath, "i", "index.txt", "json list file")
	//flag.StringVar(&jsonUrl, "j", "https://www.nhk.or.jp/radioondemand/json/0058/bangumi_0058_01.json", "target page url")
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
        //createFile(filename)

        for _, radikoDetail := range radikoData.Detail_list {
            //radikoDetail := radikoData.Detail_list[0]
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

                err := convertM3u8ToMp3(masterM3u8Path, saveDir + "/" + fileName)
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
    /*
    chunks, err := downloadChunks(masterM3u8Path)
    if err != nil {
        return err
    }

    //defer os.RemoveAll(workDirPath)
    if err := os.MkdirAll(workDirPath, 0700); err != nil {
        log.Println(err)
    }

    if err := bulkDownload(maxAttempts, maxGoroutines, chunks, workDirPath); err != nil {
        return err
    }

    */

    /*
    if err := convertTsToMP3(workDirPath, title); err != nil {
        return err
    }
    */

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
        "-acodec","libmp3lame",
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

func convertTsToMP3(aacDirPath, outputFilePath string) error {
    concatFilePath := path.Join(aacDirPath, tmpConcatAACFileName)
    if err := convertConcatAACFile(aacDirPath, concatFilePath); err != nil {
        return err
    }
    return convertAACToMP3(concatFilePath, outputFilePath)
}

func convertAACToMP3(concatFilePath, outputFilePath string) error {
    return nil
}

func concatFileNames(inputDirPath string) (string, error) {
    files, err := ioutil.ReadDir(inputDirPath)
    if err != nil {
        return "", err
    }

    sort.Sort(sliceFileInfo(files))

    var res []byte
    for _, f := range files {
        res = append(res, path.Join(inputDirPath, f.Name())...)
        res = append(res, '|')
    }
    // remove the last element "|"
    return string(res[:len(res)-1]), nil
}

func convertConcatAACFile(inputDirPath, outputFilePath string) error {
    concatFileNames, err := concatFileNames(inputDirPath)
    if err != nil {
        return err
    }

    concatArg := fmt.Sprintf("concat:%s", concatFileNames)
    f, err := newFFMPEG(concatArg)
    if err != nil {
        return err
    }

    f.setArgs("-c", "copy")
    fmt.Println(outputFilePath)
    result, err := f.execute(outputFilePath)
    if err != nil {
        log.Fatal(err)
        //return err
    }

    f.setArgs(
        "-c:a", "libmp3lame",
        "-ac", "2",
        "-q:a", "2",
    )
    result, err = f.execute(outputFilePath)
    log.Println(string(result))
    if err != nil {
        return err
    }
    return nil
}

func downloadChunks(masterM3u8Path string) ([]string, error) {
    resp, err := http.Get(masterM3u8Path)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    chunks, err := readChunks(resp.Body)
    if err != nil {
        return nil, err
    }
    return chunks, nil
}

func bulkDownload(maxAttempts, maxGoroutines int, list []string, output string) error {
    var sem = make(chan struct{}, maxGoroutines)
    var errFlag bool
    var wg sync.WaitGroup

    for _, v := range list {
        wg.Add(1)
        go func(link string) {
            defer wg.Done()

            var err error
            for i := 0; i < maxAttempts; i++ {
                sem <- struct{}{}
                err = download(link, output)
                <-sem
                if err != nil {
                    break
                }
            }
            if err != nil {
                log.Println("Failed to download: %s", err)
                errFlag = true
            }
        }(v)
    }
    wg.Wait()

    if errFlag {
        log.Println("error")
        return errors.New("Lack of asc files")
    }
    return nil
}

func download(link, output string) error {
    resp, err := http.Get(link)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    _, fileName := path.Split(link)
    name := strings.Split(fileName, "?")[0]
    file, err := os.Create(path.Join(output, name))
    fmt.Println(name)
    if err != nil {
        return err
    }

    _, err = io.Copy(file, resp.Body)
    if closeErr := file.Close(); err == nil {
        err = closeErr
    }
    return err
}

func readChunks(input io.Reader) ([]string, error) {
    playlist, listType, err := m3u8.DecodeFrom(input, true)
    if err != nil || listType != m3u8.MEDIA {
        return nil, err
    }
    p := playlist.(*m3u8.MediaPlaylist)

    var chunks []string
    for _, v := range p.Segments {
        if v != nil {
            chunks = append(chunks, v.URI)
        }
    }
    return chunks, nil
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
    Site_id string
    Program_name string
    Detail_list []DetailList
}

type DetailList struct {
    File_list []FileList
}

type FileList struct {
    Seq int
    File_id string
    File_title string
    File_name string
}

// License: https://github.com/mattn/natural#license

func compRight(ra, rb []rune) int {
	bias := 0
	la, lb := len(ra), len(rb)
	var ca, cb rune
	for i := 0; i < la || i < lb; i++ {
		if i < la {
			ca = ra[i]
		} else {
			ca = 0
		}
		if i < lb {
			cb = rb[i]
		} else {
			cb = 0
		}

		da, db := unicode.IsNumber(ca), unicode.IsNumber(cb)
		switch {
		case !da && !db:
			return bias
		case !da:
			return -1
		case !db:
			return 1
		case ca < cb:
			if bias == 0 {
				bias = -1
			}
		case ca > cb:
			if bias == 0 {
				bias = 1
			}
		case ca == 0 && cb == 0:
			return bias
		}
	}

	return 0
}

func compLeft(ra, rb []rune) int {
	la, lb := len(ra), len(rb)
	var ca, cb rune
	i := 0
	for {
		if i < la {
			ca = ra[i]
		} else {
			ca = 0
		}
		if i < lb {
			cb = rb[i]
		} else {
			cb = 0
		}

		da, db := unicode.IsNumber(ca), unicode.IsNumber(cb)
		switch {
		case !da && !db:
			return 0
		case !da:
			return -1
		case !db:
			return 1
		case ca < cb:
			return -1
		case ca > cb:
			return 1
		}
		i++
	}

	return 0
}

func naturalComp(a, b string, foldCase bool) int {
	ra, rb := []rune(a), []rune(b)
	la, lb := len(ra), len(rb)
	ia, ib := 0, 0

	for {
		if ia >= la && ib >= lb {
			return 0
		} else if ia >= la {
			return -1
		} else if ib >= lb {
			return 1
		}
		ca, cb := ra[ia], rb[ib]

		for unicode.IsSpace(ca) {
			ia++
			if ia < la {
				ca = ra[ia]
			} else {
				ca = 0
			}
		}
		for unicode.IsSpace(cb) {
			ib++
			if ib < lb {
				cb = rb[ib]
			} else {
				cb = 0
			}
		}

		if unicode.IsNumber(ca) && unicode.IsNumber(cb) {
			var r int
			if ca == '0' || cb == '0' {
				r = compLeft(ra[ia:], rb[ib:])
				if r != 0 {
					return r
				}
			} else {
				r = compRight(ra[ia:], rb[ib:])
				if r != 0 {
					return r
				}
			}
		}

		if foldCase {
			ca = unicode.ToUpper(ca)
			cb = unicode.ToUpper(cb)
		}

		if ca < cb {
			return -1
		} else if ca > cb {
			return 1
		}

		ia++
		ib++
	}

	return 0
}
