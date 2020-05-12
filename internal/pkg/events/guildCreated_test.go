package events

import (
	"os"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"purplg.com/memberchannels/internal/pkg/database"
	"purplg.com/memberchannels/internal/pkg/mock"
	"purplg.com/memberchannels/internal/pkg/variables"
)

var TEMP_DIR = os.TempDir() + "/memberchannels_testdb"

func Test_GuildCreated(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	vars := variables.New(mock.Test_Token)

	tempDB := os.TempDir() + "/" + mock.Test_DBFile
	db, err := database.Open(tempDB, logrus.NewEntry(logger))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer os.RemoveAll(tempDB)

	config := &Config{
		Vars: vars,
		Log:  logrus.NewEntry(logger),
		DB:   db,
	}
	defer config.Close()

	log := logrus.NewEntry(logger)
	log.WithFields(logrus.Fields{
		"config.Vars": config.Vars,
		"config.Log": config.Log,
		"config.DB": config.DB,
	})

	session := mock.MockSession()
	defer session.Close()

	mockGuild := &discordgo.Guild{
		ID: mock.Test_GuildID,
	}

	event := &discordgo.GuildCreate{
		Guild: mockGuild,
	}

	config.GuildCreated(session, event)

	guildDB := config.DB.AsGuild(event.Guild.ID)

	config.Log.WithFields(logrus.Fields{
		"guildID":      guildDB.GuildID,
		"categoryID":   guildDB.CategoryID(),
		"categoryName": guildDB.CategoryName(),
		"channelID":    guildDB.ChannelID(),
		"channelName":  guildDB.ChannelName(),
	}).Debugln("Current State")

	if guildDB.GuildID != mock.Test_GuildID {
		t.Error("Invalid guildID")
	}
}
