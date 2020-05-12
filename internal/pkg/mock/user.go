package mock

import "github.com/bwmarrin/discordgo"

func mockUser(id string, name string) *discordgo.User {
	return &discordgo.User{
		ID:       id,
		Username: name,
	}
}
