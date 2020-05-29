package events

import (
	"github.com/bwmarrin/discordgo"
	"purplg.com/memberchannels/internal/pkg/widget"
)

func (config *Events) GuildCreated(session *discordgo.Session, event *discordgo.GuildCreate) {
	log := config.Log.WithField("GuildID", event.Guild.ID)

	guildDB := config.DB.AsGuild(event.Guild.ID)

	categoryID := guildDB.CategoryID()
	categoryName := guildDB.CategoryName()
	listenChannelName := guildDB.ChannelName()

	if categoryName == "" {
		categoryName = config.Vars.DefaultCategoryName
	}

	if listenChannelName == "" {
		listenChannelName = config.Vars.DefaultListenName
	}

	widget := widget.New(session, log, guildDB)

	if err := widget.Spawn(categoryID, categoryName, listenChannelName); err != nil {
		log.WithError(err).Warnln("Error spawning widget")
	} else {
		config.Widgets[event.Guild.ID] = widget
	}
}
