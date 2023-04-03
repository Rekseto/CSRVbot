package commands

import (
	"csrvbot/internal/repos"
	"csrvbot/pkg"
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
)

type CsrvbotCommand struct {
	Name          string
	Description   string
	DMPermission  bool
	ThxMinValue   float64
	GiveawayHours string
	ServerRepo    repos.ServerRepo
	GiveawayRepo  repos.GiveawayRepo
	UserRepo      repos.UserRepo
}

func NewCsrvbotCommand(serverRepo *repos.ServerRepo, giveawayRepo *repos.GiveawayRepo, userRepo *repos.UserRepo) CsrvbotCommand {
	return CsrvbotCommand{
		Name:         "csrvbot",
		Description:  "Komendy konfiguracyjne i administracyjne",
		DMPermission: false,
		ThxMinValue:  0.0,
		ServerRepo:   *serverRepo,
		GiveawayRepo: *giveawayRepo,
		UserRepo:     *userRepo,
	}
}

func (h CsrvbotCommand) Register(s *discordgo.Session) {
	_, err := s.ApplicationCommandCreate(s.State.User.ID, "", &discordgo.ApplicationCommand{
		Name:         h.Name,
		Description:  h.Description,
		DMPermission: &h.DMPermission,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "settings",
				Description: "Konfiguracja giveawayów",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "giveawaychannel",
						Description: "Kanał na którym jest prezentowany zwycięzca giveawaya",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type: discordgo.ApplicationCommandOptionChannel,
								ChannelTypes: []discordgo.ChannelType{
									discordgo.ChannelTypeGuildText,
								},
								Name:        "channel",
								Description: "Kanał na którym jest prezentowany zwycięzca giveawaya",
								Required:    true,
							},
						},
					},
					{
						Name:        "thxinfochannel",
						Description: "Kanał na którym są wysyłane wszystkie thxy do rozpatrzenia",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type: discordgo.ApplicationCommandOptionChannel,
								ChannelTypes: []discordgo.ChannelType{
									discordgo.ChannelTypeGuildText,
								},
								Name:        "channel",
								Description: "Kanał na którym są wysyłane wszystkie thxy do rozpatrzenia",
								Required:    true,
							},
						},
					},
					{
						Name:        "adminrole",
						Description: "Rola, która ma dostęp do akceptowania thx i komend administracyjnych",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionRole,
								Name:        "role",
								Description: "Rola, która ma dostęp do akceptowania thx i komend administracyjnych",
								Required:    true,
							},
						},
					},
					{
						Name:        "helperrole",
						Description: "Rola którą dostanie użytkownik, gdy osiągnie daną ilość thx",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionRole,
								Name:        "role",
								Description: "Rola, która ma dostęp do akceptowania thx i komend administracyjnych",
								Required:    true,
							},
						},
					},
					{
						Name:        "helperthxamount",
						Description: "Ilość wymaganych thx do uzyskania roli helpera",
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionInteger,
								Name:        "amount",
								Description: "Ilość wymaganych thx do uzyskania roli helpera",
								Required:    true,
								MinValue:    &h.ThxMinValue,
							},
						},
					},
				},
				Type: discordgo.ApplicationCommandOptionSubCommandGroup,
			},
			{
				Name:        "delete",
				Description: "Usuwa użytkownika z obecnego giveawaya",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionUser,
						Name:        "user",
						Description: "Użytkownik, który ma zostać usunięty",
						Required:    true,
					},
				},
			},
			{
				Name:        "start",
				Description: "Rozstrzyga obecny giveaway",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
			},
			{
				Name:        "blacklist",
				Description: "Dodaje użytkownika do blacklisty możliwości udziału w giveawayu",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionUser,
						Name:        "user",
						Description: "Użytkownik, który ma zostać dodany",
						Required:    true,
					},
				},
			},
			{
				Name:        "unblacklist",
				Description: "Usuwa użytkownika z blacklisty możliwości udziału w giveawayu",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionUser,
						Name:        "user",
						Description: "Użytkownik, który ma zostać usunięty",
						Required:    true,
					},
				},
			},
			{
				Name:        "helperblacklist",
				Description: "Dodaje użytkownika do blacklisty możliwości posiadania rangi helpera",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionUser,
						Name:        "user",
						Description: "Użytkownik, który ma zostać dodany",
						Required:    true,
					},
				},
			},
			{
				Name:        "helperunblacklist",
				Description: "Usuwa użytkownika z blacklisty możliwości posiadania rangi helpera",
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionUser,
						Name:        "user",
						Description: "Użytkownik, który ma zostać usunięty",
						Required:    true,
					},
				},
			},
		},
	})
	if err != nil {
		log.Println("Could not register command", err)
	}
}

