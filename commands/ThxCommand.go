package commands

import (
	"csrvbot/internal/repos"
	"csrvbot/pkg"
	"github.com/bwmarrin/discordgo"
	"log"
)

type ThxCommand struct {
	Name          string
	Description   string
	DMPermission  bool
	GiveawayHours string
	GiveawayRepo  repos.GiveawayRepo
	UserRepo      repos.UserRepo
	ServerRepo    repos.ServerRepo
}

func NewThxCommand(giveawayRepo *repos.GiveawayRepo, userRepo *repos.UserRepo, serverRepo *repos.ServerRepo, giveawayHours string) ThxCommand {
	return ThxCommand{
		Name:          "thx",
		Description:   "Podziękowanie innemu użytkownikowi",
		DMPermission:  false,
		GiveawayRepo:  *giveawayRepo,
		UserRepo:      *userRepo,
		ServerRepo:    *serverRepo,
		GiveawayHours: giveawayHours,
	}
}

func (h ThxCommand) Register(s *discordgo.Session) {
	_, err := s.ApplicationCommandCreate(s.State.User.ID, "", &discordgo.ApplicationCommand{
		Name:         h.Name,
		Description:  h.Description,
		DMPermission: &h.DMPermission,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "Użytkownik, któremu chcesz podziękować",
				Required:    true,
			},
		},
	})
	if err != nil {
		log.Println("Could not register command", err)
	}
}

func (h ThxCommand) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guild, err := s.Guild(i.GuildID)
	if err != nil {
		log.Println("("+i.GuildID+") handleThxCommand#session.Guild", err)
		return
	}
	selectedUser := i.ApplicationCommandData().Options[0].UserValue(s)
	author := i.Member.User
	if author.ID == selectedUser.ID {
		pkg.RespondWithMessage(s, i, "Nie można dziękować sobie!")
		return
	}
	if selectedUser.Bot {
		pkg.RespondWithMessage(s, i, "Nie można dziękować botom!")
		return
	}
	if h.UserRepo.IsUserBlacklisted(i.GuildID, selectedUser.ID) {
		pkg.RespondWithMessage(s, i, "Ten użytkownik jest na czarnej liście i nie może brać udziału :(")
		return
	}
	giveaway, err := h.GiveawayRepo.GetGiveawayForGuild(i.GuildID)
	if err != nil || giveaway == nil {
		log.Println("("+i.GuildID+") Could not get giveaway", err)
		return
	}
	participants, err := h.GiveawayRepo.GetParticipantNamesForGiveaway(giveaway.Id)
	if err != nil {
		log.Println("("+i.GuildID+") Could not get participants", err)
		return
	}

	embed := pkg.ConstructThxEmbed(participants, h.GiveawayHours, selectedUser.ID, "", "wait")

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
	if err != nil {
		log.Println("("+i.GuildID+") Could not respond to interaction ("+i.ID+")", err)
		return
	}

	response, err := s.InteractionResponse(i.Interaction)
	if err != nil {
		log.Println("("+i.GuildID+") Could not fetch a response of interaction ("+i.ID+")", err)
		return
	}

	err = h.GiveawayRepo.InsertParticipant(giveaway.Id, i.GuildID, guild.Name, selectedUser.ID, selectedUser.Username, i.ChannelID, response.ID)
	if err != nil {
		log.Println("("+i.GuildID+") Could not insert participant", err)
		str := "Coś poszło nie tak przy dodawaniu podziękowania :("
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &str,
		})
		return
	}
	log.Println("(" + i.GuildID + ") " + author.Username + " has thanked " + selectedUser.Username)
	pkg.NotifyThxOnThxInfoChannel(s, h.ServerRepo, h.GiveawayRepo, i.GuildID, i.ChannelID, response.ID, selectedUser.ID, "", "wait")

	for err = s.MessageReactionAdd(i.ChannelID, response.ID, "✅"); err != nil; err = s.MessageReactionAdd(i.ChannelID, response.ID, "✅") {
	}
	for err = s.MessageReactionAdd(i.ChannelID, response.ID, "⛔"); err != nil; err = s.MessageReactionAdd(i.ChannelID, response.ID, "⛔") {
	}
}
