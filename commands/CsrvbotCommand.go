package commands

import (
	"context"
	"csrvbot/internal/repos"
	"csrvbot/internal/services"
	"csrvbot/pkg"
	"csrvbot/pkg/discord"
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
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
	CsrvClient    services.CsrvClient
}

func NewCsrvbotCommand(serverRepo *repos.ServerRepo, giveawayRepo *repos.GiveawayRepo, userRepo *repos.UserRepo, csrvClient *services.CsrvClient) CsrvbotCommand {
	return CsrvbotCommand{
		Name:         "csrvbot",
		Description:  "Komendy konfiguracyjne i administracyjne",
		DMPermission: false,
		ThxMinValue:  0.0,
		ServerRepo:   *serverRepo,
		GiveawayRepo: *giveawayRepo,
		UserRepo:     *userRepo,
		CsrvClient:   *csrvClient,
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
	// todo handle context commands for blacklists
	_, err = s.ApplicationCommandCreate(s.State.User.ID, "", &discordgo.ApplicationCommand{
		Name:         "blacklist",
		DMPermission: &h.DMPermission,
		Type:         discordgo.MessageApplicationCommand,
	})
	if err != nil {
		log.Println("Could not register context command", err)
	}

	_, err = s.ApplicationCommandCreate(s.State.User.ID, "", &discordgo.ApplicationCommand{
		Name:         "blacklist",
		DMPermission: &h.DMPermission,
		Type:         discordgo.UserApplicationCommand,
	})
	if err != nil {
		log.Println("Could not register context command", err)
	}
}

func (h CsrvbotCommand) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := pkg.CreateContext()
	member, err := s.GuildMember(i.GuildID, i.Member.User.ID)
	if err != nil {
		log.Println("("+i.GuildID+") "+"handleGiveawayReactions#s.GuildMember", err)
		return
	}
	if !discord.HasAdminPermissions(ctx, s, h.ServerRepo, member, i.GuildID) {
		discord.RespondWithMessage(s, i, "Nie masz uprawnień do tej komendy")
		return
	}

	switch i.ApplicationCommandData().Options[0].Name {
	case "settings":
		h.handleSettings(ctx, s, i)
	case "delete":
		h.handleDelete(ctx, s, i)
	case "start":
		h.handleStart(ctx, s, i)
	case "blacklist":
		h.handleBlacklist(ctx, s, i)
	case "unblacklist":
		h.handleUnblacklist(ctx, s, i)
	case "helperblacklist":
		h.handleHelperBlacklist(ctx, s, i)
	case "helperunblacklist":
		h.handleHelperUnblacklist(ctx, s, i)
	}
}

func (h CsrvbotCommand) handleSettings(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Options[0].Options[0].Name {
	case "giveawaychannel":
		h.handleGiveawayChannelSet(ctx, s, i)
	case "thxinfochannel":
		h.handleThxInfoChannelSet(ctx, s, i)
	case "adminrole":
		h.handleAdminRoleSet(ctx, s, i)
	case "helperrole":
		h.handleHelperRoleSet(ctx, s, i)
	case "helperthxamount":
		h.handleHelperThxAmountSet(ctx, s, i)
	}
}

func (h CsrvbotCommand) handleStart(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	guild, err := s.Guild(i.GuildID)
	if err != nil {
		log.Println("handleStart s.Guild(" + i.GuildID + ")")
		guild, err = s.Guild(i.GuildID)
		if err != nil {
			return
		}
	}

	// below it should be copy
	discord.FinishGiveaway(ctx, s, h.ServerRepo, h.GiveawayRepo, h.CsrvClient, guild.ID)
	discord.RespondWithMessage(s, i, "Podjęto próbę rozstrzygnięcia giveawayu")

	discord.CreateMissingGiveaways(ctx, s, h.ServerRepo, h.GiveawayRepo, guild)
}

