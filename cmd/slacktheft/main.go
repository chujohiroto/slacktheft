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
	app.Version = "0.0.1"
	app.Author = "Hiroto Chujo"
	app.Email = "hiroto@irossoftware.com"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "token, t",
			Value:  "",
			Usage:  "a Slack API token: (see: https://api.slack.com/web)",
			EnvVar: "SLACK_API_TOKEN",
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
			os.Exit(2)
		}

		err = migrate()
		if err != nil {
			fmt.Println(string(err.Error()))
			return err
		}

		err = dumpRooms(token)
		if err != nil {
			fmt.Println(string(err.Error()))
			return err
		}

		RTM(token)
		return nil
	}

	app.Run(os.Args)
}
