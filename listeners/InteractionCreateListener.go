package listeners

import (
	"context"
	"csrvbot/commands"
	"csrvbot/internal/repos"
	"csrvbot/pkg"
	"csrvbot/pkg/discord"
	"github.com/bwmarrin/discordgo"
	"log"
)

type InteractionCreateListener struct {
	GiveawayCommand commands.GiveawayCommand
	ThxCommand      commands.ThxCommand
	ThxmeCommand    commands.ThxmeCommand
	CsrvbotCommand  commands.CsrvbotCommand
	DocCommand      commands.DocCommand
	ResendCommand   commands.ResendCommand
	GiveawayRepo    repos.GiveawayRepo
}

func NewInteractionCreateListener(giveawayCommand commands.GiveawayCommand, thxCommand commands.ThxCommand, thxmeCommand commands.ThxmeCommand, csrvbotCommand commands.CsrvbotCommand, docCommand commands.DocCommand, resendCommand commands.ResendCommand, giveawayRepo *repos.GiveawayRepo) InteractionCreateListener {
	return InteractionCreateListener{
		GiveawayCommand: giveawayCommand,
		ThxCommand:      thxCommand,
		ThxmeCommand:    thxmeCommand,
		CsrvbotCommand:  csrvbotCommand,
		DocCommand:      docCommand,
		ResendCommand:   resendCommand,
		GiveawayRepo:    *giveawayRepo,
	}
}

func (h InteractionCreateListener) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := pkg.CreateContext()
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		h.handleApplicationCommands(s, i)
	case discordgo.InteractionApplicationCommandAutocomplete:
		h.handleApplicationCommandsAutocomplete(s, i)
	case discordgo.InteractionMessageComponent:
		h.handleMessageComponents(ctx, s, i)
	}
}

func (h InteractionCreateListener) handleApplicationCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {
	case "giveaway":
		h.GiveawayCommand.Handle(s, i)
	case "thx":
		h.ThxCommand.Handle(s, i)
	case "thxme":
		h.ThxmeCommand.Handle(s, i)
	case "doc":
		h.DocCommand.Handle(s, i)
	case "csrvbot":
		h.CsrvbotCommand.Handle(s, i)
	case "resend":
		h.ResendCommand.Handle(s, i)
	}
}

func (h InteractionCreateListener) handleApplicationCommandsAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {
	case "doc":
		h.DocCommand.HandleAutocomplete(s, i)
	}
}

func (h InteractionCreateListener) handleMessageComponents(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.MessageComponentData().CustomID {
	case "winnercode":
		hasWon, err := h.GiveawayRepo.HasWonGiveawayByMessageId(ctx, i.Message.ID, i.Member.User.ID)
		if err != nil {
			log.Println("("+i.GuildID+") handleMessageComponents#GiveawayRepo.HasWonGiveawayByMessageId", err)
			return
		}
		if !hasWon {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Nie wygrałeś tego giveawayu!",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				log.Println("("+i.GuildID+") handleMessageComponents#session.InteractionRespond", err)
			}
			return
		}

		code, err := h.GiveawayRepo.GetCodeForInfoMessage(ctx, i.Message.ID)
		if err != nil {
			log.Println("("+i.GuildID+") handleMessageComponents#session.InteractionRespond", err)
			return
		}

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:  discordgo.MessageFlagsEphemeral,
				Embeds: []*discordgo.MessageEmbed{discord.ConstructWinnerEmbed(code)},
			},
		})
		if err != nil {
			log.Println("("+i.GuildID+") handleMessageComponents#session.InteractionRespond", err)
			return
		}
	}
}
