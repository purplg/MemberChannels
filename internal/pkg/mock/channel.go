package mock

import "github.com/bwmarrin/discordgo"

func mockChannel(guildid string, id string, name string) *discordgo.Channel {
	return &discordgo.Channel{
		GuildID: guildid,
		ID:   id,
		Name: name,
	}
}

func mockChannels(guildid string) []*discordgo.Channel {
	return []*discordgo.Channel{
		mockChannel(guildid, Test_ChannelID,  Test_ChannelName),
		mockChannel(guildid, Test_ChannelID2, Test_ChannelName2),
		mockChannel(guildid, Test_ChannelID3, Test_ChannelName3),
	}
}
