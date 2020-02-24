package main

import (
	"fmt"
	"os"

	"github.com/slack-go/slack"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "Slacktheft"
	app.Usage = "SlackのメッセージをWebAPI と RTM APIを使ってリアルタイムにダンプするCLIツールです。"
	app.Version = "0.0.3"
	app.Author = "Hiroto Chujo"
	app.Email = "hiroto@irossoftware.com"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "token, t",
			Value:  "",
			Usage:  "SLACK API LEGACY TOKEN: (https://api.slack.com/legacy/custom-integrations/legacy-tokens)",
			EnvVar: "SLACK_API_TOKEN",
		},
		cli.BoolFlag{
			Name:  "skipdump, s",
			Usage: "過去メッセージの取得プロセスをスキップします。",
		},
		cli.BoolFlag{
			Name:  "private, p",
			Usage: "APIキー提供者のプライベートチャンネルを、#ChannelID#-privateのテーブルに保存します。 ただしRTM APIがPrivate Channelに対応してないので、リアルタイムのダンプはできません。",
		},
		cli.BoolFlag{
			Name:  "direct, d",
			Usage: "APIキー提供者のダイレクトメッセージ、個人メモを、#ChannelID#-directのテーブルに保存します。 ただしRTM APIがDMに対応してないので、リアルタイムのダンプはできません",
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "ログを詳細に出力します。",
		},
		cli.IntFlag{
			Name:  "pagesize, ps",
			Value: 1000,
			Usage: "メッセージを取得する際の一回当たりのページングサイズを指定します。",
		},
	}
	app.Action = func(context *cli.Context) error {
		token := context.String("token")

		if token == "" {
			fmt.Println("ERROR: the token flag is required...")
			fmt.Println("")
			cli.ShowAppHelp(context)
			os.Exit(2)
		}

		api := slack.New(token)
		_, err := api.AuthTest()

		if err != nil {
			fmt.Println("ERROR: Auth")
			return err
		}

		info, err := getTeamInfo(token)
		if err != nil {
			fmt.Println(string(err.Error()))
			return err
		}

		fmt.Println(info.Name)

		err = migrate(info.ID, info.Name)
		if err != nil {
			fmt.Println(string(err.Error()))
			return err
		}

		if context.Bool("skipdump") == false {
			err = dumpRooms(token, context.Int("pagesize"), info.ID, context.Bool("private"))
			if err != nil {
				fmt.Println(string(err.Error()))
				return err
			}
			if context.Bool("direct") {
				err = dumpUsers(token, context.Int("pagesize"), info.ID)
				if err != nil {
					fmt.Println(string(err.Error()))
					return err
				}
			}
		}

		RTM(token, info.ID, context.Bool("private"), context.Bool("direct"))
		return nil
	}

	app.Run(os.Args)
}
