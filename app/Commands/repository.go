package Commands

import (
	"context"
	"database/sql"
)

func NewRepository(db *sql.DB) Repository {
	return &repositoryImpl{db: db}
}

func (r *repositoryImpl) GetEmojiByName(ctx context.Context, guildId, emojiName string) string {
	var emojiId string

	err := r.db.QueryRowContext(ctx, `SELECT emojiid FROM guildemojis WHERE guildid = ? AND emojiname = ?`, guildId, emojiName).Scan(&emojiId)
	if err != nil {
		return ""
	}
	return emojiId
}

type repositoryImpl struct {
	db *sql.DB
}
