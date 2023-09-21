package settings

import (
	"context"
	"fmt"

	"DiscordBotGO/pkg/aceonline"

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
					if perms&discordgo.PermissionAdministrator != 0 {
						if i.ApplicationCommandData().Name == "set-sp-forum-channel" {
							err = svc.AddChannelId(context.Background(), i.GuildID, aceonline.SPforum, i.ChannelID)
							if err != nil {
								err = svc.UpdateChannelId(context.Background(), i.GuildID, aceonline.SPforum, i.ChannelID)
								if err != nil {
									fmt.Println(err)
									return
								}

							}

							err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
								Type: discordgo.InteractionResponseChannelMessageWithSource,
								Data: &discordgo.InteractionResponseData{
									Content: "This channel will be used for sending forum sp notification",
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
}

func Register(bot *discordgo.Session, svc Service) {
	bot.AddHandler(SetChannelId(svc))
}
