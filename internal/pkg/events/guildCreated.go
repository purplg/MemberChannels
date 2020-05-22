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

	if widgetData.CategoryName == "" {
		widgetData.CategoryName = config.Vars.DefaultCategoryName
	}

	if widgetData.ListenChannelName == "" {
		widgetData.ListenChannelName = config.Vars.DefaultListenName
	}

	widget := widget.New(session, log, guildDB)

	if err := widget.Spawn(widgetData); err != nil {
		log.WithError(err).Warnln("Error spawning widget")
	} else {
		config.Widgets[event.Guild.ID] = widget
	}
}
