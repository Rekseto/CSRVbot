package listeners

import (
	"csrvbot/internal/repos"
	"csrvbot/pkg"
	"github.com/bwmarrin/discordgo"
)

type GuildMemberUpdateListener struct {
	UserRepo repos.UserRepo
}

func NewGuildMemberUpdateListener(userRepo *repos.UserRepo) GuildMemberUpdateListener {
	return GuildMemberUpdateListener{
		UserRepo: *userRepo,
	}
}

func (h GuildMemberUpdateListener) Handle(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {
	ctx := pkg.CreateContext()
	if m.GuildID == "" { //can it be even empty?
		return
	}

	h.UserRepo.UpdateMemberSavedRoles(ctx, m.Roles, m.User.ID, m.GuildID)
}
