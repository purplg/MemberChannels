package events

import (
	"github.com/sirupsen/logrus"
	"purplg.com/memberchannels/internal/pkg/database"
	"purplg.com/memberchannels/internal/pkg/variables"
	"purplg.com/memberchannels/internal/pkg/widget"
)

type Config struct {
	Vars    *variables.Variables
	Log     *logrus.Entry
	DB      *database.DB
	Widgets map[string]*widget.Widget
}

func New(vars *variables.Variables, log *logrus.Entry, db *database.DB) *Config {
	return &Config{
		Vars: vars,
		Log:  log,
		DB:   db,
		Widgets: make(map[string]*widget.Widget),
	}
}

func (c *Config) Close() {
	for _, widget := range c.Widgets {
		widget.Close()
	}
	c.DB.Close()
}
