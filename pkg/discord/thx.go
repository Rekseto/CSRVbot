package discord

import (
	"csrvbot/internal/repos"
	"database/sql"
	"github.com/bwmarrin/discordgo"
	"log"
)

func NotifyThxOnThxInfoChannel(s *discordgo.Session, serverRepo repos.ServerRepo, giveawayRepo repos.GiveawayRepo, guildId, channelId, messageId, participantId, confirmerId, state string) {
	embed := ConstructThxNotificationEmbed(guildId, channelId, messageId, participantId, confirmerId, state)

	serverConfig, err := serverRepo.GetServerConfigForGuild(guildId)
	if err != nil {
		log.Println("("+guildId+") "+"NotifyThxOnThxInfoChannel#serverRepo.GetServerConfigForGuild", err)
		return
	}

	if serverConfig.ThxInfoChannel == "" {
		return
	}

	thxNotification, err := giveawayRepo.GetThxNotification(messageId)
	if err == sql.ErrNoRows {
		message, err := s.ChannelMessageSendEmbed(serverConfig.ThxInfoChannel, embed)
		if err != nil {
			log.Println("("+guildId+") "+"notifyThxOnThxInfoChannel#session.ChannelMessageSendEmbed Unable to send thx info!", err)
			return
		}

		err = giveawayRepo.InsertThxNotification(messageId, message.ID)
		if err != nil {
			log.Println("("+guildId+") "+"notifyThxOnThxInfoChannel#DbMap.Insert Unable to insert to database!", err)
			return
		}
		return
	}

	_, err = s.ChannelMessageEditEmbed(serverConfig.ThxInfoChannel, thxNotification.ThxNotificationMessageId, embed)
	if err != nil {
		log.Println("("+guildId+") "+"notifyThxOnThxInfoChannel#session.ChannelMessageEditEmbed Unable to edit embed!", err)
	}
}
