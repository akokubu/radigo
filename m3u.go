package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/grafov/m3u8"
)

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
	return err
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
		log.Fatalf("not support file type [%v]", t)
	}

	return p.(*m3u8.MasterPlaylist).Variants[0].URI
}
