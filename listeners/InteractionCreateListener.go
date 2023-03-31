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
	//DocCommand      commands.DocCommand
}

func NewInteractionCreateListener(giveawayCommand commands.GiveawayCommand, thxCommand commands.ThxCommand, thxmeCommand commands.ThxmeCommand, csrvbotCommand commands.CsrvbotCommand) InteractionCreateListener {
	return InteractionCreateListener{
		GiveawayCommand: giveawayCommand,
		ThxCommand:      thxCommand,
		ThxmeCommand:    thxmeCommand,
		CsrvbotCommand:  csrvbotCommand,
	}
}

func (h InteractionCreateListener) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {
	case "giveaway":
		h.GiveawayCommand.Handle(s, i)
		break
	case "thx":
		h.ThxCommand.Handle(s, i)
		break
	case "thxme":
		h.ThxmeCommand.Handle(s, i)
		break
		//case "doc":
		//	handleDocCommand(s, i)
		//	break
	case "csrvbot":
		h.CsrvbotCommand.Handle(s, i)
		break
	}

}
