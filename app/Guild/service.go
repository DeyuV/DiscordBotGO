package Guild

import (
	"context"
)

type Repository interface {
	Add(ctx context.Context, guildId, guildName string) error
	GetById(ctx context.Context, guildId string) (*Guilds, error)
	DeleteGuild(ctx context.Context, guildId string) error
	GetSlashCommands(ctx context.Context, guildId string) ([]Guildcommands, error)
	AddDefaultCommands(ctx context.Context, guildId string) error
	AddEmojis(ctx context.Context, guildId, emojiId, emojiName string, animated bool) error
	DeleteEmojis(ctx context.Context, guildId string) error
}

func NewService(repo Repository) Service {
	return &serviceImplementation{
		repo: repo,
	}
}

type serviceImplementation struct {
	repo Repository
}

func (s *serviceImplementation) DeleteGuildEmojis(ctx context.Context, guildId string) error {
	return s.repo.DeleteEmojis(ctx, guildId)
}

func (s *serviceImplementation) AddGuildEmojis(ctx context.Context, guildId, emojiId, emojiName string, animated bool) error {
	return s.repo.AddEmojis(ctx, guildId, emojiId, emojiName, animated)
}

func (s *serviceImplementation) AddDefaultCommands(ctx context.Context, guildId string) error {
	return s.repo.AddDefaultCommands(ctx, guildId)
}

func (s *serviceImplementation) GetSlashCommandsByGuildId(ctx context.Context, guildId string) ([]Guildcommands, error) {
	return s.repo.GetSlashCommands(ctx, guildId)
}

func (s *serviceImplementation) AddGuild(ctx context.Context, guildId string, guildName string) error {
	return s.repo.Add(ctx, guildId, guildName)
}

func (s *serviceImplementation) GetGuildById(ctx context.Context, guildId string) (*Guilds, error) {
	return s.repo.GetById(ctx, guildId)
}

func (s *serviceImplementation) DeleteGuild(ctx context.Context, guildId string) error {
	return s.repo.DeleteGuild(ctx, guildId)
}
