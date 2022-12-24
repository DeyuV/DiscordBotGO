package Guild

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type Service interface {
	AddGuild(ctx context.Context, guildId string, guildName string) error
	GetGuildById(ctx context.Context, guildId string) (*Guilds, error)
	DeleteGuild(ctx context.Context, guildId string) error
	GetSlashCommandsByGuildId(ctx context.Context, guildId string) ([]Guildcommands, error)
	AddDefaultCommands(ctx context.Context, guildId string) error
	AddGuildEmojis(ctx context.Context, guildId, emojiId, emojiName string, animated bool) error
	DeleteGuildEmojis(ctx context.Context, guildId string) error
}

// RegisterCommands Adds slash commands to guilds
func RegisterCommands(svc Service) func(s *discordgo.Session, r *discordgo.Ready) {
	return func(s *discordgo.Session, r *discordgo.Ready) {
		for _, guild := range r.Guilds {
			slashCommands, err := svc.GetSlashCommandsByGuildId(context.Background(), guild.ID)
			if err != nil {
				fmt.Println("Failed to get slash commands")
				fmt.Println(err)
				return
			}
			var commands []*discordgo.ApplicationCommand

			for _, slashCommand := range slashCommands {
				name := slashCommand.Command.CommandName
				description := slashCommand.Command.CommandDescription
				command := discordgo.ApplicationCommand{Name: name, Description: description}
				commands = append(commands, &command)
			}

			registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
			for i, v := range commands {
				cmd, err := s.ApplicationCommandCreate(s.State.User.ID, guild.ID, v)
				if err != nil {
					fmt.Println("Failed to create slash commands")
					fmt.Println(err)
				}
				registeredCommands[i] = cmd
			}

			for _, emoji := range guild.Emojis {
				err := svc.AddGuildEmojis(context.Background(), guild.ID, emoji.ID, emoji.Name, emoji.Animated)
				if err != nil {
					fmt.Println("Failed to add emoji")
					fmt.Println(err)
					return
				}
			}
		}
	}
}

// AddGuild Adds a guild on database when bot creates the guilds where it is a member or when invited
func AddGuild(svc Service) func(s *discordgo.Session, c *discordgo.GuildCreate) {
	return func(s *discordgo.Session, c *discordgo.GuildCreate) {
		_, err := svc.GetGuildById(context.Background(), c.ID)
		if err == nil {
			fmt.Println(err)
			return
		}

		err = svc.AddGuild(context.Background(), c.ID, c.Name)
		if err != nil {
			fmt.Println("Guild insert failed")
			fmt.Println(err)
			return
		}

		err = svc.AddDefaultCommands(context.Background(), c.ID)
		if err != nil {
			fmt.Println("Failed to insert default commands")
			fmt.Println(err)
			return
		}
		fmt.Println("Successfully added guild: " + c.Name)
	}

}

// DeleteGuild Deletes a guild from database when bot leaves a guild
func DeleteGuild(svc Service) func(s *discordgo.Session, c *discordgo.GuildDelete) {
	return func(s *discordgo.Session, c *discordgo.GuildDelete) {
		err := svc.DeleteGuild(context.Background(), c.ID)
		if err != nil {
			fmt.Println("Failed to delete guild: " + c.Name)
			return
		}
	}
}

func EmojiUpdate(svc Service) func(s *discordgo.Session, c *discordgo.GuildEmojisUpdate) {
	return func(s *discordgo.Session, c *discordgo.GuildEmojisUpdate) {
		err := svc.DeleteGuildEmojis(context.Background(), c.GuildID)
		if err != nil {
			fmt.Println("Failed to delete emojis")
			fmt.Println(err)
			return
		}
		for _, emoji := range c.Emojis {
			err := svc.AddGuildEmojis(context.Background(), c.GuildID, emoji.ID, emoji.Name, emoji.Animated)
			if err != nil {
				fmt.Println("Failed to add emoji")
				fmt.Println(err)
				return
			}
		}
	}
}

func Register(bot *discordgo.Session, svc Service) {
	bot.AddHandler(AddGuild(svc))
	bot.AddHandler(DeleteGuild(svc))
	bot.AddHandler(RegisterCommands(svc))
	bot.AddHandler(EmojiUpdate(svc))
}
