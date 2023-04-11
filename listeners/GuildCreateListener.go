package listeners

import (
	"context"
	"csrvbot/internal/repos"
	"csrvbot/pkg"
	"csrvbot/pkg/discord"
	"database/sql"
	"github.com/bwmarrin/discordgo"
	"log"
)

type GuildCreateListener struct {
	GiveawayRepo repos.GiveawayRepo
	ServerRepo   repos.ServerRepo
	UserRepo     repos.UserRepo
	Session      *discordgo.Session // TODO: remove this, replace with session from Handle method
}

func NewGuildCreateListener(session *discordgo.Session, giveawayRepo *repos.GiveawayRepo, serverRepo *repos.ServerRepo, userRepo *repos.UserRepo) GuildCreateListener {
	return GuildCreateListener{
		GiveawayRepo: *giveawayRepo,
		ServerRepo:   *serverRepo,
		UserRepo:     *userRepo,
		Session:      session,
	}
}

func (h GuildCreateListener) Handle(s *discordgo.Session, g *discordgo.GuildCreate) {
	ctx := pkg.CreateContext()
	log.Println("Registered guild (" + g.Name + "#" + g.ID + ")")

	h.createConfigurationIfNotExists(ctx, g.Guild.ID)
	discord.CreateMissingGiveaways(ctx, s, h.ServerRepo, h.GiveawayRepo, g.Guild)
	h.updateAllMembersSavedRoles(ctx, g.Guild.ID)
	discord.CheckHelpers(ctx, s, h.ServerRepo, h.GiveawayRepo, h.UserRepo, g.Guild.ID)
}

func (h GuildCreateListener) createConfigurationIfNotExists(ctx context.Context, guildID string) {
	var giveawayChannel string
	channels, _ := h.Session.GuildChannels(guildID)
	for _, channel := range channels {
		if channel.Name == "giveaway" {
			giveawayChannel = channel.ID
		}
	}

	_, err := h.ServerRepo.GetServerConfigForGuild(ctx, guildID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = h.ServerRepo.InsertServerConfig(ctx, guildID, giveawayChannel)
			if err != nil {
				log.Println("("+guildID+") Could not create server config", err)
			}
		} else {
			log.Println("("+guildID+") Could not get server config", err)
		}
	}
}

func (h GuildCreateListener) updateAllMembersSavedRoles(ctx context.Context, guildId string) {
	guildMembers := discord.GetAllMembers(h.Session, guildId)
	for _, member := range guildMembers {
		h.UserRepo.UpdateMemberSavedRoles(ctx, member.Roles, member.User.ID, guildId)
	}
}
