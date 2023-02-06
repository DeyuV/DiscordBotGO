package strategicpoint

import (
	"context"
	"database/sql"
	"time"
)

func NewRepository(db *sql.DB) Repository {
	return &repositoryImpl{db: db}
}

type repositoryImpl struct {
	db *sql.DB
}

func (r *repositoryImpl) GetEmojiByName(ctx context.Context, guildId, emojiName string) string {
	var emojiId string

	err := r.db.QueryRowContext(ctx, `SELECT emojiid FROM guildemojis WHERE guildid = ? AND emojiname = ?`, guildId, emojiName).Scan(&emojiId)
	if err != nil {
		return ""
	}
	return emojiId
}

func (r *repositoryImpl) GetChannelIdByNameAndGuildID(ctx context.Context, guildId, name string) (string, error) {
	var channelId string

	err := r.db.QueryRowContext(ctx, `SELECT channelid FROM guildchannelsid WHERE guildid = ? AND name = ?`, guildId, name).Scan(&channelId)
	if err != nil {
		return "", err
	}
	return channelId, nil
}

func (r *repositoryImpl) UpdateMessageId(ctx context.Context, guildId, name, messageId string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE guildmessagesid SET messageid = ? WHERE guildid = ? AND name = ?`, messageId, guildId, name)

	return err
}

func (r *repositoryImpl) AddMessageId(ctx context.Context, guildId, name, messageId string) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO guildmessagesid (guildid, name, messageid) VALUES (?,?,?)`, guildId, name, messageId)

	return err
}

func (r *repositoryImpl) DeleteMessageId(ctx context.Context, guildId, messageId string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM guildmessagesid WHERE guildid = ? AND messageid = ?`, guildId, messageId)

	return err
}

func (r *repositoryImpl) GetMessageIdByNameAndGuildID(ctx context.Context, guildId, name string) (string, error) {
	var messageId string

	err := r.db.QueryRowContext(ctx, `SELECT messageid FROM guildmessagesid WHERE guildid = ? AND name = ?`, guildId, name).Scan(&messageId)
	if err != nil {
		return "", err
	}
	return messageId, nil
}

func (r *repositoryImpl) AddSP(ctx context.Context, guildId, mapName, spawnTime, winningNation, userSpawning, userInteracting string) (int, error) {
	var spId int
	err := r.db.QueryRowContext(ctx, `INSERT INTO guildlogsp (guildid, map, spawntime, winningnation, userspawning, userinteracting, spdate) VALUES (?,?,?,?,?,?,?) RETURNING id`, guildId, mapName, spawnTime, winningNation, userSpawning, userInteracting, time.Now()).Scan(&spId)

	if err != nil {
		return 0, err
	}

	return spId, nil
}

func (r *repositoryImpl) DeleteSP(ctx context.Context, id int) error {
	r.db.ExecContext(context.Background(), `PRAGMA foreign_keys = ON`)
	_, err := r.db.ExecContext(ctx, `DELETE FROM guildlogsp WHERE id = ?`, id)

	return err
}

func (r *repositoryImpl) UpdateSPmap(ctx context.Context, id int, mapName string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE guildlogsp SET map = ? WHERE id = ?`, mapName, id)

	return err
}

func (r *repositoryImpl) UpdateSPspawnTime(ctx context.Context, id int, spawnTime string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE guildlogsp SET spawntime = ? WHERE id = ?`, spawnTime, id)

	return err
}

func (r *repositoryImpl) UpdateSPwinningNation(ctx context.Context, id int, winningNation string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE guildlogsp SET winningnation = ? WHERE id = ?`, winningNation, id)

	return err
}

func (r *repositoryImpl) UpdateSPmodified(ctx context.Context, id int, modified string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE guildlogsp SET modified = ? WHERE id = ?`, modified, id)

	return err
}

func (r *repositoryImpl) GetGuildId(ctx context.Context, id int) (string, error) {
	var guildId string

	err := r.db.QueryRowContext(ctx, `SELECT guildid FROM guildlogsp WHERE id = ?`, id).Scan(&guildId)

	return guildId, err
}

func (r *repositoryImpl) GetAllSPLogsByGuild(ctx context.Context, guildId string) ([]SPLogs, error) {
	var spLogs []SPLogs

	rows, err := r.db.QueryContext(ctx, `SELECT id, guildid, map, spawntime, winningnation, userspawning, userinteracting, spdate FROM guildlogsp WHERE guildid = ?`, guildId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var spLog SPLogs

		err := rows.Scan(&spLog.ID, &spLog.GuildID, &spLog.MapName, &spLog.SpawnTime, &spLog.WinningNation, &spLog.UserSpawning, &spLog.UserInteracting, &spLog.SPDate)
		if err != nil {
			return nil, err
		}

		spLogs = append(spLogs, spLog)
	}

	return spLogs, nil
}

func (r *repositoryImpl) GetSPbyGuildAndId(ctx context.Context, guildId string, spId int) error {
	row := r.db.QueryRowContext(ctx, `SELECT * FROM guildlogsp WHERE guildid = ? AND id = ?`, guildId, spId)

	return row.Err()
}
