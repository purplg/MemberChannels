package widget

import "github.com/bwmarrin/discordgo"

func (w *widget) getCategory(categoryID, defaultName string) (*discordgo.Channel, error) {
	if w.categoryChannel == nil {
		if category, err := w.session.State.Channel(categoryID); err == nil {
			return category, nil
		}

		if category, err := w.session.Channel(categoryID); err == nil {
			return category, nil
		}
	}

	return w.session.GuildChannelCreateComplex(w.guildDB.GuildID(), discordgo.GuildChannelCreateData{
		Name:      defaultName,
		Type:      discordgo.ChannelTypeGuildCategory,
		UserLimit: 1,
	})
}

func (w *widget) getListenChannel(channelID, defaultName, parentID string) (*discordgo.Channel, error) {
	if w.listenChannel == nil {
		if channel, err := w.session.State.Channel(channelID); err == nil {
			return channel, nil
		}
		if channel, err := w.session.Channel(channelID); err == nil {
			return channel, nil
		}
	}

	return w.session.GuildChannelCreateComplex(w.guildDB.GuildID(), discordgo.GuildChannelCreateData{
		Name:     defaultName,
		Type:     discordgo.ChannelTypeGuildVoice,
		ParentID: parentID,
	})
}

func (w *widget) isManagedChannel(channel *discordgo.Channel) bool {
	return channel.ParentID == w.categoryChannel.ID
}

func (w *widget) isListenChannel(channel *discordgo.Channel) bool {
	return channel.ID == w.listenChannel.ID
}

func (w *widget) isUserChannel(channel *discordgo.Channel) bool {
	return w.isManagedChannel(channel) && !w.isListenChannel(channel)
}
