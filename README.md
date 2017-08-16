radigo
====

Overview

らじる聴き逃し番組MP3ダウンロード

Go勉強用に作ってみた。

# Usage
index.txt ファイル作成

```
https://www.nhk.or.jp/radioondemand/json/0058/bangumi_0058_01.json
https://www.nhk.or.jp/radioondemand/json/0930/bangumi_0930_01.json
https://www.nhk.or.jp/radioondemand/json/0164/bangumi_0164_01.json
```

```
./radigo -i index.txt
```

program_name/file_title/file_title(n) のようにディレクトリ作ってMP3形式で保存します。<br/>
例)<br/>
新日曜名作座/料理人季蔵捕物控/料理人季蔵捕物控(1).mp3<br/>
新日曜名作座/料理人季蔵捕物控/料理人季蔵捕物控(2).mp3<br/>
新日曜名作座/料理人季蔵捕物控/料理人季蔵捕物控(3).mp3<br/>

ダウンロード済みはprogram_name.txt というファイルを作成して、file_title(n) を書き込むことで管理します。<br/>
例)<br/>
新日曜名作座.txt

```
料理人季蔵捕物控(1)
料理人季蔵捕物控(2)
料理人季蔵捕物控(3)
```

# Install

```
go get github.com/akokubu/radigo
```

2017/08/15
とりあえずコピペで動かすところまで実装。

こんなことやりたい。
* 他で使えそうな部分切り出し
* json毎にループしているところを並列実行
* tsファイルを並列ダウンロードして連結
* ログ出力
* indexファイルをjson化して設定とかも入れたい
* MP3ファイル出力先指定できるように
* テスト書く
* indexファイルに追加サブコマンド
* indexから削除サブコマンド
* indexの内容表示サブコマンド
* ダウンロード済み履歴管理
* CI
* テストカバレッジ
* ダウンロード進捗表示
* MP3変換進捗
* バージョン表記

