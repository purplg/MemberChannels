package events

import (
	"github.com/bwmarrin/discordgo"
)

func (config *Events) ChannelUpdate(session *discordgo.Session, event *discordgo.ChannelUpdate) {
	log := config.Log.WithField("GuildID", event.GuildID)
	guildDB := config.DB.AsGuild(event.GuildID)
	widget := config.Widgets[guildDB.GuildID()]
	if widget == nil {
		log.Errorln("Could not find widget for guild")
		return
	}

	if widget.IsManagedChannel(event.Channel) {
		widget.ManagedChannelChanged(event.Channel)
	}
}
