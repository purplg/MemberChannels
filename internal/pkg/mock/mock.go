package mock

import (
	"math/rand"
	"strings"
)

const (
	Test_Token        = "testtoken"
	Test_GuildID      = "testGuildID"
	Test_GuildName    = "testGuildName"
	Test_ChannelID    = "testChannelID"
	Test_ChannelID2   = "testChannelID2"
	Test_ChannelID3   = "testChannelID3"
	Test_ChannelName  = "testChannelName"
	Test_ChannelName2 = "testChannelName2"
	Test_ChannelName3 = "testChannelName3"
	Test_RoleID       = "testRoleID"
	Test_RoleID2      = "testRoleID2"
	Test_RoleID3      = "testRoleID3"
	Test_RoleName     = "testRoleName"
	Test_RoleName2    = "testRoleName2"
	Test_RoleName3    = "testRoleName3"
	Test_UserID       = "testUserID"
	Test_UserID2      = "testUserID2"
	Test_UserID3      = "testUserID3"
	Test_UserName     = "testUserName"
	Test_UserName2    = "testUserName2"
	Test_UserName3    = "testUserName3"
	Test_DBFile       = "memberchannels_testdb"
)

func RandomString() string {
	output := strings.Builder{}
	charSet := "1234567890abcdedfghijklmnopqrstABCDEFGHIJKLMNOP"
	length := 20
	for i := 0; i < length; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}
	return output.String()
}
