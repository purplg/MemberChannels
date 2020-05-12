package events

import (
	"github.com/bwmarrin/discordgo"
	"purplg.com/memberchannels/internal/pkg/widget"
)

func (config *Config) GuildCreated(session *discordgo.Session, event *discordgo.GuildCreate) {
	log := config.Log.WithField("GuildID", event.Guild.ID)
	log.Debugln("Guild created event")

	guildDB := config.DB.AsGuild(event.Guild.ID)

	w := widget.New(session, log, guildDB, config.Vars.DefaultCategoryName, config.Vars.DefaultChannelName)

	if err := w.Show(); err != nil {
		log.WithError(err).Errorln("Error creating widget")
		return
	}
	config.Widgets[guildDB.GuildID] = w
}
