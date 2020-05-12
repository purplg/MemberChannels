package mock

import "github.com/bwmarrin/discordgo"

func mockGuild(id string, name string) *discordgo.Guild {
	guild := &discordgo.Guild{
		ID:       id,
		Name:     name,
		Members:  mockMembers(id),
		Roles:    mockRoles(),
		Channels: mockChannels(id),
	}
	guild.MemberCount = len(guild.Members)
	return guild
}
