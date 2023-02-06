package settings

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

func (r *repositoryImpl) UpdateChannelId(ctx context.Context, guildId, name, channelId string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE guildchannelsid SET channelid = ? WHERE guildid = ? AND name = ?`, channelId, guildId, name)

	return err
}

func (r *repositoryImpl) AddChannelId(ctx context.Context, guildId, name, channelId string) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO guildchannelsid (guildid, name, channelid) VALUES (?,?,?)`, guildId, name, channelId)

	return err
}