func (h CsrvbotCommand) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	member, err := s.GuildMember(i.GuildID, i.Member.User.ID)
	if err != nil {
		log.Println("("+i.GuildID+") "+"handleGiveawayReactions#s.GuildMember", err)
		return
	}
	if !pkg.HasAdminPermissions(s, h.ServerRepo, member, i.GuildID) {
		pkg.RespondWithMessage(s, i, "Nie masz uprawnień do tej komendy")
		return
	}

	switch i.ApplicationCommandData().Options[0].Name {
	case "settings":
		h.handleSettings(s, i)
	case "delete":
		h.handleDelete(s, i)
	case "start":
		h.handleStart(s, i)
	case "blacklist":
		h.handleBlacklist(s, i)
	case "unblacklist":
		h.handleUnblacklist(s, i)
	case "helperblacklist":
		h.handleHelperBlacklist(s, i)
	case "helperunblacklist":
		h.handleHelperUnblacklist(s, i)
	}
}

func (h CsrvbotCommand) handleSettings(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Options[0].Options[0].Name {
	case "giveawaychannel":
		h.handleGiveawayChannelSet(s, i)
	case "thxinfochannel":
		h.handleThxInfoChannelSet(s, i)
	case "adminrole":
		h.handleAdminRoleSet(s, i)
	case "helperrole":
		h.handleHelperRoleSet(s, i)
	case "helperthxamount":
		h.handleHelperThxAmountSet(s, i)
	}
}

func (h CsrvbotCommand) handleStart(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guild, err := s.Guild(i.GuildID)
	if err != nil {
		log.Println("handleStart s.Guild(" + i.GuildID + ")")
		guild, err = s.Guild(i.GuildID)
		if err != nil {
			return
		}
	}

	pkg.FinishGiveaway(s, h.ServerRepo, h.GiveawayRepo, guild.ID)
	pkg.RespondWithMessage(s, i, "Podjęto próbę rozstrzygnięcia giveawayu")

	pkg.CreateMissingGiveaways(s, h.ServerRepo, h.GiveawayRepo, guild)
}

