package commands

import (
	"csrvbot/internal/services"
	"csrvbot/pkg/discord"
	"github.com/bwmarrin/discordgo"
	"log"
)

type DocCommand struct {
	Name         string
	Description  string
	DMPermission bool
	GithubClient services.GithubClient
}

func NewDocCommand(githubClient *services.GithubClient) DocCommand {
	return DocCommand{
		Name:         "doc",
		Description:  "Wysyła link do danego poradnika",
		DMPermission: false,
		GithubClient: *githubClient,
	}
}

func (h DocCommand) Register(s *discordgo.Session) {
	_, err := s.ApplicationCommandCreate(s.State.User.ID, "", &discordgo.ApplicationCommand{
		Name:         h.Name,
		Description:  h.Description,
		DMPermission: &h.DMPermission,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:         discordgo.ApplicationCommandOptionString,
				Name:         "nazwa",
				Description:  "Nazwa poradnika",
				Required:     true,
				Autocomplete: true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "anchor",
				Description: "Nazwa sekcji (nagłówka)",
				Required:    false,
			},
		},
	})
	if err != nil {
		log.Println("Could not register command", err)
	}
}

func (h DocCommand) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	docName := i.ApplicationCommandData().Options[0].StringValue()

	docExists, err := h.GithubClient.GetDocExists(docName)
	if err != nil {
		log.Println("Could not get doc", err)
		discord.RespondWithMessage(s, i, "Wystąpił błąd podczas wyszukiwania poradnika")
		return
	}

	if !docExists {
		discord.RespondWithMessage(s, i, "Taki poradnik nie istnieje")
		return
	}

	anchor := ""
	if len(i.ApplicationCommandData().Options) == 2 {
		anchor = "#" + i.ApplicationCommandData().Options[1].StringValue()
	}
	discord.RespondWithMessage(s, i, "<https://github.com/craftserve/docs/blob/master/"+docName+".md"+anchor+">")
}

func (h DocCommand) HandleAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	var choices []*discordgo.ApplicationCommandOptionChoice

	docs, err := h.GithubClient.GetDocs(data.Options[0].StringValue())
	if err != nil {
		log.Println("Could not get docs", err)
		return
	}

	for _, doc := range docs {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  doc,
			Value: doc,
		})
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
	if err != nil {
		log.Println("Could not respond to interaction", err)
	}
}
