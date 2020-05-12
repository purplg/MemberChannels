package widget

import "github.com/bwmarrin/discordgo"

func (w *Widget) getCategory() (*discordgo.Channel, error) {
	if w.categoryID != "" {
		if category, err := w.session.State.Channel(w.categoryID); err == nil {
			return category, nil
		}

		if category, err := w.session.Channel(w.categoryID); err == nil {
			return category, nil
		}
	}

	return w.session.GuildChannelCreateComplex(w.guildDB.GuildID(), discordgo.GuildChannelCreateData{
		Name:      w.categoryName,
		Type:      discordgo.ChannelTypeGuildCategory,
		UserLimit: 1,
	})
}

func (w *Widget) getChannel(parentID string) (*discordgo.Channel, error) {
	if w.channelID != "" {
		if channel, err := w.session.State.Channel(w.channelID); err == nil {
			return channel, nil
		}
		if channel, err := w.session.Channel(w.channelID); err == nil {
			return channel, nil
		}
	}

	return w.session.GuildChannelCreateComplex(w.guildDB.GuildID(), discordgo.GuildChannelCreateData{
		Name:      w.channelName,
		Type:      discordgo.ChannelTypeGuildVoice,
		UserLimit: 1,
		ParentID:  parentID,
	})
}
