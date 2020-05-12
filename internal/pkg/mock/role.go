package mock

import "github.com/bwmarrin/discordgo"

func mockRole(id string, name string) *discordgo.Role {
	return &discordgo.Role{
		ID:   id,
		Name: name,
	}
}

func mockRoles() []*discordgo.Role {
	return []*discordgo.Role{
		mockRole(Test_RoleID,  Test_RoleName),
		mockRole(Test_RoleID2, Test_RoleName2),
		mockRole(Test_RoleID3, Test_RoleName3),
	}
}

func mockRoleIDs(roles discordgo.Roles) []string {
	roleIDs := make([]string, len(roles))

	for i, role := range roles {
		roleIDs[i] = role.ID
	}

	return roleIDs
}
