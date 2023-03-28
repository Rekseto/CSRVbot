package repos

import (
	"gopkg.in/gorp.v2"
)

type ServerRepo struct {
	mysql *gorp.DbMap
}

func NewServerRepo(mysql *gorp.DbMap) *ServerRepo {
	mysql.AddTableWithName(ServerConfig{}, "ServerConfig").SetKeys(true, "id")

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
	err := repo.mysql.SelectOne(&serverConfig, "SELECT * FROM ServerConfig WHERE guild_id = ?", guildId)
	if err != nil {
		return nil, err
	}
	return &serverConfig, nil
}

func (repo *ServerRepo) InsertServerConfig(guildId string, giveawayChannel string) error {
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
