package settings

import (
	"context"
	"fmt"
	"strings"

	"DiscordBotGO/pkg/config"

	"github.com/bwmarrin/discordgo"
)

type Service interface {
	UpdateChannelId(ctx context.Context, guildId, name, channelId string) error
	AddChannelId(ctx context.Context, guildId, name, channelId string) error
}

func SetChannelId(svc Service) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		perms, err := s.UserChannelPermissions(i.Member.User.ID, i.ChannelID)
		if err != nil {
			fmt.Println(err)
			return
		}

		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			{
				if perms&discordgo.PermissionAdministrator != 0 {
					if i.ApplicationCommandData().Name == "set-server-status-channel" {
						err = svc.AddChannelId(context.Background(), i.GuildID, config.ServerStatus, i.ChannelID)
						if err != nil {
							if strings.Split(err.Error(), " ")[0] == "UNIQUE" {
								err = svc.UpdateChannelId(context.Background(), i.GuildID, config.ServerStatus, i.ChannelID)
								if err != nil {
									fmt.Println(err)
									return
								}
							}
						}

						err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "This will be channel used for server status commands",
								Flags:   discordgo.MessageFlagsEphemeral,
							},
						})

						if err != nil {
							fmt.Println(err)
							return
						}
					}

					if i.ApplicationCommandData().Name == "set-sp-channel" {
						err = svc.AddChannelId(context.Background(), i.GuildID, config.Strategicpoint, i.ChannelID)
						if err != nil {
							if strings.Split(err.Error(), " ")[0] == "UNIQUE" {
								err = svc.UpdateChannelId(context.Background(), i.GuildID, config.Strategicpoint, i.ChannelID)
								if err != nil {
									fmt.Println(err)
									return
								}
							}
						}

						err = svc.AddChannelId(context.Background(), i.GuildID, config.LogStrategicpoint, i.ChannelID)
						if err != nil {
							if strings.Split(err.Error(), " ")[0] == "UNIQUE" {
								err = svc.UpdateChannelId(context.Background(), i.GuildID, config.LogStrategicpoint, i.ChannelID)
								if err != nil {
									fmt.Println(err)
									return
								}
							}
						}

						err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "This will be channel used for SP commands",
								Flags:   discordgo.MessageFlagsEphemeral,
							},
						})

						if err != nil {
							fmt.Println(err)
							return
						}
					}

					if i.ApplicationCommandData().Name == "set-admin-sp-channel" {
						err = svc.AddChannelId(context.Background(), i.GuildID, config.AdminLogStrategicpoint, i.ChannelID)
						if err != nil {
							if strings.Split(err.Error(), " ")[0] == "UNIQUE" {
								err = svc.UpdateChannelId(context.Background(), i.GuildID, config.AdminLogStrategicpoint, i.ChannelID)
								if err != nil {
									fmt.Println(err)
									return
								}
							}
						}

						err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "This will be channel used for admin log SP commands",
								Flags:   discordgo.MessageFlagsEphemeral,
							},
						})

						if err != nil {
							fmt.Println(err)
							return
						}
					}
				}
			}
		}
	}
}

func Register(bot *discordgo.Session, svc Service) {
	bot.AddHandler(SetChannelId(svc))
}
