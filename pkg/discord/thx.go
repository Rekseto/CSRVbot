package discord

import (
	"context"
	"csrvbot/internal/repos"

	"github.com/bwmarrin/discordgo"
)

func NotifyThxOnThxInfoChannel(
	ctx context.Context,
	s *discordgo.Session,
	serverConfig repos.ServerConfig,
	thxNotificationMessageId *string,

	guildId, channelId,
	answerMessageId, participantId,
	confirmerId,
	state string,
) (messageId string, err error) {
	embed := ConstructThxNotificationEmbed(guildId, channelId, answerMessageId, participantId, confirmerId, state)

	// idk why but sometimes it's empty
	if serverConfig.ThxInfoChannel == "" {
		return
	}

	if thxNotificationMessageId != nil {
		_, err := s.ChannelMessageEditEmbed(serverConfig.ThxInfoChannel, *thxNotificationMessageId, embed)
		if err != nil {
			return "", err
		}
	}

	if thxNotificationMessageId == nil {
		message, err := s.ChannelMessageSendEmbed(serverConfig.ThxInfoChannel, embed)
		if err != nil {
			return "", err
		}

		return message.ID, err
	}

	return "", nil
}

func SendNotifyThxNotification(ctx context.Context) {

}
