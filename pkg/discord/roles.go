package discord

import (
	"context"
	"csrvbot/internal/repos"
	"github.com/bwmarrin/discordgo"
	"log"
)

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

func HasRoleById(member *discordgo.Member, roleId string) bool {
	for _, role := range member.Roles {
		if role == roleId {
			return true
		}
	}

	return false
}

func HasAdminPermissions(ctx context.Context, session *discordgo.Session, serverRepo repos.ServerRepo, member *discordgo.Member, guildId string) bool {
	if HasPermission(session, member, guildId, 8) {
		return true
	}
	adminRole, err := serverRepo.GetAdminRoleForGuild(ctx, guildId)
	if err != nil {
		log.Println("("+guildId+") "+"HasAdminPermissions#serverRepo.GetAdminRole", err)
		return false
	}
	if HasRoleById(member, adminRole) {
		return true
	}
	return false
}
