package listeners

import (
	"csrvbot/internal/repos"
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
	if m.GuildID == "" { //can it be even empty?
		return
	}
	h.restoreMemberRoles(s, m.Member, m.GuildID)
}

func (h GuildMemberAddListener) restoreMemberRoles(s *discordgo.Session, member *discordgo.Member, guildId string) {
	memberRoles, err := h.UserRepo.GetRolesForMember(guildId, member.User.ID)
	for _, role := range memberRoles {
		err = s.GuildMemberRoleAdd(guildId, member.User.ID, role.RoleId)
		if err != nil {
			log.Println("("+guildId+") "+"restoreMemberRoles#session.GuildMemberRoleAdd Error while adding role to member", err)
			continue
		}
	}
}
