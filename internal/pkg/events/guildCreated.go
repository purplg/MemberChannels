package events

import (
	"github.com/bwmarrin/discordgo"
	"purplg.com/memberchannels/internal/pkg/widget"
)

func (config *Events) GuildCreated(session *discordgo.Session, event *discordgo.GuildCreate) {
	log := config.Log.WithField("GuildID", event.Guild.ID)

	guildDB := config.DB.AsGuild(event.Guild.ID)

	widgetData := &widget.WidgetData{
		CategoryID:        guildDB.CategoryID(),
		CategoryName:      guildDB.CategoryName(),
		ListenChannelID:   guildDB.ChannelID(),
		ListenChannelName: guildDB.ChannelName(),
	}

	w, err := widget.New(session, log, guildDB, widgetData)
	if err != nil {
		log.WithError(err).Warnln("Error creating widget")
		return
	}
	config.Widgets[guildDB.GuildID()] = w
}
