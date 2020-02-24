# slack theft
Slackのメッセージを取得して、SQliteに保存するCLI Toolです。
実行中は、RTM APIによってリアルタイムにメッセージが追記されます。

## Usage

クローンして、ビルドするのだ！ Golangはそれなりに新しいのインストールしといてね
```
$ git clone https://github.com/ChujoHiroto/slacktheft
$ cd slacktheft
$ make build
```

```
$ bin/slacktheft -t=YOUR_SLACK_API_LEGACYTOKEN
```

プライベートチャンネルやダイレクトメッセージを含める場合
```
$ bin/slacktheft -p -d -t=YOUR_SLACK_API_LEGACYTOKEN
```

過去メッセージの取得を行わない場合
```
$ bin/slacktheft -s -t=YOUR_SLACK_API_LEGACYTOKEN
```

```
NAME:
   Slacktheft - SlackのメッセージをWebAPI と RTM APIを使ってリアルタイムにダンプするCLIツールです。

USAGE:
   slacktheft [global options] command [command options] [arguments...]

VERSION:
   0.0.3

AUTHOR:
   Hiroto Chujo <hiroto@irossoftware.com>

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --token value, -t value       SLACK API LEGACY TOKEN: (https://api.slack.com/legacy/custom-integrations/legacy-tokens) [$SLACK_API_TOKEN]
   --skipdump, -s                過去メッセージの取得プロセスをスキップします。
   --private, -p                 APIキー提供者のプライベートチャンネルを、#ChannelID#-privateのテーブルに保存します。 ただしRTM APIがPrivate Channelに対応してないので、リアルタイムのダンプはできません。
   --direct, -d                  APIキー提供者のダイレクトメッセージ、個人メモを、#ChannelID#-directのテーブルに保存します。 ただしRTM APIがDMに対応してないので、リアルタイムのダンプはできません
   --verbose                     ログを詳細に出力します。
   --pagesize value, --ps value  メッセージを取得する際の一回当たりのページングサイズを指定します。 (default: 1000)
   --help, -h                    show help
   --version, -v                 print the version
```

## Docker

Dockerで使用する際は、DockerfileのYOUR_SLACK_API_LEGACYTOKENを https://api.slack.com/legacy/custom-integrations/legacy-tokens で取得したTokenに置き換えてください。
その後以下のコマンドを実行してください。

```
$ docker-compose build
$ docker-compose up
```

Docker上のSQLiteのファイルは、ホストマシンのフォルダのdocker/app/dump/dump.dbに保存されます。

## 既知の問題点
RTM APIがPrivate Channel, DMに対応してないので、リアルタイムのダンプができない。（Event APIで可能な模様）

## Todo
- Event APIでのリアルタイムメッセージ取得の実装
- MySQLの対応
- FileのSQL Blob対応
- FileのS3 Server対応
