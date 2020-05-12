package events

import (
	"github.com/sirupsen/logrus"

	"purplg.com/memberchannels/internal/pkg/database"
	"purplg.com/memberchannels/internal/pkg/environment"
	"purplg.com/memberchannels/internal/pkg/widget"
)

type Events struct {
	Vars    *environment.Environment
	Log     *logrus.Entry
	DB      *database.DB
	Widgets map[string]*widget.Widget
}

func New(vars *environment.Environment, log *logrus.Entry, db *database.DB) *Events {
	return &Events{
		Vars:    vars,
		Log:     log,
		DB:      db,
		Widgets: make(map[string]*widget.Widget),
	}
}

func (c *Events) Close() {
	for _, widget := range c.Widgets {
		widget.Close()
	}
	c.DB.Close()
}
