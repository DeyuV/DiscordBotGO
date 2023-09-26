package strategicpoint

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewRepository(db *pgxpool.Pool) Repository {
	return &repositoryImpl{db: db}
}

type repositoryImpl struct {
	db *pgxpool.Pool
}

func (r *repositoryImpl) GetChannelIdByNameAndGuildID(ctx context.Context, guildId, name string) (string, error) {
	var channelId string

	err := r.db.QueryRow(ctx, `SELECT channelid FROM guildchannelsid WHERE guildid = $1 AND name = $2`, guildId, name).Scan(&channelId)
	if err != nil {
		return "", err
	}
	return channelId, nil
}

func (r *repositoryImpl) UpdateMessageId(ctx context.Context, guildId, name, messageId string) error {
	_, err := r.db.Exec(ctx, `UPDATE guildmessagesid SET messageid = $1 WHERE guildid = $2 AND name = $3`, messageId, guildId, name)

	return err
}

func (r *repositoryImpl) AddMessageId(ctx context.Context, guildId, name, messageId string) error {
	_, err := r.db.Exec(ctx, `INSERT INTO guildmessagesid (guildid, name, messageid) VALUES ($1,$2,$3)`, guildId, name, messageId)

	return err
}

func (r *repositoryImpl) DeleteMessageId(ctx context.Context, guildId, messageId string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM guildmessagesid WHERE guildid = $1 AND messageid = $2`, guildId, messageId)

	return err
}

func (r *repositoryImpl) GetMessageIdByNameAndGuildID(ctx context.Context, guildId, name string) (string, error) {
	var messageId string

	err := r.db.QueryRow(ctx, `SELECT messageid FROM guildmessagesid WHERE guildid = $1 AND name = $2`, guildId, name).Scan(&messageId)
	if err != nil {
		return "", err
	}
	return messageId, nil
}

func (r *repositoryImpl) AddSP(ctx context.Context, id, guildId, userSpawning string) error {
	_, err := r.db.Exec(ctx, `INSERT INTO guildlogsp (id, guildid, map, spawntime, winningnation, userspawning, userinteracting, spdate) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`, id, guildId, "-", "-", "-", userSpawning, "-", time.Now())

	return err
}

func (r *repositoryImpl) DeleteSP(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM guildlogsp WHERE id = $1`, id)

	return err
}

func (r *repositoryImpl) UpdateSP(ctx context.Context, id, mapName, spawntime, winningNation, userInteracting string) error {
	_, err := r.db.Exec(ctx, `UPDATE guildlogsp SET map = $2, spawntime = $3, winningnation = $4, userinteracting = $5  WHERE id = $1`, id, mapName, spawntime, winningNation, userInteracting)

	return err
}

func (r *repositoryImpl) GetGuildId(ctx context.Context, id int) (string, error) {
	var guildId string

	err := r.db.QueryRow(ctx, `SELECT guildid FROM guildlogsp WHERE id = $1`, id).Scan(&guildId)

	return guildId, err
}

func (r *repositoryImpl) GetAllSPLogsByGuild(ctx context.Context, guildId string) ([]SPLogs, error) {
	var spLogs []SPLogs

	rows, err := r.db.Query(ctx, `SELECT id, guildid, map, spawntime, winningnation, userspawning, userinteracting, spdate FROM guildlogsp WHERE guildid = $1`, guildId)
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

		if spLog.SPDate.Month() == time.Now().Month() && spLog.SPDate.Day() == time.Now().Day() {
			spLogs = append(spLogs, spLog)
		}
	}

	return spLogs, nil
}

func (r *repositoryImpl) GetSPbyGuildAndId(ctx context.Context, guildId, spId string) error {
	_, err := r.db.Query(ctx, `SELECT * FROM guildlogsp WHERE guildid = $1 AND id = $2`, guildId, spId)

	return err
}

func (r *repositoryImpl) UpdateChannelId(ctx context.Context, guildId, name, channelId string) error {
	_, err := r.db.Exec(ctx, `UPDATE guildchannelsid SET channelid = $1 WHERE guildid = $2 AND name = $3`, channelId, guildId, name)

	return err
}

func (r *repositoryImpl) AddChannelId(ctx context.Context, guildId, name, channelId string) error {
	_, err := r.db.Exec(ctx, `INSERT INTO guildchannelsid (guildid, name, channelid) VALUES ($1,$2,$3)`, guildId, name, channelId)

	return err
}
