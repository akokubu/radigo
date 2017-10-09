radigo
====

Overview

らじる聴き逃し番組MP3ダウンロード

Go勉強用に作ってみた。

# Usage
findex.txt ファイル作成

```
{
  "programs": [
    {
      "program_name": "FMシアター",
      "url": "https://www.nhk.or.jp/radioondemand/json/0058/bangumi_0058_01.json"
    },
    {
      "program_name": "新日曜名作座",
      "url": "https://www.nhk.or.jp/radioondemand/json/0930/bangumi_0930_01.json"
    },
    {
      "program_name": "青春アドベンチャー",
      "url": "https://www.nhk.or.jp/radioondemand/json/0164/bangumi_0164_01.json"
    },
    {
      "program_name": "特集オーディオドラマ",
      "url": "https://www.nhk.or.jp/radioondemand/json/P000025/bangumi_P000025_01.json"
    }
  ]
}
```

```
./radigo -i index.json
```

program_name/file_title/file_title_nn のようにディレクトリ作ってMP3形式で保村します。
例)
新日曜名作座/料理人季蔵捕物控/料理人季蔵捕物控_01.mp3
新日曜名作座/料理人季蔵捕物控/料理人季蔵捕物控_02.mp3
新日曜名作座/料理人季蔵捕物控/料理人季蔵捕物控_03.mp3

ダウンロード済みはprogram_name.txt というファイルを作成して、file_title_nn を書き込むことで管理します。
例)

新日曜名作座.txt

```
料理人季蔵捕物控_01
料理人季蔵捕物控_02
料理人季蔵捕物控_03
```

# Install

```
go get github.com/akokubu/radigo
```

[![Coverage Status](https://coveralls.io/repos/github/akokubu/radigo/badge.svg?branch=feature%2Fceveralls)](https://coveralls.io/github/akokubu/radigo?branch=feature%2Fceveralls)
[![CircleCI](https://circleci.com/gh/akokubu/radigo/tree/develop.svg?style=svg)](https://circleci.com/gh/akokubu/radigo/tree/develop)

2017/08/15
とりあえずコピペで動かすところまで実装。

2017/09/18
番組によってタイトルの取得方法違ったのでそれぞれに対応。

こんなことやりたい。
- [ ] 他で使えそうな部分切り出し
- [ ] json毎にループしているところを並列実行
- [ ] tsファイルを並列ダウンロードして連結
- [ ] ログ出力
- [x] indexファイルをjson化して設定とかも入れたい
- [ ] MP3ファイル出力先指定できるように
- [ ] テスト書く
- [ ] indexファイルに追加サブコマンド
- [ ] indexから削除サブコマンド
- [ ] indexの内容表示サブコマンド
- [ ] ダウンロード済み履歴管理
- [x] CI -> circleCIcで動かす様に。metalinter実行
- [x] テストカバレッジ -> coveralls
- [ ] ダウンロード進捗表示
- [ ] MP3変換進捗
- [ ] バージョン表記

