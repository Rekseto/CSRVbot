package repos

import (
	"github.com/go-gorp/gorp"
)

type ServerRepo struct {
	mysql *gorp.DbMap
}

func NewServerRepo(mysql *gorp.DbMap) *ServerRepo {
	mysql.AddTableWithName(ServerConfig{}, "server_configs").SetKeys(true, "id")

	return &ServerRepo{mysql: mysql}
}

type ServerConfig struct {
	Id                    int    `db:"id,primarykey,autoincrement"`
	GuildId               string `db:"guild_id,size:255"`
	AdminRole             string `db:"admin_role,size:255"`
	MainChannel           string `db:"main_channel,size:255"`
	ThxInfoChannel        string `db:"thx_info_channel,size:255"`
	HelperRoleName        string `db:"helper_role_name,size:255"`
	HelperRoleThxesNeeded int    `db:"helper_role_thxes_needed"`
}

func (repo *ServerRepo) GetServerConfigForGuild(guildId string) (*ServerConfig, error) {
	var serverConfig ServerConfig
	err := repo.mysql.SelectOne(&serverConfig, "SELECT * FROM server_configs WHERE guild_id = ?", guildId)
	if err != nil {
		return nil, err
	}
	return &serverConfig, nil
}

func (repo *ServerRepo) InsertServerConfig(guildId, giveawayChannel string) error {
	var serverConfig ServerConfig
	serverConfig.GuildId = guildId
	serverConfig.MainChannel = giveawayChannel
	serverConfig.AdminRole = "CraftserveBotAdmin"
	serverConfig.HelperRoleName = ""
	serverConfig.HelperRoleThxesNeeded = 0
	err := repo.mysql.Insert(&serverConfig)
	if err != nil {
		return err
	}
	return nil
}

func (repo *ServerRepo) GetAdminRoleForGuild(guildId string) (string, error) {
	serverConfig, err := repo.GetServerConfigForGuild(guildId)
	if err != nil {
		return "", err
	}
	return serverConfig.AdminRole, nil
}

func (repo *ServerRepo) GetMainChannelForGuild(guildId string) (string, error) {
	str, err := repo.mysql.SelectStr("SELECT main_channel FROM server_configs WHERE guild_id = ?", guildId)
	if err != nil {
		return "", err
	}
	return str, nil
}

func (repo *ServerRepo) SetMainChannelForGuild(guildId, channelId string) error {
	_, err := repo.mysql.Exec("UPDATE server_configs SET main_channel = ? WHERE guild_id = ?", channelId, guildId)
	if err != nil {
		return err
	}
	return nil
}

func (repo *ServerRepo) SetThxInfoChannelForGuild(guildId, channelId string) error {
	_, err := repo.mysql.Exec("UPDATE server_configs SET thx_info_channel = ? WHERE guild_id = ?", channelId, guildId)
	if err != nil {
		return err
	}
	return nil
}

func (repo *ServerRepo) SetAdminRoleForGuild(guildId, roleName string) error {
	_, err := repo.mysql.Exec("UPDATE server_configs SET admin_role = ? WHERE guild_id = ?", roleName, guildId)
	if err != nil {
		return err
	}
	return nil
}

func (repo *ServerRepo) SetHelperRoleForGuild(guildId, roleName string) error {
	_, err := repo.mysql.Exec("UPDATE server_configs SET helper_role_name = ? WHERE guild_id = ?", roleName, guildId)
	if err != nil {
		return err
	}
	return nil
}

func (repo *ServerRepo) SetHelperThxesNeededForGuild(guildId string, thxesNeeded uint64) error {
	_, err := repo.mysql.Exec("UPDATE server_configs SET helper_role_thxes_needed = ? WHERE guild_id = ?", thxesNeeded, guildId)
	if err != nil {
		return err
	}
	return nil
}
