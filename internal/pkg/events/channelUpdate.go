package events

import (
	"github.com/bwmarrin/discordgo"
)

func (config *Events) ChannelUpdate(session *discordgo.Session, event *discordgo.ChannelUpdate) {
	log := config.Log.WithField("GuildID", event.GuildID)

	widget, ok := config.Widgets[event.GuildID]
	if !ok {
		log.Errorln("Could not find widget for guild")
		return
	}

	if widget.IsListenChannel(event.Channel.ID) {
		widget.RenameListenChannel(event.Channel.Name)
	} else if widget.IsMemberChannel(event.Channel.ID) {
		widget.RenameMemberChannel(event.Channel.ID, event.Channel.Name)
	}
}
