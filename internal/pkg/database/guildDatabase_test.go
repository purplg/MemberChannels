package database

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"purplg.com/memberchannels/internal/pkg/mock"
)

func Test_GuildDatabase(t *testing.T) {
	logger := logrus.New()
	log := logrus.NewEntry(logger)
	tempDB := os.TempDir() + "/" + mock.Test_DBFile
	db, err := Open(tempDB, log)
	if err != nil {
		log.WithError(err).Fatalln()
	}

	defer db.Close()
	defer os.RemoveAll(tempDB)

	gdb := db.AsGuild(mock.Test_GuildID)

	testString := mock.RandomString()
	if gdb.SetCategoryID(testString); gdb.CategoryID() != testString {
		t.Error("CategoryID failed")
	}
	testString = mock.RandomString()
	if gdb.SetCategoryName(testString); gdb.CategoryName() != testString {
		t.Error("CategoryName failed")
	}
	testString = mock.RandomString()
	if gdb.SetChannelID(testString); gdb.ChannelID() != testString {
		t.Error("ChannelID failed")
	}
	testString = mock.RandomString()
	if gdb.SetChannelName(testString); gdb.ChannelName() != testString {
		t.Error("ChannelName failed")
	}
}
