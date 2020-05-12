package mock

import "github.com/bwmarrin/discordgo"

const (
)

func mockMember(guildid string, id string, name string) *discordgo.Member {
	return &discordgo.Member{
		GuildID: guildid,
		User:    mockUser(id, name),
	}
}

func mockMembers(guildid string) []*discordgo.Member {
	return []*discordgo.Member{
		mockMember(guildid, Test_UserID,  Test_UserName),
		mockMember(guildid, Test_UserID2, Test_UserName2),
		mockMember(guildid, Test_UserID3, Test_UserName3),
	}
}
