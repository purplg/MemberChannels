package widget

import (
	"github.com/bwmarrin/discordgo"
)

type userChannel struct {
	*discordgo.Channel
	owner    *discordgo.User
	Visitors []string
}

func (uc *userChannel) AddVisitor(userID string) {
	uc.Visitors = append(uc.Visitors, userID)
}

func (uc *userChannel) RemoveVisitor(userID string) {
	for i := 0; i < len(uc.Visitors); i++ {
		if uc.Visitors[i] == userID {
			uc.Visitors = append(uc.Visitors[:i], uc.Visitors[i+1:]...)
			return
		}
	}
}
