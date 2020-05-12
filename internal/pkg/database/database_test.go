package database

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"purplg.com/memberchannels/internal/pkg/mock"
)

func Test_Database(t *testing.T) {
	logger := logrus.New()
	log := logrus.NewEntry(logger)
	tempDB := os.TempDir() + "/" + mock.Test_DBFile
	db, err := Open(tempDB, log)
	if err != nil {
		log.WithError(err).Fatalln()
	}

	defer db.Close()
	defer os.RemoveAll(tempDB)

	db.AsGuild(mock.Test_GuildID)
}
