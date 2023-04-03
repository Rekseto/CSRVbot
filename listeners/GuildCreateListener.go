package listeners

import (
	"csrvbot/internal/repos"
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
	log.Println("Registered guild (" + g.Name + "#" + g.ID + ")")

	h.createConfigurationIfNotExists(g.Guild.ID)
	discord.CreateMissingGiveaways(s, h.ServerRepo, h.GiveawayRepo, g.Guild)
	h.updateAllMembersSavedRoles(g.Guild.ID)
	discord.CheckHelpers(s, h.ServerRepo, h.GiveawayRepo, h.UserRepo, g.Guild.ID)
}

func (h GuildCreateListener) createConfigurationIfNotExists(guildID string) {
	var giveawayChannel string
	channels, _ := h.Session.GuildChannels(guildID)
	for _, channel := range channels {
		if channel.Name == "giveaway" {
			giveawayChannel = channel.ID
		}
	}

	_, err := h.ServerRepo.GetServerConfigForGuild(guildID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = h.ServerRepo.InsertServerConfig(guildID, giveawayChannel)
			if err != nil {
				log.Println("("+guildID+") Could not create server config", err)
			}
		} else {
			log.Println("("+guildID+") Could not get server config", err)
		}
	}
}

func (h GuildCreateListener) updateAllMembersSavedRoles(guildId string) {
	guildMembers := discord.GetAllMembers(h.Session, guildId)
	for _, member := range guildMembers {
		h.UserRepo.UpdateMemberSavedRoles(member.Roles, member.User.ID, guildId)
	}
}
