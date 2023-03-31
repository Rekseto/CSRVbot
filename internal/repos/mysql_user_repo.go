package repos

import (
	"database/sql"
	"gopkg.in/gorp.v2"
	"log"
)

type UserRepo struct {
	mysql *gorp.DbMap
}

func NewUserRepo(mysql *gorp.DbMap) *UserRepo {
	mysql.AddTableWithName(Blacklist{}, "Blacklists").SetKeys(true, "id").SetUniqueTogether("guild_id", "user_id")
	mysql.AddTableWithName(MemberRole{}, "MemberRoles").SetKeys(true, "id")
	mysql.AddTableWithName(HelperBlacklist{}, "HelperBlacklist").SetKeys(true, "id").SetUniqueTogether("guild_id", "user_id")

	return &UserRepo{mysql: mysql}
}

type Blacklist struct {
	Id            int    `db:"id,primarykey,autoincrement"`
	GuildId       string `db:"guild_id,size:255"`
	UserId        string `db:"user_id,size:255"`
	BlacklisterId string `db:"blacklister_id,size:255"`
}

type MemberRole struct {
	Id       int    `db:"id,primarykey,autoincrement"`
	GuildId  string `db:"guild_id,size:255"`
	MemberId string `db:"member_id,size:255"`
	RoleId   string `db:"role_id,size:255"`
}

type HelperBlacklist struct {
	Id            int    `db:"id,primarykey,autoincrement"`
	GuildId       string `db:"guild_id,size:255"`
	UserId        string `db:"user_id,size:255"`
	BlacklisterId string `db:"blacklister_id,size:255"`
}

func (repo *UserRepo) GetRolesForMember(guildId, memberId string) ([]MemberRole, error) {
	var memberRoles []MemberRole
	_, err := repo.mysql.Select(&memberRoles, "SELECT * FROM MemberRoles WHERE guild_id = ? AND member_id = ?", guildId, memberId)
	if err != nil {
		return nil, err
	}

	return memberRoles, nil
}

func (repo *UserRepo) AddRoleForMember(guildId, memberId, roleId string) error {
	role := MemberRole{GuildId: guildId, RoleId: roleId, MemberId: memberId}
	err := repo.mysql.Insert(&role)
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepo) RemoveRoleForMember(guildId, memberId, roleId string) error {
	_, err := repo.mysql.Exec("DELETE FROM MemberRoles WHERE guild_id = ? AND member_id = ? AND role_id = ?", guildId, memberId, roleId)
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepo) IsUserHelperBlacklisted(userId string, guildId string) bool {
	var helperBlacklist HelperBlacklist
	err := repo.mysql.SelectOne(&helperBlacklist, "SELECT * FROM HelperBlacklist WHERE user_id = ? AND guild_id = ?", userId, guildId)

	if err == sql.ErrNoRows {
		return false
	}

	if err != nil {
		log.Panicln("("+userId+") checkUserBlacklist#DbMap.SelectOne", err)
	}
	return true
}

func (repo *UserRepo) UpdateMemberSavedRoles(memberRoles []string, memberId, guildId string) {
	savedRoles, err := repo.GetRolesForMember(guildId, memberId)
	if err != nil {
		log.Println("("+guildId+") "+"updateMemberSavedRoles Error while getting saved roles", err)
		return
	}
	var savedRolesIds []string
	for _, role := range savedRoles {
		savedRolesIds = append(savedRolesIds, role.RoleId)
	}

	for _, memberRole := range memberRoles {
		found := false
		for i, savedRole := range savedRolesIds {
			if savedRole == memberRole {
				found = true
				savedRolesIds[i] = ""
				break
			}
		}
		if !found {
			err = repo.AddRoleForMember(guildId, memberId, memberRole)
			if err != nil {
				log.Println("("+guildId+") Error while saving new role info", err)
				continue
			}
		}
	}

	for _, savedRole := range savedRolesIds {
		if savedRole != "" {
			err = repo.RemoveRoleForMember(guildId, memberId, savedRole)
			if err != nil {
				log.Println("("+guildId+") "+"Error while deleting info about member role", err)
				continue
			}
		}
	}
}

func (repo *UserRepo) IsUserBlacklisted(userId string, guildId string) bool {
	ret, err := repo.mysql.SelectInt("SELECT count(*) FROM Blacklists WHERE guild_id = ? AND user_id = ?", guildId, userId)
	if err != nil {
		log.Println("("+guildId+") isBlacklisted#DbMap.SelectInt", err)
		return false
	}
	if ret > 0 {
		return true
	}
	return false
}

func (repo *UserRepo) AddBlacklistForUser(userId, guildId, blacklisterId string) error {
	blacklist := Blacklist{UserId: userId, GuildId: guildId, BlacklisterId: blacklisterId}
	err := repo.mysql.Insert(&blacklist)
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepo) RemoveBlacklistForUser(userId, guildId string) error {
	_, err := repo.mysql.Exec("DELETE FROM Blacklists WHERE guild_id = ? AND user_id = ?", guildId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepo) AddHelperBlacklistForUser(userId, guildId, blacklisterId string) error {
	blacklist := HelperBlacklist{UserId: userId, GuildId: guildId, BlacklisterId: blacklisterId}
	err := repo.mysql.Insert(&blacklist)
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepo) RemoveHelperBlacklistForUser(userId, guildId string) error {
	_, err := repo.mysql.Exec("DELETE FROM HelperBlacklist WHERE guild_id = ? AND user_id = ?", guildId, userId)
	if err != nil {
		return err
	}
	return nil
}