func (h CsrvbotCommand) handleDelete(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	giveaway, err := h.GiveawayRepo.GetGiveawayForGuild(ctx, i.GuildID)
	if err != nil {
		log.Println("handleDelete h.GiveawayRepo.GetGiveawayForGuild(" + i.GuildID + ")")
		return
	}
	participants, err := h.GiveawayRepo.GetParticipantsForGiveaway(ctx, giveaway.Id)
	if err != nil {
		log.Println("handleDelete h.GiveawayRepo.GetParticipantsForGiveaway", err)
		return
	}
	selectedUser := i.ApplicationCommandData().Options[0].Options[0].UserValue(s)
	discord.RespondWithMessage(s, i, "Podjęto próbę usunięcia użytkownika z giveawayu")
	err = h.GiveawayRepo.RemoveParticipants(ctx, giveaway.Id, selectedUser.ID)
	if err != nil {
		log.Println("handleDelete h.GiveawayRepo.RemoveParticipants", err)
		return
	}
	participantNames, err := h.GiveawayRepo.GetParticipantNamesForGiveaway(ctx, giveaway.Id)
	if err != nil {
		log.Println("handleDelete h.GiveawayRepo.GetParticipantNamesForGiveaway", err)
		return
	}
	for _, participant := range participants {
		if participant.UserId != selectedUser.ID {
			return
		}
		embed := discord.ConstructThxEmbed(participantNames, h.GiveawayHours, participant.UserId, "", "reject")
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

func (h CsrvbotCommand) handleBlacklist(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	selectedUser := i.ApplicationCommandData().Options[0].Options[0].UserValue(s)
	if selectedUser.Bot {
		discord.RespondWithMessage(s, i, "Nie możesz dodać bota do blacklisty")
		return
	}
	isUserBlacklisted, err := h.UserRepo.IsUserBlacklisted(ctx, selectedUser.ID, i.GuildID)
	if err != nil {
		log.Println("handleBlacklist h.UserRepo.IsUserBlacklisted", err)
		return
	}
	if isUserBlacklisted {
		discord.RespondWithMessage(s, i, "Użytkownik jest już na blackliście")
		return
	}

	err = h.UserRepo.AddBlacklistForUser(ctx, selectedUser.ID, i.GuildID, i.Member.User.ID)
	if err != nil {
		log.Println("handleBlacklist h.UserRepo.AddBlacklistForUser", err)
		discord.RespondWithMessage(s, i, "Nie udało się dodać użytkownika do blacklisty")
		return
	}
	log.Println("(" + i.GuildID + ") " + i.Member.User.Username + " blacklisted " + selectedUser.Username)
	discord.RespondWithMessage(s, i, "Dodano użytkownika do blacklisty")
}

func (h CsrvbotCommand) handleUnblacklist(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	selectedUser := i.ApplicationCommandData().Options[0].Options[0].UserValue(s)

	isUserBlacklisted, err := h.UserRepo.IsUserBlacklisted(ctx, selectedUser.ID, i.GuildID)
	if err != nil {
		log.Println("handleUnblacklist h.UserRepo.IsUserBlacklisted", err)
		return
	}
	if !isUserBlacklisted {
		discord.RespondWithMessage(s, i, "Użytkownik nie jest na blackliście")
		return
	}
	err = h.UserRepo.RemoveBlacklistForUser(ctx, selectedUser.ID, i.GuildID)
	if err != nil {
		log.Println("handleUnblacklist h.UserRepo.RemoveBlacklistForUser", err)
		discord.RespondWithMessage(s, i, "Nie udało się usunąć użytkownika z blacklisty")
		return
	}
	log.Println("(" + i.GuildID + ") " + i.Member.User.Username + " unblacklisted " + selectedUser.Username)
	discord.RespondWithMessage(s, i, "Usunięto użytkownika z blacklisty")
}

func (h CsrvbotCommand) handleHelperBlacklist(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	selectedUser := i.ApplicationCommandData().Options[0].Options[0].UserValue(s)
	if selectedUser.Bot {
		discord.RespondWithMessage(s, i, "Nie możesz dodać bota do helper-blacklisty")
		return
	}
	isUserHelperBlacklisted, err := h.UserRepo.IsUserHelperBlacklisted(ctx, selectedUser.ID, i.GuildID)
	if err != nil {
		log.Println("handleHelperBlacklist h.UserRepo.IsUserHelperBlacklisted", err)
		return
	}
	if isUserHelperBlacklisted {
		discord.RespondWithMessage(s, i, "Użytkownik jest już na helper-blackliście")
		return
	}

	err = h.UserRepo.AddHelperBlacklistForUser(ctx, selectedUser.ID, i.GuildID, i.Member.User.ID)
	if err != nil {
		log.Println("handleHelperBlacklist h.UserRepo.AddHelperBlacklistForUser", err)
		discord.RespondWithMessage(s, i, "Nie udało się dodać użytkownika do helper-blacklisty")
		return
	}
	log.Println("(" + i.GuildID + ") " + i.Member.User.Username + " helper-blacklisted " + selectedUser.Username)
	discord.RespondWithMessage(s, i, "Użytkownik został zablokowany z możliwości zostania pomocnym")
}

func (h CsrvbotCommand) handleHelperUnblacklist(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	selectedUser := i.ApplicationCommandData().Options[0].Options[0].UserValue(s)

	isUserHelperBlacklisted, err := h.UserRepo.IsUserHelperBlacklisted(ctx, selectedUser.ID, i.GuildID)
	if err != nil {
		log.Println("handleHelperUnblacklist h.UserRepo.IsUserHelperBlacklisted", err)
		return
	}
	if !isUserHelperBlacklisted {
		discord.RespondWithMessage(s, i, "Użytkownik nie jest na helper-blackliście")
		return
	}
	err = h.UserRepo.RemoveHelperBlacklistForUser(ctx, selectedUser.ID, i.GuildID)
	if err != nil {
		log.Println("handleHelperUnblacklist h.UserRepo.RemoveHelperBlacklistForUser", err)
		discord.RespondWithMessage(s, i, "Nie udało się usunąć użytkownika z helper-blacklisty")
		return
	}
	log.Println("(" + i.GuildID + ") " + i.Member.User.Username + " helper-unblacklisted " + selectedUser.Username)
	discord.RespondWithMessage(s, i, "Użytkownik został usunięty z helper-blacklisty")
}

func (h CsrvbotCommand) handleGiveawayChannelSet(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	channelId := i.ApplicationCommandData().Options[0].Options[0].Options[0].ChannelValue(s).ID
	channel, err := s.Channel(channelId)
	if err != nil {
		log.Println("handleGiveawayChannelSet s.Channel", err)
		discord.RespondWithMessage(s, i, "Nie udało się ustawić kanału")
		return
	}
	err = h.ServerRepo.SetMainChannelForGuild(ctx, i.GuildID, channelId)
	if err != nil {
		log.Println("handleGiveawayChannelSet h.ServerRepo.SetMainChannelForGuild", err)
		discord.RespondWithMessage(s, i, "Nie udało się ustawić kanału")
		return
	}
	log.Println("(" + i.GuildID + ") " + i.Member.User.Username + " set giveaway channel to " + channel.Name + " (" + channel.ID + ")")
	discord.RespondWithMessage(s, i, "Ustawiono kanał do ogłoszeń giveawaya na "+channel.Mention())
	guild, err := s.Guild(i.GuildID)
	if err != nil {
		log.Println("handleGiveawayChannelSet s.Guild", err)
		return
	}
	discord.CreateMissingGiveaways(ctx, s, h.ServerRepo, h.GiveawayRepo, guild)
}

func (h CsrvbotCommand) handleThxInfoChannelSet(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	channelId := i.ApplicationCommandData().Options[0].Options[0].Options[0].ChannelValue(s).ID
	channel, err := s.Channel(channelId)
	if err != nil {
		log.Println("handleThxInfoChannelSet s.Channel", err)
		discord.RespondWithMessage(s, i, "Nie udało się ustawić kanału")
		return
	}
	err = h.ServerRepo.SetThxInfoChannelForGuild(ctx, i.GuildID, channelId)
	if err != nil {
		log.Println("handleThxInfoChannelSet h.ServerRepo.SetThxInfoChannelForGuild", err)
		discord.RespondWithMessage(s, i, "Nie udało się ustawić kanału")
		return
	}
	log.Println("(" + i.GuildID + ") " + i.Member.User.Username + " set thx info channel to " + channel.Name + " (" + channel.ID + ")")
	discord.RespondWithMessage(s, i, "Ustawiono kanał do powiadomień o thx na "+channel.Mention())
}

func (h CsrvbotCommand) handleAdminRoleSet(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	roleId := i.ApplicationCommandData().Options[0].Options[0].Options[0].RoleValue(s, i.GuildID).ID
	role, err := s.State.Role(i.GuildID, roleId)
	if err != nil {
		log.Println("handleAdminRoleSet s.State.Role", err)
		discord.RespondWithMessage(s, i, "Nie udało się ustawić roli")
		return
	}
	err = h.ServerRepo.SetAdminRoleIdForGuild(ctx, i.GuildID, role.ID)
	if err != nil {
		log.Println("handleAdminRoleSet h.ServerRepo.SetAdminRoleIdForGuild", err)
		discord.RespondWithMessage(s, i, "Nie udało się ustawić roli")
		return
	}
	log.Println("(" + i.GuildID + ") " + i.Member.User.Username + " set admin role to " + role.Name)
	discord.RespondWithMessage(s, i, "Ustawiono rolę admina na "+role.Name)
}

func (h CsrvbotCommand) handleHelperRoleSet(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	roleId := i.ApplicationCommandData().Options[0].Options[0].Options[0].RoleValue(s, i.GuildID).ID
	role, err := s.State.Role(i.GuildID, roleId)
	if err != nil {
		log.Println("handleHelperRoleSet s.State.Role", err)
		discord.RespondWithMessage(s, i, "Nie udało się ustawić roli")
		return
	}
	err = h.ServerRepo.SetHelperRoleIdForGuild(ctx, i.GuildID, role.ID)
	if err != nil {
		log.Println("handleHelperRoleSet h.ServerRepo.SetHelperRoleIdForGuild", err)
		discord.RespondWithMessage(s, i, "Nie udało się ustawić roli")
		return
	}
	log.Println("(" + i.GuildID + ") " + i.Member.User.Username + " set helper role to " + role.Name)
	discord.RespondWithMessage(s, i, "Ustawiono rolę helpera na "+role.Name)
}

func (h CsrvbotCommand) handleHelperThxAmountSet(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	amount := i.ApplicationCommandData().Options[0].Options[0].Options[0].UintValue()
	err := h.ServerRepo.SetHelperThxesNeededForGuild(ctx, i.GuildID, amount)
	if err != nil {
		log.Println("handleHelperThxAmountSet h.ServerRepo.SetHelperThxAmountForGuild", err)
		discord.RespondWithMessage(s, i, "Nie udało się ustawić ilości thx")
		return
	}
	log.Println("(" + i.GuildID + ") " + i.Member.User.Username + " set helper thx amount to " + strconv.FormatUint(amount, 10))
	discord.RespondWithMessage(s, i, "Ustawiono ilość thx potrzebną do uzyskania rangi helpera na "+strconv.FormatUint(amount, 10))
	discord.CheckHelpers(ctx, s, h.ServerRepo, h.GiveawayRepo, h.UserRepo, i.GuildID)
}
