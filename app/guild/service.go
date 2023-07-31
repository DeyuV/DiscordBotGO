package guild

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Repository interface {
	Add(ctx context.Context, guildId, guildName string) error
	GetById(ctx context.Context, guildId string) (*Guilds, error)
	DeleteGuild(ctx context.Context, guildId string) error
	GetSlashCommands(ctx context.Context, guildId string) ([]Guildcommands, error)
	AddDefaultCommands(ctx context.Context, guildId string) error
	AddEmojis(ctx context.Context, guildId, emojiId, emojiName string, animated bool) error
	DeleteEmojis(ctx context.Context, guildId string) error
	DeleteDefaultCommands(ctx context.Context, guildId string) error
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

func (s *serviceImplementation) DeleteDefaultCommands(ctx context.Context, guildId string) error {
	return s.repo.DeleteDefaultCommands(ctx, guildId)
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

// Delete all application commands for a guild
func (s *serviceImplementation) DeleteGuildCommands(session *discordgo.Session, guildID string) error {
	commands, err := session.ApplicationCommands(session.State.User.ID, guildID)
	if err != nil {
		return err
	}

	for _, command := range commands {
		err := session.ApplicationCommandDelete(session.State.User.ID, guildID, command.ID)
		if err != nil {
			return fmt.Errorf("error deleting command %s: %w", command.Name, err)
		}
	}

	return nil
}

func (s *serviceImplementation) AddGuildCommands(session *discordgo.Session, guildID, guildName string) error {
	slashCommands, err := s.GetSlashCommandsByGuildId(context.Background(), guildID)
	if err != nil {
		return fmt.Errorf("failed to get slash commands: %w", err)
	}
	var commands []*discordgo.ApplicationCommand

	for _, slashCommand := range slashCommands {
		var command discordgo.ApplicationCommand
		name := slashCommand.Command.CommandName
		description := slashCommand.Command.CommandDescription
		ok := true
		if name == "server-online" || name == "server-offline" || name == "server-maint" {
			ok = false
			command = discordgo.ApplicationCommand{
				Name:        name,
				Description: description,
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "message",
						Description: "Input text to show as message",
						Required:    true,
					},
				},
			}
		}

		if ok {
			command = discordgo.ApplicationCommand{Name: name, Description: description}
		}
		commands = append(commands, &command)
	}
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))

	for i, v := range commands {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, guildID, v)
		if err != nil {
			return fmt.Errorf("failed to create slash command %s for guild %s : %w", v.Name, guildName, err)
		}
		registeredCommands[i] = cmd
	}
	return nil
}
