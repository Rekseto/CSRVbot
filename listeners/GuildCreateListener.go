package listeners

import (
	"csrvbot/internal/repos"
	"csrvbot/pkg"
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
	pkg.CreateMissingGiveaways(s, h.ServerRepo, h.GiveawayRepo, g.Guild)
	h.updateAllMembersSavedRoles(g.Guild.ID)
	h.checkHelpers(g.Guild.ID)
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
	guildMembers := pkg.GetAllMembers(h.Session, guildId)
	for _, member := range guildMembers {
		h.UserRepo.UpdateMemberSavedRoles(member.Roles, member.User.ID, guildId)
	}
}

// done, but can't be used in many places fixme
func (h GuildCreateListener) checkHelpers(guildId string) {
	serverConfig, err := h.ServerRepo.GetServerConfigForGuild(guildId)
	if err != nil {
		log.Println("("+guildId+") checkHelpers#ServerRepo.GetServerConfigForGuild", err)
		return
	}
	if serverConfig.HelperRoleThxesNeeded <= 0 {
		return
	}
	if serverConfig.HelperRoleName == "" {
		return
	}

	members := pkg.GetAllMembers(h.Session, guildId)

	helpers, err := h.GiveawayRepo.GetParticipantsWithThxAmount(guildId, serverConfig.HelperRoleThxesNeeded)
	if err != nil {
		log.Println("("+guildId+") checkHelpers#GiveawayRepo.GetParticipantsWithThxAmount", err)
		return
	}

	roleId, err := pkg.GetRoleID(h.Session, guildId, serverConfig.HelperRoleName)
	if err != nil {
		log.Println("("+guildId+") checkHelpers#pkg.GetRoleID", err)
		return
	}

	for _, member := range members {
		shouldHaveRole := false
		for _, helper := range helpers {
			if h.UserRepo.IsUserHelperBlacklisted(member.User.ID, guildId) {
				shouldHaveRole = false
				//break? fixme
			}
			if helper.UserId == member.User.ID {
				shouldHaveRole = true
				break
			}
		}
		if shouldHaveRole {
			for _, memberRole := range member.Roles {
				if memberRole == roleId {
					continue
				}
			}
			err = h.Session.GuildMemberRoleAdd(guildId, member.User.ID, roleId)
			if err != nil {
				log.Println("("+guildId+") checkHelpers#session.GuildMemberRoleAdd", err)
			}
		}
	}
}
