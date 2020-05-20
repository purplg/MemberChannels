package database

import (
	"github.com/sirupsen/logrus"
	"github.com/xujiajun/nutsdb"
)

// Field names
const (
	channelPrefix = "CHAN_"
	C_OWNER       = channelPrefix + "owner"

	userPrefix    = "USER_"
	U_CHANNELNAME = userPrefix + "channelName"

	guildPrefix    = "GUILD_"
	G_CATEGORYID   = guildPrefix + "categoryID"
	G_CATEGORYNAME = guildPrefix + "categoryName"
	G_CHANNELNAME  = guildPrefix + "channelName"
	G_CHANNELID    = guildPrefix + "channelID"
)

type GuildDatabase interface {
	Log() *logrus.Entry
	GuildID() string

	// User data
	ChannelOwner(channelID string) string
	MemberChannelName(userID string) string
	SetMemberChannel(userid, channelID, channelName string)

	// * Guild data
	CategoryID() string
	CategoryName() string
	ChannelID() string
	ChannelName() string
	SetCategoryID(string)
	SetCategoryName(string)
	SetChannelID(string)
	SetChannelName(string)
}

type guildDB struct {
	*DB

	log *logrus.Entry

	// Acts as bucket id for database
	guildID string
}

func (g guildDB) Log() *logrus.Entry {
	return g.log
}

func (g guildDB) GuildID() string {
	return g.guildID
}

// -----------------------------------------------------------------------------
// User data
// -----------------------------------------------------------------------------
func (g guildDB) ChannelOwner(channelID string) string {
	value, err := g.getValue(C_OWNER + channelID)
	if err != nil {
		return ""
	}
	return value
}

func (g guildDB) MemberChannelName(userID string) string {
	value, err := g.getValue(U_CHANNELNAME + userID)
	if err != nil {
		return ""
	}
	return value
}

func (g guildDB) SetMemberChannel(userID string, channelID string, channelName string) {
	g.setValue(U_CHANNELNAME+userID, channelName)
	g.setValue(C_OWNER+channelID, userID)
}

// -----------------------------------------------------------------------------
// Guild data
// -----------------------------------------------------------------------------
func (g guildDB) CategoryID() string {
	value, err := g.getValue(G_CATEGORYID)
	if err != nil {
		return ""
	}
	return value
}

func (g guildDB) CategoryName() string {
	value, err := g.getValue(G_CATEGORYNAME)
	if err != nil {
		return ""
	}
	return value
}

func (g guildDB) ChannelID() string {
	value, err := g.getValue(G_CHANNELID)
	if err != nil {
		return ""
	}
	return value
}

func (g guildDB) ChannelName() string {
	value, err := g.getValue(G_CHANNELNAME)
	if err != nil {
		return ""
	}
	return value
}

func (g guildDB) SetCategoryID(value string) {
	g.setValue(G_CATEGORYID, value)
}

func (g guildDB) SetCategoryName(value string) {
	g.setValue(G_CATEGORYNAME, value)
}

func (g guildDB) SetChannelID(value string) {
	g.setValue(G_CHANNELID, value)
}

func (g guildDB) SetChannelName(value string) {
	g.setValue(G_CHANNELNAME, value)
}

// -----------------------------------------------------------------------------
// Util functions
// -----------------------------------------------------------------------------
func (g guildDB) getValue(key string) (string, error) {
	var value string
	if err := g.View(
		func(tx *nutsdb.Tx) error {
			if entry, err := tx.Get(g.guildID, []byte(key)); err != nil {
				return err
			} else {
				value = string(entry.Value)
				return nil
			}
		}); err != nil {
		return "", err
	}

	return value, nil
}

func (g guildDB) setValue(key string, value string) {
	if err := g.Update(
		func(tx *nutsdb.Tx) error {
			if err := tx.Put(g.guildID, []byte(key), []byte(value), 0); err != nil {
				return err
			}
			return nil
		}); err != nil {
		g.Log().WithError(err).WithFields(logrus.Fields{
			"Key":   key,
			"Value": value,
		}).Errorln("Error setting value in datastore")
	}
}
