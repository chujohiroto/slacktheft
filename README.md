# slack theft
Slackのパブリックチャンネルのメッセージを取得して、SQliteに保存するCLI Toolです。
また実行中は、RTM APIによってリアルタイムにメッセージが追記されます。

## Usage

```
$ bin/slacktheft -t=SLACKAPILEGACYTOKEN
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

## Todo
- 複数ワークスペースの対応（テーブル名をワークスペースの名前にすべきだ...）
- MySQLの対応
- FileのSQL Blob対応
- FileのS3 Server対応
- Dockerfile作成