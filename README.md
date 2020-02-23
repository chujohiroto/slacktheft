# slack theft
Slackのパブリックチャンネルのメッセージを取得して、SQliteに保存するCLI Toolです。
また実行中は、RTM APIによってリアルタイムにメッセージが追記されます。

## Usage

```
$ bin/slacktheft -t=YOUR_SLACK_API_LEGACYTOKEN
```

```
NAME:
   Slacktheft - Slack messages dump and RTM

USAGE:
   slacktheft [global options] command [command options] [arguments...]

VERSION:
   0.0.2

AUTHOR:
   Hiroto Chujo <hiroto@irossoftware.com>

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --token value, -t value     SLACK API LEGACY TOKEN: (https://api.slack.com/legacy/custom-integrations/legacy-tokens) [$SLACK_API_TOKEN]
   --skipdump, -s              過去メッセージの取得プロセスをスキップします。
   --pagesize value, -p value  メッセージを取得する際の一回当たりのページングサイズを指定します。 (default: 10000)
   --help, -h                  show help
   --version, -v               print the version
```

## Docker

Dockerで使用する際は、DockerfileのYOUR_SLACK_API_LEGACYTOKENを https://api.slack.com/legacy/custom-integrations/legacy-tokens で取得したTokenに置き換えてください。
その後以下のコマンドを実行してください。

```
$ docker-compose build
$ docker-compose up -d
```

Docker上のSQLiteのファイルは、ホストマシンのフォルダのdocker/app/dump/dump.dbに保存されます。

## Todo
- MySQLの対応
- FileのSQL Blob対応
- FileのS3 Server対応