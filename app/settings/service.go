package settings

import (
	"context"
)

type Repository interface {
	UpdateChannelId(ctx context.Context, guildId, name, channelId string) error
	AddChannelId(ctx context.Context, guildId, name, channelId string) error
}

func NewService(repo Repository) Service {
	return &serviceImplementation{
		repo: repo,
	}
}

type serviceImplementation struct {
	repo Repository
}

func (s *serviceImplementation) UpdateChannelId(ctx context.Context, guildId, name, channelId string) error {
	return s.repo.UpdateChannelId(ctx, guildId, name, channelId)
}

func (s *serviceImplementation) AddChannelId(ctx context.Context, guildId, name, channelId string) error {
	return s.repo.AddChannelId(ctx, guildId, name, channelId)
}
