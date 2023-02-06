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

			slashCommands, err := svc.GetSlashCommandsByGuildId(context.Background(), c.ID)
			if err != nil {
				fmt.Println("Failed to get slash commands")
				fmt.Println(err)
				return
			}
			var commands []*discordgo.ApplicationCommand

			for _, slashCommand := range slashCommands {
				var command discordgo.ApplicationCommand
				name := slashCommand.Command.CommandName
				description := slashCommand.Command.CommandDescription
				if name == "server-online" || name == "server-offline" || name == "server-maint" {
					switch name {
					case "server-online":
						description = "Set channel name to 🟢┃game-info + custom message"
					case "server-offline":
						description = "Set channel name to 🔴┃game-info❕ + custom message"
					case "server-maint":
						description = "Set channel name to 🟠┃game-info❕ + custom message"
					}

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

				switch name {
				case "add-sp":
					command = discordgo.ApplicationCommand{
						Name:        name,
						Description: description,
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionString,
								Name:        "map-name",
								Description: "Input full map name",
								Required:    true,
							},
							{
								Type:        discordgo.ApplicationCommandOptionString,
								Name:        "spawn-time",
								Description: "Input in format <t:(unix):R>",
								Required:    true,
							},
							{
								Type:        discordgo.ApplicationCommandOptionString,
								Name:        "winning-nation",
								Description: "Input short name of winning nation (ani/bcu)",
								Required:    true,
							},
						},
					}
				case "modify-sp":
					command = discordgo.ApplicationCommand{
						Name:        name,
						Description: description,
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionString,
								Name:        "id",
								Description: "Input SP id",
								Required:    true,
							},
							{
								Type:        discordgo.ApplicationCommandOptionString,
								Name:        "map-name",
								Description: "Input full map name",
								Required:    true,
							},
							{
								Type:        discordgo.ApplicationCommandOptionString,
								Name:        "spawn-time",
								Description: "Input in format <t:(unix):R>",
								Required:    true,
							},
							{
								Type:        discordgo.ApplicationCommandOptionString,
								Name:        "winning-nation",
								Description: "Input short name of winning nation (ani/bcu)",
								Required:    true,
							},
						},
					}
				case "delete-sp":
					command = discordgo.ApplicationCommand{
						Name:        name,
						Description: description,
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionString,
								Name:        "id",
								Description: "Input SP id",
								Required:    true,
							},
						},
					}
				default:
					command = discordgo.ApplicationCommand{Name: name, Description: description}
				}

				commands = append(commands, &command)
			}
			registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))

			for i, v := range commands {
				cmd, err := s.ApplicationCommandCreate(s.State.User.ID, c.ID, v)
				if err != nil {
					fmt.Println("Failed to create slash commands for " + c.Name)
					fmt.Println(err)
				}
				registeredCommands[i] = cmd
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
