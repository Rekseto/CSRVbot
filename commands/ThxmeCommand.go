package commands

import (
	"csrvbot/internal/repos"
	"csrvbot/pkg"
	"github.com/bwmarrin/discordgo"
	"log"
)

type ThxmeCommand struct {
	Name          string
	Description   string
	DMPermission  bool
	GiveawayHours string
	GiveawayRepo  repos.GiveawayRepo
	UserRepo      repos.UserRepo
	ServerRepo    repos.ServerRepo
}

func NewThxmeCommand(giveawayRepo *repos.GiveawayRepo, userRepo *repos.UserRepo, serverRepo *repos.ServerRepo, giveawayHours string) ThxmeCommand {
	return ThxmeCommand{
		Name:          "thxme",
		Description:   "Poproszenie użytkownika o podziękowanie",
		DMPermission:  false,
		GiveawayRepo:  *giveawayRepo,
		UserRepo:      *userRepo,
		ServerRepo:    *serverRepo,
		GiveawayHours: giveawayHours,
	}
}

func (h ThxmeCommand) Register(s *discordgo.Session) {
	_, err := s.ApplicationCommandCreate(s.State.User.ID, "", &discordgo.ApplicationCommand{
		Name:         h.Name,
		Description:  h.Description,
		DMPermission: &h.DMPermission,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "Użytkownik, którego chcesz poprosić o podziękowanie",
				Required:    true,
			},
		},
	})
	if err != nil {
		log.Println("Could not register command", err)
	}
}

func (h ThxmeCommand) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guild, err := s.Guild(i.GuildID)
	if err != nil {
		log.Println("("+i.GuildID+") handleThxmeCommand#session.Guild", err)
		return
	}
	selectedUser := i.ApplicationCommandData().Options[0].UserValue(s)
	author := i.Member.User
	if author.ID == selectedUser.ID {
		pkg.RespondWithMessage(s, i, "Nie można poprosić o podziękowanie samego siebie!")
		return
	}
	if selectedUser.Bot {
		pkg.RespondWithMessage(s, i, "Nie można prosić o podziękowanie bota!")
		return
	}
	isUserBlacklisted, err := h.UserRepo.IsUserBlacklisted(i.GuildID, selectedUser.ID)
	if err != nil {
		log.Println("("+i.GuildID+") handleThxmeCommand#UserRepo.IsUserBlacklisted", err)
		return
	}
	if isUserBlacklisted {
		pkg.RespondWithMessage(s, i, "Nie możesz poprosić o podziękowanie, gdyż jesteś na czarnej liście!")
		return
	}

	pkg.RespondWithMessage(s, i, selectedUser.Mention()+", czy chcesz podziękować użytkownikowi "+author.Mention()+"?")

	response, err := s.InteractionResponse(i.Interaction)
	if err != nil {
		log.Println("("+i.GuildID+") Could not fetch a response of interaction ("+i.ID+")", err)
		return
	}

	err = h.GiveawayRepo.InsertParticipantCandidate(i.GuildID, guild.Name, author.ID, author.Username, selectedUser.ID, selectedUser.Username, i.ChannelID, response.ID)
	if err != nil {
		log.Println("("+i.GuildID+") Could not insert participant candidate", err)
		str := "Coś poszło nie tak przy dodawaniu kandydata do podziękowania :("
		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &str,
		})
		return
	}
	log.Println("(" + i.GuildID + ") " + author.Username + " has requested thx from " + selectedUser.Username)

	for err = s.MessageReactionAdd(i.ChannelID, response.ID, "✅"); err != nil; err = s.MessageReactionAdd(i.ChannelID, response.ID, "✅") {
	}
	for err = s.MessageReactionAdd(i.ChannelID, response.ID, "⛔"); err != nil; err = s.MessageReactionAdd(i.ChannelID, response.ID, "⛔") {
	}
}
