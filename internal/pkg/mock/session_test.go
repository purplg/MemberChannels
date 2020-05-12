package mock

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

func Test_MockSession(t *testing.T) {
	session := MockSession()

	_, err := session.User(Test_UserID)
	if err != nil {
		t.Error(err)
	}

	_, err = session.Guild(Test_GuildID)
	if err != nil {
		t.Error(err)
	}

	_, err = session.Channel(Test_ChannelID)
	if err != nil {
		t.Error(err)
	}

	_, err = session.GuildChannelCreate(Test_GuildID, Test_ChannelName, discordgo.ChannelTypeGuildCategory)
	if err != nil {
		t.Error(err)
	}

	defer session.Close()
}
