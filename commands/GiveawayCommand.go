package commands

import (
	"csrvbot/internal/repos"
	"csrvbot/pkg/discord"
	"github.com/bwmarrin/discordgo"
	"log"
)

type GiveawayCommand struct {
	Name          string
	Description   string
	DMPermission  bool
	GiveawayHours string
	GiveawayRepo  repos.GiveawayRepo
}

func NewGiveawayCommand(giveawayRepo *repos.GiveawayRepo, giveawayHours string) GiveawayCommand {
	return GiveawayCommand{
		Name:          "giveaway",
		Description:   "Wyświetla zasady giveawaya",
		DMPermission:  false,
		GiveawayRepo:  *giveawayRepo,
		GiveawayHours: giveawayHours,
	}
}

func (h GiveawayCommand) Register(s *discordgo.Session) {
	_, err := s.ApplicationCommandCreate(s.State.User.ID, "", &discordgo.ApplicationCommand{
		Name:         h.Name,
		Description:  h.Description,
		DMPermission: &h.DMPermission,
	})
	if err != nil {
		log.Println("Could not register command", err)
	}
}

func (h GiveawayCommand) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	giveaway, err := h.GiveawayRepo.GetGiveawayForGuild(i.GuildID)
	if err != nil {
		log.Println("("+i.GuildID+") Could not get giveaway", err)
		return
	}
	if giveaway == nil {
		log.Println("(" + i.GuildID + ") Could not get giveaway")
		return
	}
	participants, err := h.GiveawayRepo.GetParticipantNamesForGiveaway(giveaway.Id)
	if err != nil {
		log.Println("("+i.GuildID+") Could not get participants", err)
		return
	}
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				discord.ConstructInfoEmbed(participants, h.GiveawayHours),
			},
		},
	})
	if err != nil {
		log.Println("("+i.GuildID+") Could not respond to interaction ("+i.ID+")", err)
		return
	}
}
