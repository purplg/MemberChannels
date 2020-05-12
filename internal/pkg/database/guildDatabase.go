package database

import (
	"github.com/sirupsen/logrus"
	"github.com/xujiajun/nutsdb"
)

// Field names
const (
	channelPrefix  = "CHAN_"
	C_OWNER        = channelPrefix + "owner"

	userPrefix     = "USER_"
	U_CHANNELNAME  = userPrefix + "channelName"

	guildPrefix    = "GUILD_"
	G_CATEGORYID   = guildPrefix + "categoryID"
	G_CATEGORYNAME = guildPrefix + "categoryName"
	G_CHANNELNAME  = guildPrefix + "channelName"
	G_CHANNELID    = guildPrefix + "channelID"
)

type GuildDB struct {
	*DB

	Log *logrus.Entry

	// Acts as bucket id for database
	GuildID string
}

func (gs *GuildDB) ChannelOwner(channelID string) string {
	value, err := gs.getValue(C_OWNER + channelID)
	if err != nil {
		return ""
	}
	return value
}

func (gs *GuildDB) UserChannelName(userID string) string {
	value, err := gs.getValue(U_CHANNELNAME + userID)
	if err != nil {
		return ""
	}
	return value
}

func (gs *GuildDB) SetUserChannel(userID string, channelID string, channelName string) {
	gs.setValue(U_CHANNELNAME+userID, channelName)
	gs.setValue(C_OWNER+channelID, userID)
}

func (gs *GuildDB) SetCategoryID(value string) {
	gs.setValue(G_CATEGORYID, value)
}

func (gs *GuildDB) CategoryID() string {
	value, err := gs.getValue(G_CATEGORYID)
	if err != nil {
		return ""
	}
	return value
}

func (gs *GuildDB) SetCategoryName(value string) {
	gs.setValue(G_CATEGORYNAME, value)
}

func (gs *GuildDB) CategoryName() string {
	value, err := gs.getValue(G_CATEGORYNAME)
	if err != nil {
		return ""
	}
	return value
}

func (gs *GuildDB) SetChannelID(value string) {
	gs.setValue(G_CHANNELID, value)
}

func (gs *GuildDB) ChannelID() string {
	value, err := gs.getValue(G_CHANNELID)
	if err != nil {
		return ""
	}
	return value
}

func (gs *GuildDB) SetChannelName(value string) {
	gs.setValue(G_CHANNELNAME, value)
}

func (gs *GuildDB) ChannelName() string {
	value, err := gs.getValue(G_CHANNELNAME)
	if err != nil {
		return ""
	}
	return value
}

func (gs *GuildDB) getValue(key string) (string, error) {
	var value string
	if err := gs.View(
		func(tx *nutsdb.Tx) error {
			if entry, err := tx.Get(gs.GuildID, []byte(key)); err != nil {
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

func (gs *GuildDB) setValue(key string, value string) {
	if err := gs.Update(
		func(tx *nutsdb.Tx) error {
			if err := tx.Put(gs.GuildID, []byte(key), []byte(value), 0); err != nil {
				return err
			}
			return nil
		}); err != nil {
		gs.Log.WithError(err).WithFields(logrus.Fields{
			"Key":   key,
			"Value": value,
		}).Errorln("Error setting value in datastore")
	}
}
