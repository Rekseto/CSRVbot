package pkg

import (
	"csrvbot/internal/repos"
	"errors"
	"github.com/bwmarrin/discordgo"
	"log"
)

func GetRoleID(session *discordgo.Session, guildId string, roleName string) (string, error) {
	guild, err := session.Guild(guildId)
	if err != nil {
		return "", err
	}

	roles := guild.Roles
	for _, role := range roles {
		if role.Name == roleName {
			return role.ID, nil
		}
	}

	return "", errors.New("no " + roleName + " role available")
}

func HasPermission(session *discordgo.Session, member *discordgo.Member, guildId string, permission int64) bool {
	g, err := session.Guild(guildId)
	if err != nil {
		log.Println("("+guildId+") "+"hasPermisson#session.Guild", err)
		return false
	}
	if g.OwnerID == member.User.ID {
		return true
	}
	for _, roleId := range member.Roles {
		role, err := session.State.Role(guildId, roleId)
		if err != nil {
			log.Println("("+guildId+") "+"hasPermisson#session.State.Role", err)
			return false
		}
		if role.Permissions&permission != 0 {
			return true
		}
	}
	return false
}

func HasRole(session *discordgo.Session, member *discordgo.Member, roleName, guildId string) bool {
	adminRole, err := GetRoleID(session, guildId, roleName)
	if err != nil {
		log.Println("("+guildId+") "+"hasRole#getRoleID("+roleName+")", err)
		return false
	}

	for _, role := range member.Roles {
		if role == adminRole {
			return true
		}
	}

	return false
}

func HasRoleById(member *discordgo.Member, roleId string) bool {
	for _, role := range member.Roles {
		if role == roleId {
			return true
		}
	}

	return false
}

func HasAdminPermissions(session *discordgo.Session, serverRepo repos.ServerRepo, member *discordgo.Member, guildId string) bool {
	if HasPermission(session, member, guildId, 8) {
		return true
	}
	adminRole, err := serverRepo.GetAdminRoleForGuild(guildId)
	if err != nil {
		log.Println("("+guildId+") "+"HasAdminPermissions#serverRepo.GetAdminRole", err)
		return false
	}
	if HasRole(session, member, adminRole, guildId) {
		return true
	}
	return false
}
