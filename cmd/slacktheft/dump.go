package main

import (
	"fmt"

	"github.com/slack-go/slack"
)

func getTeamInfo(token string) (*slack.TeamInfo, error) {
	api := slack.New(
		token,
		slack.OptionDebug(false),
	)
	info, err := api.GetTeamInfo()

	if err != nil {
		return nil, err
	}

	return info, nil
}

func dumpUsers(token string, count int, workspaceid string) error {
	api := slack.New(
		token,
		slack.OptionDebug(false),
		//slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)

	fmt.Println("dump user information")
	users, err := api.GetUsers()
	if err != nil {
		return err
	}

	fmt.Println("dump direct message")
	ims, err := api.GetIMChannels()
	//fmt.Println(ims)

	for _, im := range ims {
		for _, user := range users {
			if im.User == user.ID {
				fmt.Println("dump DM with " + user.Name)
				err = dumpChannel(api, im.ID, user.Name, "dm", count, workspaceid, true)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func dumpRooms(token string, count int, workspaceid string, isprivate bool) error {
	api := slack.New(
		token,
		slack.OptionDebug(false),
		//slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)
	// Dump Channels
	fmt.Println("dump public channel")
	_, err := dumpChannels(api, count, workspaceid, isprivate)

	if err != nil {
		return err
	}

	// Dump Private Groups
	fmt.Println("dump private channel")
	_, err = dumpGroups(api, count, workspaceid, isprivate)

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

func dumpGroups(api *slack.Client, count int, workspaceid string, isprivate bool) ([]slack.Group, error) {
	groups, err := api.GetGroups(false)

	if err != nil {
		return nil, err
	}

	if len(groups) == 0 {
		var groups []slack.Group
		return groups, nil
	}

	for _, group := range groups {
		fmt.Println("dump channel " + group.Name)
		dumpChannel(api, group.ID, group.Name, "group", count, workspaceid, isprivate)
	}

	return groups, nil
}

func dumpChannels(api *slack.Client, count int, workspaceid string, isprivate bool) ([]slack.Channel, error) {
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
		err = dumpChannel(api, channel.ID, channel.Name, "channel", count, workspaceid, isprivate)
		if err != nil {
			return nil, err
		}
	}

	return channels, nil
}

func dumpChannel(api *slack.Client, id, name, channelType string, count int, workspaceid string, isprivate bool) error {
	var messages []slack.Message
	var err error

	// Todo Group、Private Channel、DMのオプション実装
	if channelType == "group" && isprivate {
		messages, err = fetchGroupHistory(api, id, count)
		if err != nil {
			return err
		}
	} else if channelType == "dm" && isprivate {
		messages, err = fetchDirectMessageHistory(api, id, count)
		if err != nil {
			return err
		}
	} else {
		messages, err = fetchChannelHistory(api, id, count)
		if err != nil {
			return err
		}
	}

	if len(messages) == 0 {
		fmt.Printf("Zero Message")
		return nil
	}

	// もし時系列ソートの実装したいときに まあtimestampをキーにしている間は実装しなくていいでしょう
	// sort.Sort(byTimestamp(messages))

	for _, message := range messages {
		message.Channel = name
		if channelType == "group" {
			err := insertPrivate(message, workspaceid)
			if err != nil {
				return err
			}
		} else if channelType == "dm" {
			err := insertDirect(message, workspaceid)
			if err != nil {
				return err
			}
		} else {
			err := insert(message, workspaceid)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func fetchChannelHistory(api *slack.Client, ID string, count int) ([]slack.Message, error) {
	historyParams := slack.NewHistoryParameters()
	historyParams.Count = count

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

func fetchGroupHistory(api *slack.Client, ID string, count int) ([]slack.Message, error) {
	historyParams := slack.NewHistoryParameters()
	historyParams.Count = count

	// Fetch History
	history, err := api.GetGroupHistory(ID, historyParams)

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
		history, err = api.GetGroupHistory(ID, historyParams)
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

func fetchDirectMessageHistory(api *slack.Client, ID string, count int) ([]slack.Message, error) {
	historyParams := slack.NewHistoryParameters()
	historyParams.Count = count

	// Fetch History
	history, err := api.GetIMHistory(ID, historyParams)
	messages := history.Messages
	if len(messages) == 0 {
		return messages, nil
	}
	latest := messages[len(messages)-1].Timestamp
	for {
		if history.HasMore != true {
			break
		}

		historyParams.Latest = latest
		history, err = api.GetIMHistory(ID, historyParams)

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
