package widget

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"purplg.com/memberchannels/internal/pkg/database"
	"purplg.com/memberchannels/internal/pkg/mock"
)

func Test_Widget(t *testing.T) {
	session := mock.MockSession()
	log := logrus.NewEntry(logrus.New())
	tempDB := os.TempDir() + "/" + mock.Test_DBFile
	db, err := database.Open(tempDB, log)
	if err != nil {
		log.WithError(err).Fatalln()
	}

	defer db.Close()
	defer os.RemoveAll(tempDB)

	guildDB := db.AsGuild(mock.Test_GuildID)

	New(session, log, guildDB, mock.Test_ChannelName, mock.Test_ChannelName2)
}
