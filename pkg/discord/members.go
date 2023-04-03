package discord

import (
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

func CheckHelpers(session *discordgo.Session, serverRepo repos.ServerRepo, giveawayRepo repos.GiveawayRepo, userRepo repos.UserRepo, guildId string) {
	serverConfig, err := serverRepo.GetServerConfigForGuild(guildId)
	if err != nil {
		log.Println("("+guildId+") checkHelpers#ServerRepo.GetServerConfigForGuild", err)
		return
	}
	if serverConfig.HelperRoleThxesNeeded <= 0 {
		return
	}
	if serverConfig.HelperRoleName == "" {
		return
	}

	members := GetAllMembers(session, guildId)

	helpers, err := giveawayRepo.GetParticipantsWithThxAmount(guildId, serverConfig.HelperRoleThxesNeeded)
	if err != nil {
		log.Println("("+guildId+") checkHelpers#GiveawayRepo.GetParticipantsWithThxAmount", err)
		return
	}

	roleId, err := GetRoleID(session, guildId, serverConfig.HelperRoleName)
	if err != nil {
		log.Println("("+guildId+") checkHelpers#pkg.GetRoleID", err)
		return
	}

	for _, member := range members {
		shouldHaveRole := false
		for _, helper := range helpers {
			isHelperBlacklisted, err := userRepo.IsUserHelperBlacklisted(member.User.ID, guildId)
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
		hasRole := HasRoleById(member, roleId)
		if shouldHaveRole {
			if hasRole {
				continue
			}
			log.Println("Adding helper role to " + member.User.Username + " (" + member.User.ID + ")")
			err = session.GuildMemberRoleAdd(guildId, member.User.ID, roleId)
			if err != nil {
				log.Println("("+guildId+") checkHelpers#session.GuildMemberRoleAdd", err)
			}
		} else {
			if !hasRole {
				continue
			}
			log.Println("Removing helper role from " + member.User.Username + " (" + member.User.ID + ")")
			err = session.GuildMemberRoleRemove(guildId, member.User.ID, roleId)
			if err != nil {
				log.Println("("+guildId+") checkHelpers#session.GuildMemberRoleRemove", err)
			}
		}
	}
}

func CheckHelper(session *discordgo.Session, serverRepo repos.ServerRepo, giveawayRepo repos.GiveawayRepo, userRepo repos.UserRepo, guildId, memberId string) {
	serverConfig, err := serverRepo.GetServerConfigForGuild(guildId)
	if err != nil {
		log.Println("("+guildId+") checkHelper#ServerRepo.GetServerConfigForGuild", err)
		return
	}
	if serverConfig.HelperRoleThxesNeeded <= 0 {
		return
	}
	if serverConfig.HelperRoleName == "" {
		return
	}

	member, err := session.GuildMember(guildId, memberId)
	if err != nil {
		log.Println("("+guildId+") checkHelper#session.GuildMember", err)
		return
	}

	roleId, err := GetRoleID(session, guildId, serverConfig.HelperRoleName)
	if err != nil {
		log.Println("("+guildId+") checkHelper#GetRoleID", err)
		return
	}

	hasHelperAmount, err := giveawayRepo.HasThxAmount(guildId, memberId, serverConfig.HelperRoleThxesNeeded)
	if err != nil {
		log.Println("("+guildId+") checkHelper#GiveawayRepo.HasThxAmount", err)
		return
	}
	isHelperBlacklisted, err := userRepo.IsUserHelperBlacklisted(memberId, guildId)
	if err != nil {
		log.Println("("+guildId+") checkHelper#UserRepo.IsUserHelperBlacklisted", err)
		return
	}
	hasRole := HasRoleById(member, roleId)
	if !hasHelperAmount {
		if hasRole {
			log.Println("Removing helper role from " + member.User.Username + " (" + member.User.ID + ")")
			err = session.GuildMemberRoleRemove(guildId, memberId, roleId)
			if err != nil {
				log.Println("("+guildId+") checkHelper#session.GuildMemberRoleRemove", err)
			}
		}
		return
	}

	if isHelperBlacklisted && hasRole {
		log.Println("Removing helper role from " + member.User.Username + " (" + member.User.ID + ")")
		err = session.GuildMemberRoleRemove(guildId, memberId, roleId)
		if err != nil {
			log.Println("("+guildId+") checkHelper#session.GuildMemberRoleRemove", err)
		}
		return
	}
	if hasRole {
		return
	}
	log.Println("Adding helper role to " + member.User.Username + " (" + member.User.ID + ")")
	err = session.GuildMemberRoleAdd(guildId, memberId, roleId)
	if err != nil {
		log.Println("("+guildId+") checkHelper#session.GuildMemberRoleAdd", err)
	}

}
