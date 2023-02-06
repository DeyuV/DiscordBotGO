package serverstatus

import (
	"GOdiscordBOT/pkg/config"
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Service interface {
	GetChannelIdByNameAndGuildID(ctx context.Context, guildId, name string) (string, error)
	GetMessageIdByNameAndGuildID(ctx context.Context, guildId, name string) (string, error)
	UpdateChannelId(ctx context.Context, guildId, name, channelId string) error
	UpdateMessageId(ctx context.Context, guildId, name, messageId string) error
	AddChannelId(ctx context.Context, guildId, name, channelId string) error
	AddMessageId(ctx context.Context, guildId, name, messageId string) error

	InteractionResponse(session *discordgo.Session, i *discordgo.InteractionCreate, name, content string) error
	StrikethroughMessage(session *discordgo.Session, i *discordgo.InteractionCreate, messageName, serverStatusChannelID string) error
}

func ServerStatus(svc Service) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		perms, err := s.UserChannelPermissions(i.Member.User.ID, i.ChannelID)
		if err != nil {
			fmt.Println(err)
			return
		}

		if perms&discordgo.PermissionAdministrator != 0 {
			switch i.Type {
			case discordgo.InteractionApplicationCommand:
				{
					if i.ApplicationCommandData().Name == "server-online" {
						err := svc.InteractionResponse(s, i, "🟢┃game-info", "You changed channel name to 🟢┃game-info")
						if err != nil {
							fmt.Println(err)
							return
						}

						serverStatusChannelID, err := svc.GetChannelIdByNameAndGuildID(context.Background(), i.GuildID, config.ServerStatus)
						if err != nil {
							fmt.Println(err)
							return
						}

						message, err := s.ChannelMessageSend(serverStatusChannelID, "<t:"+strconv.Itoa(int(time.Now().Unix()))+":R> 🟢 **"+i.ApplicationCommandData().Options[0].StringValue()+"**")
						if err != nil {
							fmt.Println(err)
							return
						}

						err = svc.StrikethroughMessage(s, i, config.ServerOnline, serverStatusChannelID)
						if err != nil {
							if err == sql.ErrNoRows {
								fmt.Println("server online message does not exist")
							} else {
								fmt.Println(err)
							}
						}

						err = svc.AddMessageId(context.Background(), i.GuildID, config.ServerOnline, message.ID)
						if err != nil {
							if strings.Split(err.Error(), " ")[0] == "UNIQUE" {
								err = svc.UpdateMessageId(context.Background(), i.GuildID, config.ServerOnline, message.ID)
								if err != nil {
									fmt.Println(err)
									return
								}
							}
						}

						err = svc.StrikethroughMessage(s, i, config.ServerOffline, serverStatusChannelID)
						if err != nil {
							if err == sql.ErrNoRows {
								fmt.Println("server offfline message does not exist")
							} else {
								fmt.Println(err)
							}
						}

						err = svc.StrikethroughMessage(s, i, config.ServerMaintenance, serverStatusChannelID)
						if err != nil {
							if err == sql.ErrNoRows {
								fmt.Println("server maintenance message does not exist")
							} else {
								fmt.Println(err)
							}
						}
					}

					if i.ApplicationCommandData().Name == "server-offline" {
						err := svc.InteractionResponse(s, i, "🔴┃game-info❕", "You changed channel name to 🔴┃game-info❕")
						if err != nil {
							fmt.Println(err)
							return
						}

						serverStatusChannelID, err := svc.GetChannelIdByNameAndGuildID(context.Background(), i.GuildID, config.ServerStatus)
						if err != nil {
							fmt.Println(err)
							return
						}

						message, err := s.ChannelMessageSend(serverStatusChannelID, "<t:"+strconv.Itoa(int(time.Now().Unix()))+":R> 🔴 **"+i.ApplicationCommandData().Options[0].StringValue()+"**")
						if err != nil {
							fmt.Println(err)
							return
						}

						err = svc.StrikethroughMessage(s, i, config.ServerOffline, serverStatusChannelID)
						if err != nil {
							if err == sql.ErrNoRows {
								fmt.Println("server offfline message does not exist")
							} else {
								fmt.Println(err)
							}
						}

						err = svc.AddMessageId(context.Background(), i.GuildID, config.ServerOffline, message.ID)
						if err != nil {
							if strings.Split(err.Error(), " ")[0] == "UNIQUE" {
								err = svc.UpdateMessageId(context.Background(), i.GuildID, config.ServerOffline, message.ID)
								if err != nil {
									fmt.Println(err)
									return
								}
							}
						}

						err = svc.StrikethroughMessage(s, i, config.ServerOnline, serverStatusChannelID)
						if err != nil {
							if err == sql.ErrNoRows {
								fmt.Println("server online message does not exist")
							} else {
								fmt.Println(err)
							}
						}

						err = svc.StrikethroughMessage(s, i, config.ServerMaintenance, serverStatusChannelID)
						if err != nil {
							if err == sql.ErrNoRows {
								fmt.Println("server maintenance message does not exist")
							} else {
								fmt.Println(err)
							}
						}
					}

					if i.ApplicationCommandData().Name == "server-maint" {
						err := svc.InteractionResponse(s, i, "🟠┃game-info❕", "You changed channel name to 🟠┃game-info❕")
						if err != nil {
							fmt.Println(err)
							return
						}

						serverStatusChannelID, err := svc.GetChannelIdByNameAndGuildID(context.Background(), i.GuildID, config.ServerStatus)
						if err != nil {
							fmt.Println(err)
							return
						}

						message, err := s.ChannelMessageSend(serverStatusChannelID, "<t:"+strconv.Itoa(int(time.Now().Unix()))+":R> 🟠 **"+i.ApplicationCommandData().Options[0].StringValue()+"**")
						if err != nil {
							fmt.Println(err)
							return
						}

						err = svc.StrikethroughMessage(s, i, config.ServerMaintenance, serverStatusChannelID)
						if err != nil {
							if err == sql.ErrNoRows {
								fmt.Println("server maintenance message does not exist")
							} else {
								fmt.Println(err)
							}
						}

						err = svc.AddMessageId(context.Background(), i.GuildID, config.ServerMaintenance, message.ID)
						if err != nil {
							if strings.Split(err.Error(), " ")[0] == "UNIQUE" {
								err = svc.UpdateMessageId(context.Background(), i.GuildID, config.ServerMaintenance, message.ID)
								if err != nil {
									fmt.Println(err)
									return
								}
							}
						}

						err = svc.StrikethroughMessage(s, i, config.ServerOffline, serverStatusChannelID)
						if err != nil {
							if err == sql.ErrNoRows {
								fmt.Println("server offfline message does not exist")
							} else {
								fmt.Println(err)
							}
						}

						err = svc.StrikethroughMessage(s, i, config.ServerOnline, serverStatusChannelID)
						if err != nil {
							if err == sql.ErrNoRows {
								fmt.Println("server online message does not exist")
							} else {
								fmt.Println(err)
							}
						}
					}
				}
			}
		}
	}
}

func Register(bot *discordgo.Session, svc Service) {
	bot.AddHandler(ServerStatus(svc))
}
