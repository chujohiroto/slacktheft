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
	app.Usage = "Slack messages dump and RTM"
	app.Version = "0.0.2"
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
		cli.IntFlag{
			Name:  "pagesize, p",
			Value: 10000,
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
			err = dumpRooms(token, context.Int("pagesize"), info.ID)
			if err != nil {
				fmt.Println(string(err.Error()))
				return err
			}
		}

		RTM(token, info.ID)
		return nil
	}

	app.Run(os.Args)
}
