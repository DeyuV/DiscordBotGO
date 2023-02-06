package serverstatus

import (
	"context"
	"database/sql"
)

func NewRepository(db *sql.DB) Repository {
	return &repositoryImpl{db: db}
}

type repositoryImpl struct {
	db *sql.DB
}

func (r *repositoryImpl) GetChannelIdByNameAndGuildID(ctx context.Context, guildId, name string) (string, error) {
	var channelId string

	err := r.db.QueryRowContext(ctx, `SELECT channelid FROM guildchannelsid WHERE guildid = ? AND name = ?`, guildId, name).Scan(&channelId)
	if err != nil {
		return "", err
	}
	return channelId, nil
}

func (r *repositoryImpl) GetMessageIdByNameAndGuildID(ctx context.Context, guildId, name string) (string, error) {
	var messageId string

	err := r.db.QueryRowContext(ctx, `SELECT messageid FROM guildmessagesid WHERE guildid = ? AND name = ?`, guildId, name).Scan(&messageId)
	if err != nil {
		return "", err
	}
	return messageId, nil
}

func (r *repositoryImpl) UpdateChannelId(ctx context.Context, guildId, name, channelId string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE guildchannelsid SET channelid = ? WHERE guildid = ? AND name = ?`, channelId, guildId, name)

	return err
}

func (r *repositoryImpl) UpdateMessageId(ctx context.Context, guildId, name, messageId string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE guildmessagesid SET messageid = ? WHERE guildid = ? AND name = ?`, messageId, guildId, name)

	return err
}

func (r *repositoryImpl) AddMessageId(ctx context.Context, guildId, name, messageId string) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO guildmessagesid (guildid, name, messageid) VALUES (?,?,?)`, guildId, name, messageId)

	return err
}

func (r *repositoryImpl) AddChannelId(ctx context.Context, guildId, name, channelId string) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO guildchannelsid (guildid, name, channelid) VALUES (?,?,?)`, guildId, name, channelId)

	return err
}
