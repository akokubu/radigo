package main

import (
	"testing"
)

func TestGetCount(t *testing.T) {
	var tests = []struct {
		s, want string
	}{
		{"「1」", "01"},
		{"「2」", "02"},
		{"「10」", "10"},
	}
	for _, c := range tests {
		got := getCount(c.s)
		if got != c.want {
			t.Errorf("getCount(%q) == %q, want %q", c.s, got, c.want)
		}
	}
}

func TestGetFileInfo_adventure(t *testing.T) {
	main := radikoData{}
	programName := "青春アドベンチャー"
	file := fileList{
		FileName:  "https://nhks-vh.akamaihd.net/i/radioondemand/r/0164/s/stream_0164_262127bf3a6dd84c60385fa0f4e89f3e.mp4/master.m3u8",
		FileTitle: "第1回",
	}
	detail := detailList{
		Headline: "風の向こうへ駆け抜けろ",
	}
	got := getFileInfo(main, programName, detail, file)
	want := fileInfo{
		title:     "風の向こうへ駆け抜けろ",
		fileTitle: "風の向こうへ駆け抜けろ_01",
		fileName:  "https://nhks-vh.akamaihd.net/i/radioondemand/r/0164/s/stream_0164_262127bf3a6dd84c60385fa0f4e89f3e.mp4/master.m3u8",
	}
	if got != want {
		t.Errorf("getFileInfo() == %q, want %q", got, want)
	}
}

func TestGetFileInfo_fmtheater(t *testing.T) {
	main := radikoData{}
	programName := "FMシアター"
	file := fileList{
		FileName:  "https://nhks-vh.akamaihd.net/i/radioondemand/r/0058/s/stream_0058_b6bec48cc2ab5bc5568fe0e5eb5188ed.mp4/master.m3u8",
		FileTitle: "「さよならサリバン先生」",
	}
	detail := detailList{
		Headline: "",
	}
	got := getFileInfo(main, programName, detail, file)
	want := fileInfo{
		title:     "さよならサリバン先生",
		fileTitle: "さよならサリバン先生",
		fileName:  "https://nhks-vh.akamaihd.net/i/radioondemand/r/0058/s/stream_0058_b6bec48cc2ab5bc5568fe0e5eb5188ed.mp4/master.m3u8",
	}
	if got != want {
		t.Errorf("getFileInfo() == %q, want %q", got, want)
	}
}

func TestGetFileInfo_sunday(t *testing.T) {
	main := radikoData{}
	programName := "新日曜名作座"
	file := fileList{
		FileName:  "https://nhks-vh.akamaihd.net/i/radioondemand/r/0930/s/stream_0930_fcb3c25fa336b4b2079d976c28b29b34.mp4/master.m3u8",
		FileTitle: "「多摩川物語(4)」",
	}
	detail := detailList{
		Headline: "",
	}
	got := getFileInfo(main, programName, detail, file)
	want := fileInfo{
		title:     "多摩川物語",
		fileTitle: "多摩川物語_04",
		fileName:  "https://nhks-vh.akamaihd.net/i/radioondemand/r/0930/s/stream_0930_fcb3c25fa336b4b2079d976c28b29b34.mp4/master.m3u8",
	}
	if got != want {
		t.Errorf("getFileInfo() == %q, want %q", got, want)
	}
}

func TestGetFileInfo_special(t *testing.T) {
	main := radikoData{
		ProgramName: "特集オーディオドラマ",
	}
	programName := "特集オーディオドラマ"
	file := fileList{
		FileName:  "https://nhks-vh.akamaihd.net/i/radioondemand/r/P000025/s/stream_P000025_2b1c9dcfca3f50abe499af115c559b3f.mp4/master.m3u://nhks-vh.akamaihd.net/i/radioondemand/r/P000025/s/stream_P000025_2b1c9dcfca3f50abe499af115c559b3f.mp4/master.m3u8",
		FileTitle: "「アシマの銃、セギルの草笛」",
	}
	detail := detailList{
		Headline: "",
	}
	got := getFileInfo(main, programName, detail, file)
	want := fileInfo{
		title:     "アシマの銃、セギルの草笛",
		fileTitle: "アシマの銃、セギルの草笛",
		fileName:  "https://nhks-vh.akamaihd.net/i/radioondemand/r/P000025/s/stream_P000025_2b1c9dcfca3f50abe499af115c559b3f.mp4/master.m3u://nhks-vh.akamaihd.net/i/radioondemand/r/P000025/s/stream_P000025_2b1c9dcfca3f50abe499af115c559b3f.mp4/master.m3u8",
	}
	if got != want {
		t.Errorf("getFileInfo() == %q, want %q", got, want)
	}
}

func TestGetFileInfo_special_old(t *testing.T) {
	main := radikoData{
		ProgramName: "特集オーディオドラマ「ピンザの島」",
	}
	programName := "特集オーディオドラマ"
	file := fileList{
		FileName:  "https://nhks-vh.akamaihd.net/i/radioondemand/r/P000025/s/stream_P000025_2b1c9dcfca3f50abe499af115c559b3f.mp4/master.m3u://nhks-vh.akamaihd.net/i/radioondemand/r/P000025/s/stream_P000025_2b1c9dcfca3f50abe499af115c559b3f.mp4/master.m3u8",
		FileTitle: "2017年8月13日(日)",
	}
	detail := detailList{
		Headline: "",
	}
	got := getFileInfo(main, programName, detail, file)
	want := fileInfo{
		title:     "017年8月13日(日",
		fileTitle: "017年8月13日(日",
		fileName:  "https://nhks-vh.akamaihd.net/i/radioondemand/r/P000025/s/stream_P000025_2b1c9dcfca3f50abe499af115c559b3f.mp4/master.m3u://nhks-vh.akamaihd.net/i/radioondemand/r/P000025/s/stream_P000025_2b1c9dcfca3f50abe499af115c559b3f.mp4/master.m3u8",
	}
	if got != want {
		t.Errorf("getFileInfo() == %q, want %q", got, want)
	}
}
