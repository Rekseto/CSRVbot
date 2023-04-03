package discord

import (
	"csrvbot/internal/repos"
	"csrvbot/pkg"
	"github.com/bwmarrin/discordgo"
	"log"
	"math/rand"
	"time"
)

func FinishGiveaway(s *discordgo.Session, serverRepo repos.ServerRepo, giveawayRepo repos.GiveawayRepo, guildId string) {
	giveaway, err := giveawayRepo.GetGiveawayForGuild(guildId)
	if giveaway == nil || err != nil {
		log.Println("("+guildId+") Could not get giveaway", err)
		return
	}
	_, err = s.Guild(giveaway.GuildId)
	if err != nil {
		log.Println("("+guildId+") finishGiveaway#session.Guild", err)
		return
	}

	giveawayChannelId, err := serverRepo.GetMainChannelForGuild(guildId)
	if err != nil {
		log.Println("("+guildId+") finishGiveaway#serverRepo.GetMainChannelForGuild", err)
		return
	}

	participants, err := giveawayRepo.GetParticipantsForGiveaway(giveaway.Id)
	if err != nil {
		log.Println("("+guildId+") finishGiveaway#DbMap.Select", err)
	}

	if participants == nil || len(participants) == 0 {
		message, err := s.ChannelMessageSend(giveawayChannelId, "Dzisiaj nikt nie wygrywa, ponieważ nikt nie był w loterii.")
		if err != nil {
			log.Println("("+guildId+") finishGiveaway#session.ChannelMessageSend", err)
		}
		err = giveawayRepo.UpdateGiveaway(giveaway, message.ID, "", "", "")
		if err != nil {
			log.Println("("+guildId+") finishGiveaway#DbMap.Update", err)
		}
		log.Println("(" + guildId + ") Giveaway ended without any participants.")
		return
	}

	code, err := pkg.GetCSRVCode()
	if err != nil {
		log.Println("("+guildId+") finishGiveaway#getCSRVCode", err)
		_, err = s.ChannelMessageSend(giveawayChannelId, "Błąd API Craftserve, nie udało się pobrać kodu!")
		if err != nil {
			return
		}
		return
	}
	rand.New(rand.NewSource(time.Now().UnixNano()))
	winner := participants[rand.Intn(len(participants))]

	member, err := s.GuildMember(guildId, winner.UserId)
	if err != nil {
		log.Println("("+guildId+") finishGiveaway#session.GuildMember", err)
		return
	}
	dmEmbed := ConstructWinnerEmbed(code)
	dm, err := s.UserChannelCreate(winner.UserId)
	if err != nil {
		log.Println("("+guildId+") finishGiveaway#session.UserChannelCreate", err)
	}
	_, err = s.ChannelMessageSendEmbed(dm.ID, dmEmbed)

	mainEmbed := ConstructChannelWinnerEmbed(member.User.Username)
	message, err := s.ChannelMessageSendComplex(giveawayChannelId, &discordgo.MessageSend{
		Embed: mainEmbed,
		Components: []discordgo.MessageComponent{
			&discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					&discordgo.Button{
						Label:    "Kliknij tutaj, aby wyświetlić kod",
						Style:    discordgo.SuccessButton,
						CustomID: "winnercode",
						Emoji: discordgo.ComponentEmoji{
							Name: "🎉",
						},
					},
				},
			},
		},
	})

	if err != nil {
		log.Println("("+guildId+") finishGiveaway#session.ChannelMessageSendEmbed", err)
	}

	err = giveawayRepo.UpdateGiveaway(giveaway, message.ID, code, winner.UserId, member.User.Username)
	if err != nil {
		log.Println("("+guildId+") finishGiveaway#DbMap.Update", err)
	}
	log.Println("(" + guildId + ") " + member.User.Username + " has won the giveaway (Code: " + code + ").")

}

func FinishGiveaways(s *discordgo.Session, giveawayRepo repos.GiveawayRepo, serverRepo repos.ServerRepo) {
	giveaways, err := giveawayRepo.GetUnfinishedGiveaways()
	if err != nil {
		log.Println("finishGiveaways#giveawayRepo.GetUnfinishedGiveaways", err)
		return
	}
	for _, giveaway := range giveaways {
		FinishGiveaway(s, serverRepo, giveawayRepo, giveaway.GuildId)
		guild, err := s.Guild(giveaway.GuildId)
		if err == nil {
			CreateMissingGiveaways(s, serverRepo, giveawayRepo, guild)
		}
	}

}

func CreateMissingGiveaways(s *discordgo.Session, serverRepo repos.ServerRepo, giveawayRepo repos.GiveawayRepo, guild *discordgo.Guild) {
	serverConfig, err := serverRepo.GetServerConfigForGuild(guild.ID)
	if err != nil {
		log.Println("("+guild.ID+") createMissingGiveaways#ServerRepo.GetServerConfigForGuild", err)
		return
	}
	giveawayChannelId := serverConfig.MainChannel
	_, err = s.Channel(giveawayChannelId)
	if err != nil {
		log.Println("("+guild.ID+") createMissingGiveaways#Session.Channel", err)
		return
	}
	giveaway, err := giveawayRepo.GetGiveawayForGuild(guild.ID)
	if giveaway == nil && err == nil {
		err := giveawayRepo.InsertGiveaway(guild.ID, guild.Name)
		if err != nil {
			log.Panicln("createMissingGiveaways#DbMap.Insert", err)
			return
		}
	}

}