package Commands

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

func NewRepository(db *pgxpool.Pool) Repository {
	return &repositoryImpl{db: db}
}

func (r *repositoryImpl) GetEmojiByName(ctx context.Context, guildId, emojiName string) string {
	var emojiId string

	err := r.db.QueryRow(ctx, `SELECT emojiid FROM guildemojis WHERE guildid = $1 AND emojiname = $2`, guildId, emojiName).Scan(&emojiId)
	if err != nil {
		return ""
	}
	return emojiId
}

type repositoryImpl struct {
	db *pgxpool.Pool
}
