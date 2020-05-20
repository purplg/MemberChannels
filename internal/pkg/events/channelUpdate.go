package events

import (
	"github.com/bwmarrin/discordgo"
)

func (config *Events) ChannelUpdate(session *discordgo.Session, event *discordgo.ChannelUpdate) {
	log := config.Log.WithField("GuildID", event.GuildID)
	guildDB := config.DB.AsGuild(event.GuildID)

	if widget, ok := config.Widgets[guildDB.GuildID()]; ok {
		if widget.IsManagedChannel(event.Channel) {
			widget.ChannelChangedEvent(event.Channel)
		}
	} else {
		log.Errorln("Could not find widget for guild")
	}
}
