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
		widgetData.ListenChannelName = config.Vars.DefaultChannelName
	}

	if w, err := widget.New(session, log, guildDB, widgetData); err != nil {
		log.WithError(err).Warnln("Error creating widget")
	} else {
		config.Widgets[guildDB.GuildID()] = w
	}
}
