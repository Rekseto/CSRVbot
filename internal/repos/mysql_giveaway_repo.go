package repos

import (
	"database/sql"
	"gopkg.in/gorp.v2"
	"log"
	"time"
)

type GiveawayRepo struct {
	mysql *gorp.DbMap
}

func NewGiveawayRepo(mysql *gorp.DbMap) *GiveawayRepo {
	mysql.AddTableWithName(Giveaway{}, "Giveaways").SetKeys(true, "id")
	mysql.AddTableWithName(Participant{}, "Participants").SetKeys(true, "id")
	mysql.AddTableWithName(ParticipantCandidate{}, "ParticipantCandidates").SetKeys(true, "id")
	mysql.AddTableWithName(ThxNotification{}, "ThxNotifications").SetKeys(true, "id")

	return &GiveawayRepo{mysql: mysql}
}

type Giveaway struct {
	Id            int            `db:"id, primarykey, autoincrement"`
	StartTime     time.Time      `db:"start_time"`
	EndTime       *time.Time     `db:"end_time"`
	GuildId       string         `db:"guild_id,size:255"`
	GuildName     string         `db:"guild_name,size:255"`
	WinnerId      sql.NullString `db:"winner_id,size:255"`
	WinnerName    sql.NullString `db:"winner_name,size:255"`
	InfoMessageId sql.NullString `db:"info_message_id,size:255"`
	Code          sql.NullString `db:"code,size:255"`
}

type Participant struct {
	Id           int            `db:"id, primarykey, autoincrement"`
	UserName     string         `db:"user_name,size:255"`
	UserId       string         `db:"user_id,size:255"`
	GiveawayId   int            `db:"giveaway_id"`
	CreateTime   time.Time      `db:"create_time"`
	GuildName    string         `db:"guild_name,size:255"`
	GuildId      string         `db:"guild_id,size:255"`
	MessageId    string         `db:"message_id,size:255"`
	ChannelId    string         `db:"channel_id,size:255"`
	IsAccepted   sql.NullBool   `db:"is_accepted"`
	AcceptTime   *time.Time     `db:"accept_time"`
	AcceptUser   sql.NullString `db:"accept_user,size:255"`
	AcceptUserId sql.NullString `db:"accept_user_id,size:255"`
}

type ParticipantCandidate struct {
	Id                    int          `db:"id, primarykey, autoincrement"`
	CandidateId           string       `db:"candidate_id,size:255"`
	CandidateName         string       `db:"candidate_name,size:255"`
	CandidateApproverId   string       `db:"candidate_approver_id,size:255"`
	CandidateApproverName string       `db:"candidate_approver_name,size:255"`
	GuildId               string       `db:"guild_id,size:255"`
	GuildName             string       `db:"guild_name,size:255"`
	MessageId             string       `db:"message_id,size:255"`
	ChannelId             string       `db:"channel_id,size:255"`
	IsAccepted            sql.NullBool `db:"is_accepted"`
	AcceptTime            *time.Time   `db:"accept_time"`
}

type ThxNotification struct {
	Id                       int    `db:"id,primarykey,autoincrement"`
	MessageId                string `db:"message_id,size:255"`
	ThxNotificationMessageId string `db:"thx_notification_message_id,size:255"`
}

type ParticipantWithThxAmount struct {
	UserId    string `db:"user_id,size:255"`
	ThxAmount int    `db:"amount"`
}

func (repo *GiveawayRepo) GetGiveawayForGuild(guildId string) (*Giveaway, error) {
	var giveaway Giveaway
	err := repo.mysql.SelectOne(&giveaway, "SELECT * FROM Giveaways WHERE guild_id = ? AND end_time IS NULL", guildId)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &giveaway, nil
}

func (repo *GiveawayRepo) GetParticipantNamesForGiveaway(giveawayId int) ([]string, error) {
	var participants []Participant
	_, err := repo.mysql.Select(&participants, "SELECT user_name FROM Participants WHERE giveaway_id = ? AND is_accepted = true", giveawayId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		log.Println("getParticipantsNames#DbMap.Select", err)
		return nil, err
	}
	names := make([]string, len(participants))
	for i := range participants {
		names[i] = participants[i].UserName
	}
	return names, nil
}

func (repo *GiveawayRepo) InsertGiveaway(guildId string, guildName string) error {
	giveaway := &Giveaway{
		StartTime: time.Now(),
		GuildId:   guildId,
		GuildName: guildName,
	}
	err := repo.mysql.Insert(giveaway)
	if err != nil {
		return err
	}
	return nil
}

func (repo *GiveawayRepo) InsertParticipant(giveawayId int, guildId, guildName, userId, userName, channelId, messageId string) error {
	participant := &Participant{
		UserId:     userId,
		GiveawayId: giveawayId,
		CreateTime: time.Now(),
		GuildId:    guildId,
		ChannelId:  channelId,
		GuildName:  guildName,
		UserName:   userName,
		MessageId:  messageId,
	}
	err := repo.mysql.Insert(participant)
	if err != nil {
		return err
	}
	return nil
}

func (repo *GiveawayRepo) InsertParticipantCandidate(guildId, guildName, candidateId, candidateName, approverId, approverName, channelId, messageId string) error {
	participantCandidate := &ParticipantCandidate{
		CandidateName:         candidateName,
		CandidateId:           candidateId,
		CandidateApproverName: approverName,
		CandidateApproverId:   approverId,
		GuildName:             guildName,
		GuildId:               guildId,
		ChannelId:             channelId,
		MessageId:             messageId,
	}
	err := repo.mysql.Insert(participantCandidate)
	if err != nil {
		return err
	}
	return nil
}

