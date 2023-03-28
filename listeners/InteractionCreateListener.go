package listeners

import (
	"csrvbot/commands"
	"github.com/bwmarrin/discordgo"
)

type InteractionCreateListener struct {
	GiveawayCommand commands.GiveawayCommand
	ThxCommand      commands.ThxCommand
	ThxmeCommand    commands.ThxmeCommand
	//DocCommand      commands.DocCommand
	//CsrvBotCommands commands.CsrvBotCommands
}

func NewInteractionCreateListener(giveawayCommand commands.GiveawayCommand, thxCommand commands.ThxCommand, thxmeCommand commands.ThxmeCommand) InteractionCreateListener {
	return InteractionCreateListener{
		GiveawayCommand: giveawayCommand,
		ThxCommand:      thxCommand,
		ThxmeCommand:    thxmeCommand,
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
		//case "csrvbot":
		//	handleCsrvBotCommands(s, i)
		//	break
	}

}
