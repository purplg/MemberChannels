package events

import (
	"github.com/bwmarrin/discordgo"
)

func (config *Events) ChannelUpdate(session *discordgo.Session, event *discordgo.ChannelUpdate) {
	guildDB := config.DB.AsGuild(event.GuildID)
	categoryID := guildDB.CategoryID()

	// Rename category
	if event.Channel.ID == categoryID {
		guildDB.SetCategoryName(event.Channel.Name)
		return
	}

	// Otherwise, we don't care about channels without a category (aka parent)
	if event.Channel.ParentID != categoryID {
		return
	}

	// Okay, now we know we got a managed channel
	channelID := guildDB.ChannelID()

	// Rename the listening channel
	if event.Channel.ID == channelID {
		guildDB.SetChannelName(event.Channel.Name)
		return
	}

	// Otherwise, it's a user channel
	userID := guildDB.ChannelOwner(event.Channel.ID)
	guildDB.SetUserChannel(userID, event.Channel.ID, event.Channel.Name)
}
