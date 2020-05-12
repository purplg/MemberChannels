package database

import (
	"github.com/sirupsen/logrus"
	"github.com/xujiajun/nutsdb"
)

var db *DB

type DB struct {
	*nutsdb.DB
	Log *logrus.Entry
}

func Open(filename string, log *logrus.Entry) (*DB, error) {
	if db == nil {
		dboptions := nutsdb.DefaultOptions
		dboptions.Dir = filename

		if innerDB, err := nutsdb.Open(dboptions); err != nil {
			return nil, err
		} else {
			return &DB{innerDB, log}, nil
		}
	}

	return db, nil
}

func (db *DB) AsGuild(guildID string) GuildDatabase {
	return guildDB{
		DB:      db,
		log:     db.Log.WithField("guildID", guildID),
		guildID: guildID,
	}
}
