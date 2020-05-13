package widget

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"purplg.com/memberchannels/internal/pkg/database"
)

// A Widget consists of a Category channel and a designated Voice channel
// The Voice channel waits for a User to join. When a user joins,
// the Widget will:
// - Create a new voice channel with parameters:
//   - Name the channel appropriately
//   - Permissions so that the creating user has more control
//   - Parented to the Category channel
// - Save the channel name in case it was generated
// - Move player to new channel
type Widget struct {
	session *discordgo.Session
	log     *logrus.Entry
	guildDB database.GuildDatabase

	categoryChannel *discordgo.Channel
	listenChannel   *discordgo.Channel
	userChannels    map[string]*discordgo.Channel // map[userID]
	usersInChannel  map[string][]string
}

// Only used to initialize a new Widget
type WidgetData struct {
	CategoryID        string
	CategoryName      string
	ListenChannelID   string
	ListenChannelName string
}

// Just initialize values to prepare the widget
func New(session *discordgo.Session, log *logrus.Entry, guildDB database.GuildDatabase, data *WidgetData) (*Widget, error) {
	w := &Widget{
		session:         session,
		log:             log,
		guildDB:         guildDB,
		categoryChannel: nil,
		listenChannel:   nil,
		userChannels:    make(map[string]*discordgo.Channel),
	}

	category, err := w.getCategory(data.CategoryID, data.CategoryName)
	if err != nil {
		return nil, err
	}
	w.categoryChannel = category

	listenChannel, err := w.getListenChannel(data.ListenChannelID, data.ListenChannelName, w.categoryChannel.ID)
	if err != nil {
		return nil, err
	}
	w.listenChannel = listenChannel

	w.guildDB.SetCategoryID(w.categoryChannel.ID)
	w.guildDB.SetCategoryName(w.categoryChannel.Name)
	w.guildDB.SetChannelID(w.listenChannel.ID)
	w.guildDB.SetChannelName(w.listenChannel.Name)

	return w, nil
}

// Create a new channel for user
func (w *Widget) NewUserChannel(user *discordgo.User) {
	// Look up the saved channel name for user
	channelName := w.guildDB.UserChannelName(user.ID)

	// Or generate one if none found
	if channelName == "" {
		channelName = fmt.Sprintf("%s's channel", user.Username)
	}

	// Send API request to create the voice channel
	if channel, err := w.session.GuildChannelCreateComplex(w.guildDB.GuildID(), discordgo.GuildChannelCreateData{
		Name: channelName,
		Type: discordgo.ChannelTypeGuildVoice,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				ID:    user.ID,
				Type:  "member",
				Allow: 16,
				Deny:  0,
			},
		},
		ParentID: w.categoryChannel.ID,
	}); err != nil {
		w.log.WithError(err).Errorln("Failed to create user channel")
	} else {
		w.guildDB.SetUserChannel(user.ID, channel.ID, channelName)
		w.session.GuildMemberMove(w.guildDB.GuildID(), user.ID, channel.ID)
	}
}

// A hack to cleanup all empty channels within category
func (w *Widget) Sweep() error {
	var (
		guild    *discordgo.Guild
		channels []*discordgo.Channel
		err      error
	)

	if guild, err = w.session.State.Guild(w.guildDB.GuildID()); err != nil {
		if guild, err = w.session.Guild(w.guildDB.GuildID()); err != nil {
			return err
		}
	}

	channels, err = w.session.GuildChannels(w.guildDB.GuildID())
	if err != nil {
		return err
	}

	categoryID := w.guildDB.CategoryID()
	listeningChannelID := w.guildDB.ChannelID()

	// Loop through all channels in guild
	for _, channel := range channels {
		if channel.ID == listeningChannelID {
			continue
		}
		// If this is a managed channel
		if channel.ParentID == categoryID {
			// Count the number of users in it
			userCount := 0
			for _, vs := range guild.VoiceStates {
				if vs.ChannelID == channel.ID {
					userCount++
				}
			}
			// If it's empty, remove it
			if userCount == 0 {
				w.session.ChannelDelete(channel.ID)
			}
		}
	}
	return nil
}

func (w *Widget) Close() {
	w.Sweep()
	w.session.ChannelDelete(w.listenChannel.ID)
}
