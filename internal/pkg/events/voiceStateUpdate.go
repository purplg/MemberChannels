package events

import (
	"github.com/bwmarrin/discordgo"
)

func (config *Events) VoiceStateUpdate(session *discordgo.Session, event *discordgo.VoiceStateUpdate) {
	log := config.Log.WithField("GuildID", event.GuildID)

	guildDB := config.DB.AsGuild(event.GuildID)

	if widget, ok := config.Widgets[guildDB.GuildID()]; ok {
		widget.UserVoiceEvent(event)
	} else {
		log.Errorln("Could not find widget for guild")
	}
}