func (h CsrvbotCommand) handleDelete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	giveaway, err := h.GiveawayRepo.GetGiveawayForGuild(i.GuildID)
	if err != nil {
		log.Println("handleDelete h.GiveawayRepo.GetGiveawayForGuild(" + i.GuildID + ")")
		return
	}
	participants, err := h.GiveawayRepo.GetParticipantsForGiveaway(giveaway.Id)
	if err != nil {
		log.Println("handleDelete h.GiveawayRepo.GetParticipantsForGiveaway", err)
		return
	}
	selectedUser := i.ApplicationCommandData().Options[0].Options[0].UserValue(s)
	pkg.RespondWithMessage(s, i, "Podjęto próbę usunięcia użytkownika z giveawayu")
	err = h.GiveawayRepo.RemoveParticipants(giveaway.Id, selectedUser.ID)
	if err != nil {
		log.Println("handleDelete h.GiveawayRepo.RemoveParticipants", err)
		return
	}
	participantNames, err := h.GiveawayRepo.GetParticipantNamesForGiveaway(giveaway.Id)
	if err != nil {
		log.Println("handleDelete h.GiveawayRepo.GetParticipantNamesForGiveaway", err)
		return
	}
	for _, participant := range participants {
		if participant.UserId != selectedUser.ID {
			return
		}
		embed := pkg.ConstructThxEmbed(participantNames, h.GiveawayHours, participant.UserId, "", "reject")
		_, err = s.ChannelMessageEditEmbed(participant.ChannelId, participant.MessageId, embed)
		if err != nil {
			log.Println("("+i.GuildID+") Could not update message", err)
			return
		}
		acceptUserId := participant.AcceptUserId
		if !acceptUserId.Valid {
			return
		}
		err = s.MessageReactionRemove(participant.ChannelId, participant.MessageId, "✅", acceptUserId.String)
		if err != nil {
			log.Println("("+i.GuildID+") Could not remove reaction", err)
			return
		}

	}
}

func (h CsrvbotCommand) handleBlacklist(s *discordgo.Session, i *discordgo.InteractionCreate) {
	selectedUser := i.ApplicationCommandData().Options[0].Options[0].UserValue(s)
	if selectedUser.Bot {
		pkg.RespondWithMessage(s, i, "Nie możesz dodać bota do blacklisty")
		return
	}
	isUserBlacklisted, err := h.UserRepo.IsUserBlacklisted(selectedUser.ID, i.GuildID)
	if err != nil {
		log.Println("handleBlacklist h.UserRepo.IsUserBlacklisted", err)
		return
	}
	if isUserBlacklisted {
		pkg.RespondWithMessage(s, i, "Użytkownik jest już na blackliście")
		return
	}

	err = h.UserRepo.AddBlacklistForUser(selectedUser.ID, i.GuildID, i.Member.User.ID)
	if err != nil {
		log.Println("handleBlacklist h.UserRepo.AddBlacklistForUser", err)
		pkg.RespondWithMessage(s, i, "Nie udało się dodać użytkownika do blacklisty")
		return
	}
	log.Println("(" + i.GuildID + ") " + i.Member.User.Username + " blacklisted " + selectedUser.Username)
	pkg.RespondWithMessage(s, i, "Dodano użytkownika do blacklisty")
}

func (h CsrvbotCommand) handleUnblacklist(s *discordgo.Session, i *discordgo.InteractionCreate) {
	selectedUser := i.ApplicationCommandData().Options[0].Options[0].UserValue(s)

	isUserBlacklisted, err := h.UserRepo.IsUserBlacklisted(selectedUser.ID, i.GuildID)
	if err != nil {
		log.Println("handleUnblacklist h.UserRepo.IsUserBlacklisted", err)
		return
	}
	if !isUserBlacklisted {
		pkg.RespondWithMessage(s, i, "Użytkownik nie jest na blackliście")
		return
	}
	err = h.UserRepo.RemoveBlacklistForUser(selectedUser.ID, i.GuildID)
	if err != nil {
		log.Println("handleUnblacklist h.UserRepo.RemoveBlacklistForUser", err)
		pkg.RespondWithMessage(s, i, "Nie udało się usunąć użytkownika z blacklisty")
		return
	}
	log.Println("(" + i.GuildID + ") " + i.Member.User.Username + " unblacklisted " + selectedUser.Username)
	pkg.RespondWithMessage(s, i, "Usunięto użytkownika z blacklisty")
}

