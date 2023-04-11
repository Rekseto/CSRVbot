package discord

import (
	"context"
	"csrvbot/internal/repos"
	"github.com/bwmarrin/discordgo"
	"log"
)

func GetAllMembers(session *discordgo.Session, guildId string) []*discordgo.Member {
	after := ""
	var allMembers []*discordgo.Member
	for {
		members, err := session.GuildMembers(guildId, after, 1000)
		if err != nil {
			log.Println("("+guildId+") getAllMembers#session.GuildMembers", err)
			return nil
		}
		allMembers = append(allMembers, members...)
		if len(members) != 1000 {
			break
		}
		after = members[999].User.ID
	}
	return allMembers
}

func CheckHelpers(ctx context.Context, session *discordgo.Session, serverRepo repos.ServerRepo, giveawayRepo repos.GiveawayRepo, userRepo repos.UserRepo, guildId string) {
	serverConfig, err := serverRepo.GetServerConfigForGuild(ctx, guildId)
	if err != nil {
		log.Println("("+guildId+") checkHelpers#ServerRepo.GetServerConfigForGuild", err)
		return
	}
	if serverConfig.HelperRoleThxesNeeded <= 0 {
		return
	}
	if serverConfig.HelperRoleId == "" {
		return
	}

	members := GetAllMembers(session, guildId)

	helpers, err := giveawayRepo.GetParticipantsWithThxAmount(ctx, guildId, serverConfig.HelperRoleThxesNeeded)
	if err != nil {
		log.Println("("+guildId+") checkHelpers#GiveawayRepo.GetParticipantsWithThxAmount", err)
		return
	}

	for _, member := range members {
		shouldHaveRole := false
		for _, helper := range helpers {
			isHelperBlacklisted, err := userRepo.IsUserHelperBlacklisted(ctx, member.User.ID, guildId)
			if err != nil {
				log.Println("("+guildId+") checkHelpers#UserRepo.IsUserHelperBlacklisted", err)
				continue
			}
			if isHelperBlacklisted {
				shouldHaveRole = false
				break
			}
			if helper.UserId == member.User.ID {
				shouldHaveRole = true
				break
			}
		}
		hasRole := HasRoleById(member, serverConfig.HelperRoleId)
		if shouldHaveRole {
			if hasRole {
				continue
			}
			log.Println("Adding helper role to " + member.User.Username + " (" + member.User.ID + ")")
			err = session.GuildMemberRoleAdd(guildId, member.User.ID, serverConfig.HelperRoleId)
			if err != nil {
				log.Println("("+guildId+") checkHelpers#session.GuildMemberRoleAdd", err)
			}
		} else {
			if !hasRole {
				continue
			}
			log.Println("Removing helper role from " + member.User.Username + " (" + member.User.ID + ")")
			err = session.GuildMemberRoleRemove(guildId, member.User.ID, serverConfig.HelperRoleId)
			if err != nil {
				log.Println("("+guildId+") checkHelpers#session.GuildMemberRoleRemove", err)
			}
		}
	}
}

func CheckHelper(ctx context.Context, session *discordgo.Session, serverRepo repos.ServerRepo, giveawayRepo repos.GiveawayRepo, userRepo repos.UserRepo, guildId, memberId string) {
	serverConfig, err := serverRepo.GetServerConfigForGuild(ctx, guildId)
	if err != nil {
		log.Println("("+guildId+") checkHelper#ServerRepo.GetServerConfigForGuild", err)
		return
	}
	if serverConfig.HelperRoleThxesNeeded <= 0 {
		return
	}
	if serverConfig.HelperRoleId == "" {
		return
	}

	member, err := session.GuildMember(guildId, memberId)
	if err != nil {
		log.Println("("+guildId+") checkHelper#session.GuildMember", err)
		return
	}

	hasHelperAmount, err := giveawayRepo.HasThxAmount(ctx, guildId, memberId, serverConfig.HelperRoleThxesNeeded)
	if err != nil {
		log.Println("("+guildId+") checkHelper#GiveawayRepo.HasThxAmount", err)
		return
	}
	isHelperBlacklisted, err := userRepo.IsUserHelperBlacklisted(ctx, memberId, guildId)
	if err != nil {
		log.Println("("+guildId+") checkHelper#UserRepo.IsUserHelperBlacklisted", err)
		return
	}
	hasRole := HasRoleById(member, serverConfig.HelperRoleId)
	if !hasHelperAmount {
		if hasRole {
			log.Println("Removing helper role from " + member.User.Username + " (" + member.User.ID + ")")
			err = session.GuildMemberRoleRemove(guildId, memberId, serverConfig.HelperRoleId)
			if err != nil {
				log.Println("("+guildId+") checkHelper#session.GuildMemberRoleRemove", err)
			}
		}
		return
	}

	if isHelperBlacklisted && hasRole {
		log.Println("Removing helper role from " + member.User.Username + " (" + member.User.ID + ")")
		err = session.GuildMemberRoleRemove(guildId, memberId, serverConfig.HelperRoleId)
		if err != nil {
			log.Println("("+guildId+") checkHelper#session.GuildMemberRoleRemove", err)
		}
		return
	}
	if hasRole {
		return
	}
	log.Println("Adding helper role to " + member.User.Username + " (" + member.User.ID + ")")
	err = session.GuildMemberRoleAdd(guildId, memberId, serverConfig.HelperRoleId)
	if err != nil {
		log.Println("("+guildId+") checkHelper#session.GuildMemberRoleAdd", err)
	}

}
