package guild

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
	DeleteGuildCommands(session *discordgo.Session, guildID string) error
	AddGuildCommands(session *discordgo.Session, guildID, guildName string) error
	DeleteDefaultCommands(ctx context.Context, guildId string) error
}

// AddGuild Adds a guild on database when bot creates the guilds where it is a member or when invited
func AddGuild(svc Service) func(s *discordgo.Session, c *discordgo.GuildCreate) {
	return func(s *discordgo.Session, c *discordgo.GuildCreate) {
		g, err := svc.GetGuildById(context.Background(), c.ID)

		if g == nil && err != nil {
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

			err = svc.AddGuildCommands(s, c.ID, c.Name)
			if err != nil {
				fmt.Println(err)
				return
			}

			for _, emoji := range c.Emojis {
				err := svc.AddGuildEmojis(context.Background(), c.ID, emoji.ID, emoji.Name, emoji.Animated)
				if err != nil {
					fmt.Println("Failed to add emoji for " + c.Name)
					fmt.Println(err)
					return
				}
			}

			fmt.Println("Successfully added guild: " + c.Name)
		} else {
			go func() {
				err = svc.DeleteDefaultCommands(context.Background(), c.ID)
				if err != nil {
					fmt.Println(err)
					return
				}

				err = svc.DeleteGuildCommands(s, c.ID)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println("Deleted commands for guild: " + c.Name)

				err = svc.AddDefaultCommands(context.Background(), c.ID)
				if err != nil {
					fmt.Println("Failed to insert default commands")
					fmt.Println(err)
					return
				}

				err = svc.AddGuildCommands(s, c.ID, c.Name)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println("Added commands for guild: " + c.Name)
			}()
		}
	}
}

// DeleteGuild Deletes a guild from database when bot leaves a guild
func DeleteGuild(svc Service) func(s *discordgo.Session, c *discordgo.GuildDelete) {
	return func(s *discordgo.Session, c *discordgo.GuildDelete) {
		err := svc.DeleteGuild(context.Background(), c.ID)
		if err != nil {
			fmt.Println("Failed to delete guild")
			return
		}
	}
}

func EmojiUpdate(svc Service) func(s *discordgo.Session, c *discordgo.GuildEmojisUpdate) {
	return func(s *discordgo.Session, c *discordgo.GuildEmojisUpdate) {
		// Maybe there is a better way to do this
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
	bot.AddHandler(EmojiUpdate(svc))
}
