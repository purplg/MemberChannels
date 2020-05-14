package events

import (
	"github.com/bwmarrin/discordgo"
)

func (config *Events) VoiceStateUpdate(session *discordgo.Session, event *discordgo.VoiceStateUpdate) {
	log := config.Log.WithField("GuildID", event.GuildID)

	guildDB := config.DB.AsGuild(event.GuildID)
	widget := config.Widgets[guildDB.GuildID()]
	if widget == nil {
		log.Errorln("Could not find widget for guild")
		return
	}

	user, err := session.User(event.UserID)
	if err != nil {
		log.WithError(err).WithField("UserID", user.ID).Warnln("Could not find user")
		return
	}

	if event.ChannelID == "" {
		widget.UserDisconnected(user)
		return
	}

	channel, err := session.Channel(event.ChannelID)
	if err != nil {
		log.WithError(err).WithField("ChannelID", channel.ID).Warnln("Could not find channel")
		return
	}

	if widget.IsManagedChannel(channel) {
		widget.UserJoinedManagedChannel(user, channel)
	} else {
		widget.UserDisconnected(user)
	}
}
