package serverstatus

import (
	"context"

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

func (r *repositoryImpl) GetMessageIdByNameAndGuildID(ctx context.Context, guildId, name string) (string, error) {
	var messageId string

	err := r.db.QueryRow(ctx, `SELECT messageid FROM guildmessagesid WHERE guildid = $1 AND name = $2`, guildId, name).Scan(&messageId)
	if err != nil {
		return "", err
	}
	return messageId, nil
}

func (r *repositoryImpl) UpdateChannelId(ctx context.Context, guildId, name, channelId string) error {
	_, err := r.db.Exec(ctx, `UPDATE guildchannelsid SET channelid = $1 WHERE guildid = $2 AND name = $3`, channelId, guildId, name)

	return err
}

func (r *repositoryImpl) UpdateMessageId(ctx context.Context, guildId, name, messageId string) error {
	_, err := r.db.Exec(ctx, `UPDATE guildmessagesid SET messageid = $1 WHERE guildid = $2 AND name = $3`, messageId, guildId, name)

	return err
}

func (r *repositoryImpl) AddMessageId(ctx context.Context, guildId, name, messageId string) error {
	_, err := r.db.Exec(ctx, `INSERT INTO guildmessagesid (guildid, name, messageid) VALUES ($1,$2,$3)`, guildId, name, messageId)

	return err
}

func (r *repositoryImpl) AddChannelId(ctx context.Context, guildId, name, channelId string) error {
	_, err := r.db.Exec(ctx, `INSERT INTO guildchannelsid (guildid, name, channelid) VALUES ($1,$2,$3)`, guildId, name, channelId)

	return err
}
