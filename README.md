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
   0.0.1

AUTHOR:
   Hiroto Chujo <hiroto@irossoftware.com>

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --token value, -t value  a Slack API token: (see: https://api.slack.com/web) [$SLACK_API_TOKEN]
   --help, -h               show help
   --version, -v            print the version
```

## Todo
- MySQLの対応
- FileのSQL Blob対応
- FileのS3 Server対応
- Dockerfile作成