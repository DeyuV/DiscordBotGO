package serverstatus

import (
	"DiscordBotGO/pkg/config"
	"context"
	"errors"

	"github.com/bwmarrin/discordgo"
)

type Repository interface {
	GetChannelIdByNameAndGuildID(ctx context.Context, guildId, name string) (string, error)
	GetMessageIdByNameAndGuildID(ctx context.Context, guildId, name string) (string, error)
	UpdateChannelId(ctx context.Context, guildId, name, channelId string) error
	UpdateMessageId(ctx context.Context, guildId, name, messageId string) error
	AddMessageId(ctx context.Context, guildId, name, messageId string) error
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

func (s *serviceImplementation) InteractionResponse(session *discordgo.Session, i *discordgo.InteractionCreate, name, content string) error {
	serverStatusChannelId, err := s.GetChannelIdByNameAndGuildID(context.Background(), i.GuildID, config.ServerStatus)
	if err != nil {
		return err
	}

	if serverStatusChannelId != i.ChannelID {
		err = session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Wrong channel for server status commands",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})

		if err != nil {
			return err
		}

		return errors.New("wrong-channel")
	}

	channel, err := session.Channel(i.ChannelID)
	if err != nil {
		return err
	}

	editChannel := discordgo.ChannelEdit{
		Name:     name,
		Position: channel.Position,
	}

	_, err = session.ChannelEdit(i.ChannelID, &editChannel)
	if err != nil {
		return err
	}

	err = session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *serviceImplementation) StrikethroughMessage(session *discordgo.Session, i *discordgo.InteractionCreate, messageName, serverStatusChannelID string) error {
	messageId, err := s.GetMessageIdByNameAndGuildID(context.Background(), i.GuildID, messageName)
	if err != nil {
		return err
	}

	message, err := session.ChannelMessage(i.ChannelID, messageId)
	if err != nil {
		return err
	}

	if message.ID != "" {
		if message.Content[0] != '~' {
			session.ChannelMessageEdit(serverStatusChannelID, message.ID, "~~"+message.Content+"~~")
		}
	}

	return nil
}

func (s *serviceImplementation) GetChannelIdByNameAndGuildID(ctx context.Context, guildId, name string) (string, error) {
	return s.repo.GetChannelIdByNameAndGuildID(ctx, guildId, name)
}

func (s *serviceImplementation) GetMessageIdByNameAndGuildID(ctx context.Context, guildId, name string) (string, error) {
	return s.repo.GetMessageIdByNameAndGuildID(ctx, guildId, name)
}

func (s *serviceImplementation) UpdateChannelId(ctx context.Context, guildId, name, channelId string) error {
	return s.repo.UpdateChannelId(ctx, guildId, name, channelId)
}

func (s *serviceImplementation) UpdateMessageId(ctx context.Context, guildId, name, messageId string) error {
	return s.repo.UpdateMessageId(ctx, guildId, name, messageId)
}

func (s *serviceImplementation) AddChannelId(ctx context.Context, guildId, name, channelId string) error {
	return s.repo.AddChannelId(ctx, guildId, name, channelId)
}

func (s *serviceImplementation) AddMessageId(ctx context.Context, guildId, name, messageId string) error {
	return s.repo.AddMessageId(ctx, guildId, name, messageId)
}
