package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/slack-go/slack"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v1"
)

func exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

var globaldbmap gorp.DbMap
var isOpenDB bool

func openDBMap() (gorp.DbMap, error) {

	if isOpenDB {
		return globaldbmap, nil
	}

	// Todo MYSQL も選択できるように
	db, err := sql.Open("sqlite3", "./dump/dump.db")
	if err != nil {
		return gorp.DbMap{}, err
	}
	globaldbmap = gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	isOpenDB = true
	return globaldbmap, nil
}

func migrate(workspaceid string, workspacename string) (err error) {
	if !exists("dump") {
		if err := os.Mkdir("dump", 0777); err != nil {
			return err
		}
	}

	err = createTable(workspaceid)
	err = createTable(workspaceid + "-private")
	err = createTable(workspaceid + "-direct")

	if err != nil {
		return err
	}

	insertWorkspace(workspaceid, workspacename)
	return nil
}

func createTable(name string) error {
	dbmap, err := openDBMap()

	dbmap.AddTableWithName(Message{}, name).SetKeys(false, "Timestamp")

	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		return err
	}

	return nil
}

func insertWorkspace(workspaceid string, workspacename string) error {
	dbmap, err := openDBMap()

	dbmap.AddTableWithName(Workspace{}, "Workspaces").SetKeys(false, "ID")
	dbmap.DropTables()
	err = dbmap.CreateTablesIfNotExists()
	dbmap.Insert(&Workspace{ID: workspaceid, Name: workspacename})
	if err != nil {
		return err
	}

	return nil
}

func insert(message slack.Message, workspaceid string) (err error) {
	ev := slack.MessageEvent(message)

	/*
		SlackAPIのtypeとDBのtypeでmappingする
		本当にダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイ
		ダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイ
		ダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイ
	*/
	msg := mappedModel(ev)

	dbmap, err := openDBMap()

	dbmap.AddTableWithName(Message{}, workspaceid).SetKeys(false, "Timestamp")

	var pr Message
	err = dbmap.SelectOne(&pr, "select * from ? where Timestamp=?", workspaceid, msg.Timestamp)

	if err == nil {
		return nil
	}

	err = dbmap.Insert(&msg)

	if err != nil {
		fmt.Println(string(err.Error()))
		return nil
	}

	return nil
}

func insertPrivate(message slack.Message, workspaceid string) (err error) {
	ev := slack.MessageEvent(message)
	msg := mappedModel(ev)

	dbmap, err := openDBMap()

	dbmap.AddTableWithName(Message{}, workspaceid+"-private").SetKeys(false, "Timestamp")

	var pr Message
	err = dbmap.SelectOne(&pr, "select * from ?-private where Timestamp=?", workspaceid, msg.Timestamp)

	if err == nil {
		return nil
	}

	err = dbmap.Insert(&msg)

	if err != nil {
		fmt.Println(string(err.Error()))
		return nil
	}

	return nil
}

func insertDirect(message slack.Message, workspaceid string) (err error) {
	ev := slack.MessageEvent(message)
	msg := mappedModel(ev)

	dbmap, err := openDBMap()
	dbmap.AddTableWithName(Message{}, workspaceid+"-direct").SetKeys(false, "Timestamp")

	var pr Message
	err = dbmap.SelectOne(&pr, "select * from ?-direct where Timestamp=?", workspaceid, msg.Timestamp)

	if err == nil {
		return nil
	}

	err = dbmap.Insert(&msg)

	if err != nil {
		fmt.Println(string(err.Error()))
		return nil
	}

	return nil
}

// reflectionが使えればこの辺マシになるはずだが、、、
func mappedModel(message slack.MessageEvent) Message {
	msg := Message{}
	msg.BotID = message.BotID
	msg.Channel = message.Channel
	msg.DeleteOriginal = message.DeleteOriginal
	msg.DeletedTimestamp = message.DeletedTimestamp
	msg.EventTimestamp = message.EventTimestamp
	msg.Hidden = message.Hidden
	msg.Inviter = message.Inviter
	msg.IsStarred = message.IsStarred
	msg.ItemType = message.ItemType
	msg.LastRead = message.LastRead
	msg.Name = message.Name
	msg.OldName = message.OldName
	msg.ParentUserId = message.ParentUserId
	msg.Purpose = message.Purpose
	msg.ReplaceOriginal = message.ReplaceOriginal
	msg.ReplyCount = message.ReplyCount
	msg.ReplyTo = message.ReplyTo
	msg.ResponseType = message.ResponseType
	msg.SubType = message.SubType
	msg.Subscribed = message.Subscribed
	msg.Team = message.Team
	msg.Text = message.Text
	msg.ThreadTimestamp = message.ThreadTimestamp
	msg.Timestamp = message.Timestamp
	msg.Topic = message.Topic
	msg.Type = message.Type
	msg.UnreadCount = message.UnreadCount
	msg.Upload = message.Upload
	msg.User = message.User
	msg.Username = message.Username

	if message.Icons != nil {
		msg.IconEmoji = message.Icons.IconEmoji
		msg.IconURL = message.Icons.IconURL
	}

	return msg
}
