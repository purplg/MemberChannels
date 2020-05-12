package events

import (
	"github.com/bwmarrin/discordgo"
	"purplg.com/memberchannels/internal/pkg/widget"
)

func (config *Events) GuildCreated(session *discordgo.Session, event *discordgo.GuildCreate) {
	log := config.Log.WithField("GuildID", event.Guild.ID)

	guildDB := config.DB.AsGuild(event.Guild.ID)

	w := widget.New(session, log, guildDB, config.Vars.DefaultCategoryName, config.Vars.DefaultChannelName)
	config.Widgets[guildDB.GuildID()] = w

	if err := w.Show(); err != nil {
		log.WithError(err).Errorln("Error creating widget")
		return
	}
}