func (h CsrvbotCommand) handleHelperBlacklist(s *discordgo.Session, i *discordgo.InteractionCreate) {
	selectedUser := i.ApplicationCommandData().Options[0].Options[0].UserValue(s)
	if selectedUser.Bot {
		pkg.RespondWithMessage(s, i, "Nie możesz dodać bota do helper-blacklisty")
		return
	}
	isUserHelperBlacklisted, err := h.UserRepo.IsUserHelperBlacklisted(selectedUser.ID, i.GuildID)
	if err != nil {
		log.Println("handleHelperBlacklist h.UserRepo.IsUserHelperBlacklisted", err)
		return
	}
	if isUserHelperBlacklisted {
		pkg.RespondWithMessage(s, i, "Użytkownik jest już na helper-blackliście")
		return
	}

	err = h.UserRepo.AddHelperBlacklistForUser(selectedUser.ID, i.GuildID, i.Member.User.ID)
	if err != nil {
		log.Println("handleHelperBlacklist h.UserRepo.AddHelperBlacklistForUser", err)
		pkg.RespondWithMessage(s, i, "Nie udało się dodać użytkownika do helper-blacklisty")
		return
	}
	log.Println("(" + i.GuildID + ") " + i.Member.User.Username + " helper-blacklisted " + selectedUser.Username)
	pkg.RespondWithMessage(s, i, "Użytkownik został zablokowany z możliwości zostania pomocnym")
}

func (h CsrvbotCommand) handleHelperUnblacklist(s *discordgo.Session, i *discordgo.InteractionCreate) {
	selectedUser := i.ApplicationCommandData().Options[0].Options[0].UserValue(s)

	isUserHelperBlacklisted, err := h.UserRepo.IsUserHelperBlacklisted(selectedUser.ID, i.GuildID)
	if err != nil {
		log.Println("handleHelperUnblacklist h.UserRepo.IsUserHelperBlacklisted", err)
		return
	}
	if !isUserHelperBlacklisted {
		pkg.RespondWithMessage(s, i, "Użytkownik nie jest na helper-blackliście")
		return
	}
	err = h.UserRepo.RemoveHelperBlacklistForUser(selectedUser.ID, i.GuildID)
	if err != nil {
		log.Println("handleHelperUnblacklist h.UserRepo.RemoveHelperBlacklistForUser", err)
		pkg.RespondWithMessage(s, i, "Nie udało się usunąć użytkownika z helper-blacklisty")
		return
	}
	log.Println("(" + i.GuildID + ") " + i.Member.User.Username + " helper-unblacklisted " + selectedUser.Username)
	pkg.RespondWithMessage(s, i, "Użytkownik został usunięty z helper-blacklisty")
}

func (h CsrvbotCommand) handleGiveawayChannelSet(s *discordgo.Session, i *discordgo.InteractionCreate) {
	channelId := i.ApplicationCommandData().Options[0].Options[0].Options[0].ChannelValue(s).ID
	channel, err := s.Channel(channelId)
	if err != nil {
		log.Println("handleGiveawayChannelSet s.Channel", err)
		pkg.RespondWithMessage(s, i, "Nie udało się ustawić kanału")
		return
	}
	err = h.ServerRepo.SetMainChannelForGuild(i.GuildID, channelId)
	if err != nil {
		log.Println("handleGiveawayChannelSet h.ServerRepo.SetMainChannelForGuild", err)
		pkg.RespondWithMessage(s, i, "Nie udało się ustawić kanału")
		return
	}
	log.Println("(" + i.GuildID + ") " + i.Member.User.Username + " set giveaway channel to " + channel.Name + " (" + channel.ID + ")")
	pkg.RespondWithMessage(s, i, "Ustawiono kanał do ogłoszeń giveawaya na "+channel.Mention())
	guild, err := s.Guild(i.GuildID)
	if err != nil {
		log.Println("handleGiveawayChannelSet s.Guild", err)
		return
	}
	pkg.CreateMissingGiveaways(s, h.ServerRepo, h.GiveawayRepo, guild)
}