func (repo *GiveawayRepo) GetParticipantsWithThxAmount(guildId string, minThxAmount int) ([]ParticipantWithThxAmount, error) {
	var helpers []ParticipantWithThxAmount
	_, err := repo.mysql.Select(&helpers, "SELECT * FROM (SELECT user_id, count(*) AS amount FROM Participants WHERE guild_id=? AND is_accepted=1 GROUP BY user_id) AS a WHERE amount > ?", guildId, minThxAmount)
	if err != nil {
		return nil, err
	}
	return helpers, nil
}

func (repo *GiveawayRepo) GetParticipantsForGiveaway(giveawayId int) ([]Participant, error) {
	var participants []Participant
	_, err := repo.mysql.Select(&participants, "SELECT * FROM Participants WHERE giveaway_id = ? AND is_accepted = true", giveawayId)
	if err != nil {
		return nil, err
	}
	return participants, nil
}

func (repo *GiveawayRepo) GetThxNotification(messageId string) (*ThxNotification, error) {
	var notification ThxNotification
	err := repo.mysql.SelectOne(&notification, "SELECT * FROM ThxNotifications WHERE message_id = ?", messageId)
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

func (repo *GiveawayRepo) InsertThxNotification(thxMessageId string, notificationMessageId string) error {
	notification := &ThxNotification{
		MessageId:                thxMessageId,
		ThxNotificationMessageId: notificationMessageId,
	}
	err := repo.mysql.Insert(notification)
	if err != nil {
		return err
	}
	return nil
}

func (repo *GiveawayRepo) IsThxMessage(messageId string) bool {
	ret, err := repo.mysql.SelectInt("SELECT count(*) FROM Participants WHERE message_id = ?", messageId)
	if err != nil {
		log.Println("IsThxMessage#DbMap.SelectInt", err)
		return false
	}

	return ret == 1
}

func (repo *GiveawayRepo) IsThxmeMessage(messageId string) bool {
	ret, err := repo.mysql.SelectInt("SELECT count(*) FROM ParticipantCandidates WHERE message_id = ?", messageId)
	if err != nil {
		log.Println("IsThxMessage#DbMap.SelectInt", err)
		return false
	}

	return ret == 1
}

func (repo *GiveawayRepo) GetParticipant(messageId string) (*Participant, error) {
	var participant Participant
	err := repo.mysql.SelectOne(&participant, "SELECT * FROM Participants WHERE message_id = ?", messageId) //fixme error when selecting accepted thx sql: Scan error on column index 10, name "accept_time": unsupported Scan, storing driver.Value type []uint8 into type *time.Time
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &participant, nil
}

func (repo *GiveawayRepo) UpdateParticipant(participant *Participant, acceptUserId, acceptUsername string, isAccepted bool) error {
	now := time.Now()
	participant.AcceptTime = &now
	participant.IsAccepted = sql.NullBool{Bool: isAccepted, Valid: true}
	participant.AcceptUserId = sql.NullString{String: acceptUserId, Valid: true}
	participant.AcceptUser = sql.NullString{String: acceptUsername, Valid: true}
	_, err := repo.mysql.Update(participant)
	if err != nil {
		return err
	}
	return nil
}

func (repo *GiveawayRepo) GetParticipantCandidate(messageId string) (*ParticipantCandidate, error) {
	var participantCandidate ParticipantCandidate
	err := repo.mysql.SelectOne(&participantCandidate, "SELECT * FROM ParticipantCandidates WHERE message_id = ?", messageId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &participantCandidate, nil
}

func (repo *GiveawayRepo) UpdateParticipantCandidate(participantCandidate *ParticipantCandidate, isAccepted bool) error {
	now := time.Now()
	participantCandidate.AcceptTime = &now
	participantCandidate.IsAccepted = sql.NullBool{Bool: isAccepted, Valid: true}
	_, err := repo.mysql.Update(participantCandidate)
	if err != nil {
		return err
	}
	return nil
}

func (repo *GiveawayRepo) UpdateGiveaway(giveaway *Giveaway, messageId, code, winnerId, winnerName string) error {
	now := time.Now()
	giveaway.EndTime = &now
	giveaway.InfoMessageId = sql.NullString{String: messageId, Valid: true}
	giveaway.Code = sql.NullString{String: code, Valid: true}
	giveaway.WinnerId = sql.NullString{String: winnerId, Valid: true}
	giveaway.WinnerName = sql.NullString{String: winnerName, Valid: true}
	_, err := repo.mysql.Update(giveaway)
	if err != nil {
		return err
	}
	return nil
}

func (repo *GiveawayRepo) GetUnfinishedGiveaways() ([]Giveaway, error) {
	var giveaways []Giveaway
	_, err := repo.mysql.Select(&giveaways, "SELECT * FROM Giveaways WHERE end_time IS NULL")
	if err != nil {
		return nil, err
	}
	return giveaways, nil
}

func (repo *GiveawayRepo) RemoveParticipants(giveawayId int, participantId string) error {
	_, err := repo.mysql.Exec("UPDATE Participants SET is_accepted=false WHERE giveaway_id = ? AND user_id = ?", giveawayId, participantId)
	if err != nil {
		return err
	}
	return nil
}
