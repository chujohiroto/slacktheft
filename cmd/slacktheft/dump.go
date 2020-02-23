package main

import (
	"fmt"

	"github.com/slack-go/slack"
)

func dumpRooms(token string) error {
	api := slack.New(
		token,
		slack.OptionDebug(false),
		//slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)

	// Dump Channels
	fmt.Println("dump public channel")
	_, err := dumpChannels(api)

	if err != nil {
		return err
	}

	// Channelのテーブル定義
	/*for _, channel := range channels {
		fmt.Printf("Channel ID: " + channel.ID)
	}*/

	// Dump Private Groups
	/*
		fmt.Println("dump private channel")
		groups := dumpGroups(api, dir, rooms)

		if len(groups) > 0 {
			for _, group := range groups {
				channel := slack.Channel{}
				channel.ID = group.ID
				channel.Name = group.Name
				channel.Created = group.Created
				channel.Creator = group.Creator
				channel.IsArchived = group.IsArchived
				channel.IsChannel = true
				channel.IsGeneral = false
				channel.IsMember = true
				channel.LastRead = group.LastRead
				channel.Latest = group.Latest
				channel.Members = group.Members
				channel.NumMembers = group.NumMembers
				channel.Purpose = group.Purpose
				channel.Topic = group.Topic
				channel.UnreadCount = group.UnreadCount
				channel.UnreadCountDisplay = group.UnreadCountDisplay
				channels = append(channels, channel)
			}
		}*/

	return nil
}

func dumpChannels(api *slack.Client) ([]slack.Channel, error) {
	// Todo アーカイブ除去フラグのオプション実装
	channels, err := api.GetChannels(false)

	if err != nil {
		return nil, err
	}

	if len(channels) == 0 {
		var channels []slack.Channel
		return channels, nil
	}

	for _, channel := range channels {
		fmt.Println("dump channel " + channel.Name)
		err = dumpChannel(api, channel.ID, channel.Name, "channel")
		if err != nil {
			return nil, err
		}
	}

	return channels, nil
}

func dumpChannel(api *slack.Client, id, name, channelType string) error {
	var messages []slack.Message

	// Todo Group、Private Channel、DMのオプション実装
	/*if channelType == "group" {
		channelPath = path.Join("private_channel", name)
		messages = fetchGroupHistory(api, id)
	} else if channelType == "dm" {
		channelPath = path.Join("direct_message", name)
		messages = fetchDirectMessageHistory(api, id)
	} else {*/
	//channelPath = path.Join("channel", name)

	messages, err := fetchChannelHistory(api, id)
	if err != nil {
		return err
	}

	//}

	if len(messages) == 0 {
		fmt.Printf("Zero Message")
		return nil
	}

	//もし時系列ソートの実装したいときに
	//sort.Sort(byTimestamp(messages))

	for _, message := range messages {
		message.Channel = name
		err = insert(message)
		if err != nil {
			return err
		}
	}

	return nil
}

func fetchChannelHistory(api *slack.Client, ID string) ([]slack.Message, error) {
	historyParams := slack.NewHistoryParameters()
	historyParams.Count = 10000

	// Fetch History
	history, err := api.GetChannelHistory(ID, historyParams)

	if err != nil {
		return nil, err
	}

	messages := history.Messages
	latest := messages[len(messages)-1].Timestamp
	for {
		if history.HasMore != true {
			break
		}

		historyParams.Latest = latest
		history, err = api.GetChannelHistory(ID, historyParams)

		if err != nil {
			return nil, err
		}

		length := len(history.Messages)
		if length > 0 {
			latest = history.Messages[length-1].Timestamp
			messages = append(messages, history.Messages...)
		}

	}

	return messages, nil
}
