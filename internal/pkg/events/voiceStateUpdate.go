package events

import (
	"github.com/bwmarrin/discordgo"
)

func (config *Events) VoiceStateUpdate(session *discordgo.Session, event *discordgo.VoiceStateUpdate) {
	log := config.Log.WithField("GuildID", event.GuildID)

	guildDB := config.DB.AsGuild(event.GuildID)

	widget, ok := config.Widgets[guildDB.GuildID()]
	if !ok {
		log.Errorln("Could not find widget for guild")
		return
	}

	widget.UserLeft(event.UserID)

	if event.ChannelID == "" {
		return
	}

	if widget.IsListenChannel(event.ChannelID) {
		widget.UserRequestChannel(event.UserID)
	} else {
		widget.UserJoined(event)
	}
}
