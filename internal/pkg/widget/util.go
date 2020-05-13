package widget

import "github.com/bwmarrin/discordgo"

func (w *Widget) getCategory(categoryID, defaultName string) (*discordgo.Channel, error) {
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

func (w *Widget) getListenChannel(channelID, defaultName, parentID string) (*discordgo.Channel, error) {
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
