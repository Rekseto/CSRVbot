package listeners

import (
	"context"
	"csrvbot/internal/repos"
	"csrvbot/pkg"
	"github.com/bwmarrin/discordgo"
	"log"
)

type GuildMemberAddListener struct {
	UserRepo repos.UserRepo
}

func NewGuildMemberAddListener(userRepo *repos.UserRepo) GuildMemberAddListener {
	return GuildMemberAddListener{
		UserRepo: *userRepo,
	}
}

func (h GuildMemberAddListener) Handle(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	ctx := pkg.CreateContext()
	if m.GuildID == "" { //can it be even empty?
		return
	}
	h.restoreMemberRoles(ctx, s, m.Member, m.GuildID)
}

func (h GuildMemberAddListener) restoreMemberRoles(ctx context.Context, s *discordgo.Session, member *discordgo.Member, guildId string) {
	memberRoles, err := h.UserRepo.GetRolesForMember(ctx, guildId, member.User.ID)
	for _, role := range memberRoles {
		err = s.GuildMemberRoleAdd(guildId, member.User.ID, role.RoleId)
		if err != nil {
			log.Println("("+guildId+") "+"restoreMemberRoles#session.GuildMemberRoleAdd Error while adding role to member", err)
			continue
		}
	}
}
