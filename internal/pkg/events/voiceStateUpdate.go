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
	channel, _ := session.Channel(event.ChannelID)

	widget.UserSwitchedChannel(user, channel)
}
