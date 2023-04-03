package listeners

import (
	"csrvbot/internal/repos"
	"csrvbot/pkg"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
)

type MessageReactionAddListener struct {
	GiveawayHours string
	UserRepo      repos.UserRepo
	GiveawayRepo  repos.GiveawayRepo
	ServerRepo    repos.ServerRepo
}

func NewMessageReactionAddListener(giveawayHours string, userRepo *repos.UserRepo, giveawayRepo *repos.GiveawayRepo, serverRepo *repos.ServerRepo) MessageReactionAddListener {
	return MessageReactionAddListener{
		GiveawayHours: giveawayHours,
		UserRepo:      *userRepo,
		GiveawayRepo:  *giveawayRepo,
		ServerRepo:    *serverRepo,
	}
}

func (h MessageReactionAddListener) Handle(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == s.State.User.ID {
		return
	}

	if r.ChannelID == r.UserID {
		return
	}

	if r.Emoji.Name != "✅" && r.Emoji.Name != "⛔" {
		err := s.MessageReactionRemove(r.ChannelID, r.MessageID, r.Emoji.Name, r.UserID)
		if err != nil {
			log.Println("("+r.GuildID+") "+"handleGiveawayReactions#s.MessageReactionRemove", err)
		}
		return
	}

	member, err := s.GuildMember(r.GuildID, r.UserID)
	if err != nil {
		log.Println("("+r.GuildID+") "+"handleGiveawayReactions#s.GuildMember", err)
		return
	}

	isThxMessage, err := h.GiveawayRepo.IsThxMessage(r.MessageID)
	if err != nil {
		log.Println("("+r.GuildID+") "+"handleGiveawayReactions#h.GiveawayRepo.IsThxMessage", err)
		return
	}
	isThxmeMessage, err := h.GiveawayRepo.IsThxmeMessage(r.MessageID)
	if err != nil {
		log.Println("("+r.GuildID+") "+"handleGiveawayReactions#h.GiveawayRepo.IsThxmeMessage", err)
		return
	}

	if isThxMessage {
		if !pkg.HasAdminPermissions(s, h.ServerRepo, member, r.GuildID) {
			err = s.MessageReactionRemove(r.ChannelID, r.MessageID, r.Emoji.Name, r.UserID)
			if err != nil {
				log.Println("("+r.GuildID+") "+"handleGiveawayReactions#s.MessageReactionRemove", err)
			}
			return
		}
		// reactionists...
		participant, err := h.GiveawayRepo.GetParticipant(r.MessageID)
		if err != nil {
			log.Println("handleGiveawayReactions#GetParticipant", err)
			return
		}

		switch r.Emoji.Name {
		case "✅":
			err := h.GiveawayRepo.UpdateParticipant(participant, r.UserID, member.User.Username, true)
			if err != nil {
				log.Println("handleGiveawayReactions#UpdateParticipant", err)
				return
			}
			log.Println(member.User.Username + "(" + member.User.ID + ") zaakceptował udział " + participant.UserName + "(" + participant.UserId + ") w giveawayu o ID " + fmt.Sprintf("%d", participant.GiveawayId))

			participants, err := h.GiveawayRepo.GetParticipantNamesForGiveaway(participant.GiveawayId)
			if err != nil {
				log.Println("("+r.GuildID+") Could not get participants", err)
				return
			}

			embed := pkg.ConstructThxEmbed(participants, h.GiveawayHours, participant.UserId, r.UserID, "confirm")

			_, err = s.ChannelMessageEditEmbed(r.ChannelID, r.MessageID, embed)
			if err != nil {
				log.Println("("+r.GuildID+") Could not update message", err)
				return
			}
			pkg.NotifyThxOnThxInfoChannel(s, h.ServerRepo, h.GiveawayRepo, r.GuildID, r.ChannelID, r.MessageID, participant.UserId, r.UserID, "confirm")
			pkg.CheckHelper(s, h.ServerRepo, h.GiveawayRepo, h.UserRepo, r.GuildID, participant.UserId)
			break
		case "⛔":
			err := h.GiveawayRepo.UpdateParticipant(participant, r.UserID, member.User.Username, false)
			if err != nil {
				log.Println("handleGiveawayReactions#UpdateParticipant", err)
				return
			}
			log.Println(member.User.Username + "(" + member.User.ID + ") odrzucił udział " + participant.UserName + "(" + participant.UserId + ") w giveawayu o ID " + fmt.Sprintf("%d", participant.GiveawayId))

			participants, err := h.GiveawayRepo.GetParticipantNamesForGiveaway(participant.GiveawayId)
			if err != nil {
				log.Println("("+r.GuildID+") Could not get participants", err)
				return
			}

			embed := pkg.ConstructThxEmbed(participants, h.GiveawayHours, participant.UserId, r.UserID, "reject")

			_, err = s.ChannelMessageEditEmbed(r.ChannelID, r.MessageID, embed)
			if err != nil {
				log.Println("("+r.GuildID+") Could not update message", err)
				return
			}
			pkg.NotifyThxOnThxInfoChannel(s, h.ServerRepo, h.GiveawayRepo, r.GuildID, r.ChannelID, r.MessageID, participant.UserId, r.UserID, "reject")
			pkg.CheckHelper(s, h.ServerRepo, h.GiveawayRepo, h.UserRepo, r.GuildID, participant.UserId)
			break
		}
	} else if isThxmeMessage {
		candidate, err := h.GiveawayRepo.GetParticipantCandidate(r.MessageID)
		if err != nil {
			log.Println("handleGiveawayReactions#GetParticipant", err)
			return
		}

		if r.UserID != candidate.CandidateApproverId {
			err = s.MessageReactionRemove(r.ChannelID, r.MessageID, r.Emoji.Name, r.UserID)
			if err != nil {
				log.Println("("+r.GuildID+") "+"handleThxme#s.MessageReactionRemove", err)
			}
			return
		}
		// reactionists...
		switch r.Emoji.Name {
		case "✅":
			err := h.GiveawayRepo.UpdateParticipantCandidate(candidate, true)
			if err != nil {
				log.Println("handleGiveawayReactions#UpdateParticipant", err)
				return
			}
			log.Println(candidate.CandidateApproverName + "(" + candidate.CandidateApproverId + ") zaakceptował prosbe o thx uzytkownika " + candidate.CandidateName + "(" + candidate.CandidateId + ")")

			giveaway, err := h.GiveawayRepo.GetGiveawayForGuild(r.GuildID)
			if err != nil || giveaway == nil {
				log.Println("("+r.GuildID+") Could not get giveaway", err)
				return
			}
			participants, err := h.GiveawayRepo.GetParticipantNamesForGiveaway(giveaway.Id)
			if err != nil {
				log.Println("("+r.GuildID+") Could not get participants", err)
				return
			}

			embed := pkg.ConstructThxEmbed(participants, h.GiveawayHours, candidate.CandidateId, "", "wait")

			err = s.MessageReactionRemove(r.ChannelID, r.MessageID, r.Emoji.Name, r.UserID)

			_, err = s.ChannelMessageEdit(r.ChannelID, r.MessageID, "Prośba o podziękowanie zaakceptowana przez: "+member.User.Mention())
			if err != nil {
				log.Println("("+r.GuildID+") Could not update message", err)
				return
			}
			_, err = s.ChannelMessageEditEmbed(r.ChannelID, r.MessageID, embed)
			if err != nil {
				log.Println("("+r.GuildID+") Could not update message", err)
				return
			}

			guild, err := s.Guild(r.GuildID)
			if err != nil {
				log.Println("("+r.GuildID+") handleThxCommand#session.Guild", err)
				return
			}

			err = h.GiveawayRepo.InsertParticipant(giveaway.Id, r.GuildID, guild.Name, candidate.CandidateId, candidate.CandidateName, r.ChannelID, r.MessageID)
			if err != nil {
				log.Println("("+r.GuildID+") Could not insert participant", err)
				_, err = s.ChannelMessageSend(r.ChannelID, "Coś poszło nie tak przy dodawaniu podziękowania :(")
				return
			}
			log.Println("(" + r.GuildID + ") " + member.User.Username + " has thanked " + candidate.CandidateName)
			pkg.NotifyThxOnThxInfoChannel(s, h.ServerRepo, h.GiveawayRepo, r.GuildID, r.ChannelID, r.MessageID, candidate.CandidateId, "", "wait")

			break
		case "⛔":
			err := h.GiveawayRepo.UpdateParticipantCandidate(candidate, false)
			if err != nil {
				log.Println("handleGiveawayReactions#UpdateParticipant", err)
				return
			}
			log.Println(candidate.CandidateApproverName + "(" + candidate.CandidateApproverId + ") odrzucił prosbe o thx uzytkownika " + candidate.CandidateName + "(" + candidate.CandidateId + ")")

			break
		}
	}

}
