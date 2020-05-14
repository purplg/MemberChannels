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

type Widget interface {
	UserJoinedManagedChannel(user *discordgo.User, channel *discordgo.Channel)
	UserDisconnected(user *discordgo.User)
	ManagedChannelChanged(channel *discordgo.Channel)

	IsManagedChannel(channel *discordgo.Channel) bool

	Close()
}

type widget struct {
	session *discordgo.Session
	log     *logrus.Entry
	guildDB database.GuildDatabase

	categoryChannel *discordgo.Channel
	listenChannel   *discordgo.Channel

	currentChannel map[string]*userChannel // map[userID] -> userChannel
	userChannels   map[string]*userChannel // map[channelID] -> userChannel
}

// Only used to initialize a new Widget
type WidgetData struct {
	CategoryID        string
	CategoryName      string
	ListenChannelID   string
	ListenChannelName string
}

type userChannel struct {
	channel   *discordgo.Channel
	owner     *discordgo.User
	userCount uint8
}

// Just initialize values to prepare the widget
func New(session *discordgo.Session, log *logrus.Entry, guildDB database.GuildDatabase, data *WidgetData) (Widget, error) {
	w := &widget{
		session:         session,
		log:             log,
		guildDB:         guildDB,
		categoryChannel: nil,
		listenChannel:   nil,
		currentChannel:  make(map[string]*userChannel),
		userChannels:    make(map[string]*userChannel),
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

func (w *widget) UserJoinedManagedChannel(user *discordgo.User, channel *discordgo.Channel) {
	log := w.log.WithFields(logrus.Fields{
		"user":    user.Username,
		"channel": channel.Name,
	})
	log.Debugln("UserJoinedManagedChannel")

	if _, ok := w.currentChannel[user.ID]; ok {
		w.UserDisconnected(user)
	}

	if w.isListenChannel(channel) {
		if userChan, err := w.newUserChannel(user); err != nil {
			log.WithError(err).Errorln("Error creating User Channel")
		} else {
			w.userChannels[userChan.channel.ID] = userChan
			w.session.GuildMemberMove(w.guildDB.GuildID(), user.ID, userChan.channel.ID)
		}
	} else {
		if userChan, ok := w.userChannels[channel.ID]; ok {
			w.currentChannel[user.ID] = userChan
			userChan.userCount++
		}
	}
}

func (w *widget) UserDisconnected(user *discordgo.User) {
	log := w.log.WithFields(logrus.Fields{
		"user": user.Username,
	})
	log.Debugln("UserDisconnected")

	if userChan, ok := w.currentChannel[user.ID]; ok {
		userChan.userCount--
		log.Debugln("User Channel found")
		if userChan.userCount == 0 {
			log.Debugln("User Channel deleted")
			w.session.ChannelDelete(userChan.channel.ID)
		}
	} else {
		log.Debugln("User Channel doesn't exist")
	}
}

func (w *widget) ManagedChannelChanged(channel *discordgo.Channel) {
	w.log.Warnln("Not implemented yet")
}

func (w *widget) IsManagedChannel(channel *discordgo.Channel) bool {
	return channel.ParentID == w.categoryChannel.ID
}

func (w *widget) Close() {
	w.session.ChannelDelete(w.listenChannel.ID)
}

// Create a new channel for user
func (w *widget) newUserChannel(user *discordgo.User) (*userChannel, error) {
	// Look up the saved channel name for user
	channelName := w.guildDB.UserChannelName(user.ID)

	// Or generate one if none found
	if channelName == "" {
		channelName = fmt.Sprintf("%s's channel", user.Username)
	}

	// Send API request to create the voice channel
	channel, err := w.session.GuildChannelCreateComplex(w.guildDB.GuildID(), discordgo.GuildChannelCreateData{
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
	})
	if err != nil {
		return nil, err
	}

	return &userChannel{
		channel:   channel,
		owner:     user,
		userCount: 0,
	}, nil
}

func (w *widget) getCategory(categoryID, defaultName string) (*discordgo.Channel, error) {
	if category, err := w.session.Channel(categoryID); err == nil {
		return category, nil
	}

	return w.session.GuildChannelCreateComplex(w.guildDB.GuildID(), discordgo.GuildChannelCreateData{
		Name:      defaultName,
		Type:      discordgo.ChannelTypeGuildCategory,
		UserLimit: 1,
	})
}

func (w *widget) getListenChannel(channelID, defaultName, parentID string) (*discordgo.Channel, error) {
	if channel, err := w.session.Channel(channelID); err == nil {
		return channel, nil
	}

	return w.session.GuildChannelCreateComplex(w.guildDB.GuildID(), discordgo.GuildChannelCreateData{
		Name:     defaultName,
		Type:     discordgo.ChannelTypeGuildVoice,
		ParentID: parentID,
	})
}

func (w *widget) isListenChannel(channel *discordgo.Channel) bool {
	return channel.ID == w.listenChannel.ID
}