func (h CsrvbotCommand) handleThxInfoChannelSet(s *discordgo.Session, i *discordgo.InteractionCreate) {
	channelId := i.ApplicationCommandData().Options[0].Options[0].Options[0].ChannelValue(s).ID
	channel, err := s.Channel(channelId)
	if err != nil {
		log.Println("handleThxInfoChannelSet s.Channel", err)
		pkg.RespondWithMessage(s, i, "Nie udało się ustawić kanału")
		return
	}
	err = h.ServerRepo.SetThxInfoChannelForGuild(i.GuildID, channelId)
	if err != nil {
		log.Println("handleThxInfoChannelSet h.ServerRepo.SetThxInfoChannelForGuild", err)
		pkg.RespondWithMessage(s, i, "Nie udało się ustawić kanału")
		return
	}
	log.Println("(" + i.GuildID + ") " + i.Member.User.Username + " set thx info channel to " + channel.Name + " (" + channel.ID + ")")
	pkg.RespondWithMessage(s, i, "Ustawiono kanał do powiadomień o thx na "+channel.Mention())
}

func (h CsrvbotCommand) handleAdminRoleSet(s *discordgo.Session, i *discordgo.InteractionCreate) {
	roleId := i.ApplicationCommandData().Options[0].Options[0].Options[0].RoleValue(s, i.GuildID).ID
	role, err := s.State.Role(i.GuildID, roleId)
	if err != nil {
		log.Println("handleAdminRoleSet s.State.Role", err)
		pkg.RespondWithMessage(s, i, "Nie udało się ustawić roli")
		return
	}
	err = h.ServerRepo.SetAdminRoleForGuild(i.GuildID, role.Name)
	if err != nil {
		log.Println("handleAdminRoleSet h.ServerRepo.SetAdminRoleForGuild", err)
		pkg.RespondWithMessage(s, i, "Nie udało się ustawić roli")
		return
	}
	log.Println("(" + i.GuildID + ") " + i.Member.User.Username + " set admin role to " + role.Name)
	pkg.RespondWithMessage(s, i, "Ustawiono rolę admina na "+role.Name)
}

func (h CsrvbotCommand) handleHelperRoleSet(s *discordgo.Session, i *discordgo.InteractionCreate) {
	roleId := i.ApplicationCommandData().Options[0].Options[0].Options[0].RoleValue(s, i.GuildID).ID
	role, err := s.State.Role(i.GuildID, roleId)
	if err != nil {
		log.Println("handleHelperRoleSet s.State.Role", err)
		pkg.RespondWithMessage(s, i, "Nie udało się ustawić roli")
		return
	}
	err = h.ServerRepo.SetHelperRoleForGuild(i.GuildID, role.Name)
	if err != nil {
		log.Println("handleHelperRoleSet h.ServerRepo.SetHelperRoleForGuild", err)
		pkg.RespondWithMessage(s, i, "Nie udało się ustawić roli")
		return
	}
	log.Println("(" + i.GuildID + ") " + i.Member.User.Username + " set helper role to " + role.Name)
	pkg.RespondWithMessage(s, i, "Ustawiono rolę helpera na "+role.Name)
}

func (h CsrvbotCommand) handleHelperThxAmountSet(s *discordgo.Session, i *discordgo.InteractionCreate) {
	amount := i.ApplicationCommandData().Options[0].Options[0].Options[0].UintValue()
	err := h.ServerRepo.SetHelperThxesNeededForGuild(i.GuildID, amount)
	if err != nil {
		log.Println("handleHelperThxAmountSet h.ServerRepo.SetHelperThxAmountForGuild", err)
		pkg.RespondWithMessage(s, i, "Nie udało się ustawić ilości thx")
		return
	}
	log.Println("(" + i.GuildID + ") " + i.Member.User.Username + " set helper thx amount to " + strconv.FormatUint(amount, 10))
	pkg.RespondWithMessage(s, i, "Ustawiono ilość thx potrzebną do uzyskania rangi helpera na "+strconv.FormatUint(amount, 10))
	pkg.CheckHelpers(s, h.ServerRepo, h.GiveawayRepo, h.UserRepo, i.GuildID)
}
