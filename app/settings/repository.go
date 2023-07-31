package settings

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

func (r *repositoryImpl) UpdateChannelId(ctx context.Context, guildId, name, channelId string) error {
	_, err := r.db.Exec(ctx, `UPDATE guildchannelsid SET channelid = $1 WHERE guildid = $2 AND name = $3`, channelId, guildId, name)

	return err
}

func (r *repositoryImpl) AddChannelId(ctx context.Context, guildId, name, channelId string) error {
	_, err := r.db.Exec(ctx, `INSERT INTO guildchannelsid (guildid, name, channelid) VALUES ($1,$2,$3)`, guildId, name, channelId)

	return err
}
