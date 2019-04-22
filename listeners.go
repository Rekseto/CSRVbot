package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func OnMessageReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if !isThxMessage(r.MessageID) {
		return
	}
	if r.UserID == s.State.User.ID {
		return
	}
	member, _ := s.GuildMember(r.GuildID, r.UserID)
	if hasRole(member, config.AdminRole) {
		participant := getParticipantByMessageId(r.MessageID)
		participant.AcceptTime.Time = time.Now()
		participant.AcceptUser.String = member.User.Username
		participant.AcceptUserId.String = r.UserID
		if r.Emoji.Name == "thumbsup" {
			participant.IsAccepted.Bool = true
		} else if r.Emoji.Name == "thumbsdown" {
			participant.IsAccepted.Bool = false
		}
		participant.update()
	}
}

func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore own messages
	if m.Author.ID == s.State.User.ID {
		return
	}
	// ignore other bots messages
	if m.Author.Bot {
		return
	}

	if !strings.HasPrefix(m.Content, "!") {
		return
	}

	// remove prefix
	m.Content = m.Content[1:]

	args := strings.Fields(m.Content)
	if args[0] == "thx" {
		if len(args) != 2 {
			printGiveawayInfo(&m.ChannelID, &m.GuildID)
			return
		}
		match, _ := regexp.Match("<@[!]?[0-9]*>", []byte(args[1]))
		if !match {
			printGiveawayInfo(&m.ChannelID, &m.GuildID)
			return
		}

		fmt.Print(m.Content)
		return
	}
	if args[0] == "giveaway" {
		printGiveawayInfo(&m.ChannelID, &m.GuildID)
		return
	}
	if args[0] == "csrvbot" {
		if len(args) == 2 {
			if args[1] == "info" {
				printServerInfo(&m.ChannelID, &m.GuildID)
				return
			}
			if args[1] == "start" {
				//			forceStart <- m.GuildID
				return
			}
			if args[1] == "delete" {
				member, err := s.GuildMember(m.GuildID, m.Message.Author.ID)
				if err != nil {
					fmt.Println(err)
					return
				}
				if !hasRole(member, config.AdminRole) {
					_, _ = s.ChannelMessageSend(m.ChannelID, "Brak uprawnień.")
					return
				}
				if len(args) == 2 {
					_, err := s.ChannelMessageSend(m.ChannelID, "Musisz podać ID użytkownika!")
					if err != nil {
						fmt.Println(err)
					}
					return
				}
				deleteFromGiveaway(args[2], m.GuildID)
			}
			if args[1] == "blacklist" {
				member, err := s.GuildMember(m.GuildID, m.Message.Author.ID)
				if err != nil {
					fmt.Println(err)
					return
				}
				if !hasRole(member, config.AdminRole) {
					_, _ = s.ChannelMessageSend(m.ChannelID, "Brak uprawnień.")
					return
				}
				if len(args) == 2 {
					_, err := s.ChannelMessageSend(m.ChannelID, "Musisz podać ID użytkownika!")
					if err != nil {
						fmt.Println(err)
					}
					return
				}
				blacklistUser(args[2], m.GuildID)
			}
		}
		_, err := s.ChannelMessageSend(m.ChannelID, "!csrvbot <delete|resend|start|blacklist|info>")
		if err != nil {
			fmt.Println(err)
		}
	}
}
