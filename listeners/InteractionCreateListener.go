package listeners

import (
	"csrvbot/commands"
	"github.com/bwmarrin/discordgo"
)

type InteractionCreateListener struct {
	GiveawayCommand commands.GiveawayCommand
	ThxCommand      commands.ThxCommand
	ThxmeCommand    commands.ThxmeCommand
	CsrvbotCommand  commands.CsrvbotCommand
	DocCommand      commands.DocCommand
}

func NewInteractionCreateListener(giveawayCommand commands.GiveawayCommand, thxCommand commands.ThxCommand, thxmeCommand commands.ThxmeCommand, csrvbotCommand commands.CsrvbotCommand, docCommand commands.DocCommand) InteractionCreateListener {
	return InteractionCreateListener{
		GiveawayCommand: giveawayCommand,
		ThxCommand:      thxCommand,
		ThxmeCommand:    thxmeCommand,
		CsrvbotCommand:  csrvbotCommand,
		DocCommand:      docCommand,
	}
}

func (h InteractionCreateListener) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		h.handleApplicationCommands(s, i)
	case discordgo.InteractionApplicationCommandAutocomplete:
		h.handleApplicationCommandsAutocomplete(s, i)
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
	}
}

func (h InteractionCreateListener) handleApplicationCommandsAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {
	case "doc":
		h.DocCommand.HandleAutocomplete(s, i)
	}
}
