package events

import (
	"github.com/bwmarrin/discordgo"
)

func (config *Events) VoiceStateUpdate(session *discordgo.Session, event *discordgo.VoiceStateUpdate) {
	log := config.Log.WithField("GuildID", event.GuildID)

	lastVoiceState := config.voiceStateCache[event.UserID]
	config.voiceStateCache[event.UserID] = event.VoiceState

	widget := config.Widgets[event.GuildID]
	if widget == nil {
		log.Errorln("Could not find widget for guild")
		return
	}

	// Check if user left a previously member channel
	userLeftChannel := lastVoiceState != nil && lastVoiceState.ChannelID != event.ChannelID
	if userLeftChannel && widget.IsMemberChannel(lastVoiceState.ChannelID) {
		widget.UserLeft(event.UserID, lastVoiceState.ChannelID)
	}

	// Join new channel
	if widget.IsListenChannel(event.ChannelID) {
		widget.UserRequestChannel(event.UserID)
	} else if widget.IsMemberChannel(event.ChannelID) {
		widget.UserJoined(event.UserID, event.ChannelID)
	}
}
