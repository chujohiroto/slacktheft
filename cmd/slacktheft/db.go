package main

import (
	"database/sql"

	"github.com/slack-go/slack"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v1"
)

func migrate() (err error) {
	// Todo MYSQL も選択できるように
	db, err := sql.Open("sqlite3", "./dump.db")
	if err != nil {
		return err
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	dbmap.AddTableWithName(Message{}, "message").SetKeys(false, "Timestamp")
	dbmap.DropTables()
	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		return err
	}

	defer db.Close()

	return nil
}

func insert(message slack.Message) (err error) {
	ev := slack.MessageEvent(message)

	// mapする　本当にダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイダルイ
	msg := mappedModel(ev)
	// Todo MYSQL も選択できるように
	db, err := sql.Open("sqlite3", "./dump.db")
	if err != nil {
		return err
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	dbmap.AddTableWithName(Message{}, "message").SetKeys(false, "Timestamp")

	var pr Message
	err = dbmap.SelectOne(&pr, "select * from message where Timestamp=?", msg.Timestamp)

	if err == nil {
		return nil
	}

	err = dbmap.Insert(&msg)

	if err != nil {
		return err
	}

	defer db.Close()

	return nil
}

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
