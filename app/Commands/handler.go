package Commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	logSpMessageID        string
	logSpChannelID        string
	spHistory             = false
	serverOnlineMessage   Message
	serverOfflineMessage  Message
	serverMaintMessage    Message
	serverStatusChannelID string
)

type Service interface {
	GetEmojiByName(ctx context.Context, guildId, emojiName string) string
	ToggleCommand(ctx context.Context, name string)
	GetImageURL(name string) string
}

func ToggleCommands(svc Service) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {

	}
}

func ResetSpLog(s *discordgo.Session) {
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  0x000000,
		Title:  "STRATEGIC POINT HISTORY",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Map: ",
				Value:  "----------",
				Inline: true,
			},
			{
				Name:   "Spawn Time: ",
				Value:  "----------",
				Inline: true,
			},
			{
				Name:   "Winning Nation: ",
				Value:  "----------",
				Inline: true,
			},
		},
	}

	_, err := s.ChannelMessageEditEmbed(logSpChannelID, logSpMessageID, embed)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func InitResetSpLog(s *discordgo.Session) {
	t := time.Now()
	n := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 0, 0, t.Location())
	d := n.Sub(t)
	if d < 0 {
		n = n.Add(24 * time.Hour)
		d = n.Sub(t)
	}
	for {
		time.Sleep(d)
		d = 24 * time.Hour
		ResetSpLog(s)
	}
}

func HandleSpNotification(svc Service) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			if m.Content == "Please insert a time bellow 61" {
				go func() {
					time.Sleep(10 * time.Second)
					err := s.ChannelMessageDelete(m.ChannelID, m.ID)
					if err != nil {
						return
					}
				}()
			}

			roles, err := s.GuildRoles(m.GuildID)
			if err != nil {
				return
			}
			var mentionRole string
			for _, role := range roles {
				if role.Name == "SP Notifications" {
					mentionRole = role.Mention()
				}
			}

			if m.Content == mentionRole {
				err := s.MessageReactionAdd(m.ChannelID, m.Message.ID, "won:"+svc.GetEmojiByName(context.Background(), m.GuildID, "won"))
				if err != nil {
					fmt.Println(err)
					return
				}
				err = s.MessageReactionAdd(m.ChannelID, m.Message.ID, "lost:"+svc.GetEmojiByName(context.Background(), m.GuildID, "lost"))
				if err != nil {
					fmt.Println(err)
					return
				}
				err = s.MessageReactionAdd(m.ChannelID, m.Message.ID, "dislike:"+svc.GetEmojiByName(context.Background(), m.GuildID, "dislike"))
				if err != nil {
					fmt.Println(err)
					return
				}
				t, _ := strconv.Atoi(strings.Split(m.Message.Embeds[0].Fields[1].Value, " ")[0])

				go func() {
					for t != 0 {
						time.Sleep(1 * time.Minute)
						t--
						embed := &discordgo.MessageEmbed{
							Author: &discordgo.MessageEmbedAuthor{},
							Color:  m.Embeds[0].Color,
							Title:  m.Embeds[0].Title,
							Fields: []*discordgo.MessageEmbedField{
								m.Embeds[0].Fields[0],
								{
									Name:   "Time remaining: ",
									Value:  strconv.Itoa(t) + " minutes",
									Inline: true,
								},
							},
							Thumbnail: m.Embeds[0].Thumbnail,
							Footer:    m.Embeds[0].Footer,
						}

						_, err = s.ChannelMessageEditEmbed(m.ChannelID, m.ID, embed)
						if err != nil {
							return
						}
					}
					mapName := m.Embeds[0].Fields[0].Value
					err = s.ChannelMessageDelete(m.ChannelID, m.ID)
					if err != nil {
						return
					}
					logMessage, err := s.ChannelMessage(logSpChannelID, logSpMessageID)
					if err != nil {
						return
					}

					var winningNationShort string
					var winningNationLong string
					if m.Embeds[0].Color == 0x00FFFF {
						winningNationShort = "ani"
						winningNationLong = "Arlington National Influence"
					} else {
						winningNationShort = "bcu"
						winningNationLong = "Bygeniou City United"
					}

					embed := &discordgo.MessageEmbed{
						Author: &discordgo.MessageEmbedAuthor{},
						Color:  0x000000,
						Title:  "STRATEGIC POINT HISTORY",
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:   "Map: ",
								Value:  strings.ReplaceAll(logMessage.Embeds[0].Fields[0].Value, "-", "") + "\n" + mapName,
								Inline: true,
							},
							{
								Name:   "Spawn Time: ",
								Value:  strings.ReplaceAll(logMessage.Embeds[0].Fields[1].Value, "-", "") + "\n" + "<t:" + strconv.Itoa(int(time.Now().Add(time.Hour*time.Duration(1*-1)).Unix())) + ":R>",
								Inline: true,
							},
							{
								Name:   "Winning Nation: ",
								Value:  strings.ReplaceAll(logMessage.Embeds[0].Fields[2].Value, "-", "") + "\n" + "<:" + winningNationShort + ":" + svc.GetEmojiByName(context.Background(), m.GuildID, winningNationShort) + "> " + winningNationLong,
								Inline: true,
							},
						},
					}

					_, err = s.ChannelMessageEditEmbed(logSpChannelID, logSpMessageID, embed)
					if err != nil {
						fmt.Println(err)
						return
					}
				}()
			}
		}
	}
}

func HandleSpNotificationReactions(svc Service) func(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	return func(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
		if m.Member.User.ID != s.State.User.ID {
			if m.ChannelID == "1039910012228747375" {
				if m.Emoji.Name == "dislike" {
					err := s.ChannelMessageDelete(m.ChannelID, m.MessageID)
					if err != nil {
						fmt.Println(err)
						return
					}
				}
				if m.Emoji.Name == "won" || m.Emoji.Name == "lost" {
					message, err := s.ChannelMessage(m.ChannelID, m.MessageID)
					if err != nil {
						return
					}

					logMessage, err := s.ChannelMessage(logSpChannelID, logSpMessageID)
					if err != nil {
						return
					}

					var winningNationShort string
					var winningNationLong string
					if m.Emoji.Name == "won" {
						winningNationShort = "ani"
						winningNationLong = "Arlington National Influence"
					} else {
						winningNationShort = "bcu"
						winningNationLong = "Bygeniou City United"
					}

					value, _ := strconv.Atoi(strings.Split(message.Embeds[0].Fields[1].Value, " ")[0])
					value = 60 - value

					embed := &discordgo.MessageEmbed{
						Author: &discordgo.MessageEmbedAuthor{},
						Color:  0x000000,
						Title:  "STRATEGIC POINT HISTORY",
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:   "Map: ",
								Value:  strings.ReplaceAll(logMessage.Embeds[0].Fields[0].Value, "-", "") + "\n" + message.Embeds[0].Fields[0].Value,
								Inline: true,
							},
							{
								Name:   "Spawn Time: ",
								Value:  strings.ReplaceAll(logMessage.Embeds[0].Fields[1].Value, "-", "") + "\n" + "<t:" + strconv.Itoa(int(time.Now().Add(time.Minute*time.Duration(1*-value)).Unix())) + ":R>",
								Inline: true,
							},
							{
								Name:   "Winning Nation: ",
								Value:  strings.ReplaceAll(logMessage.Embeds[0].Fields[2].Value, "-", "") + "\n" + "<:" + winningNationShort + ":" + svc.GetEmojiByName(context.Background(), m.GuildID, winningNationShort) + "> " + winningNationLong,
								Inline: true,
							},
						},
					}

					_, err = s.ChannelMessageEditEmbed(logSpChannelID, logSpMessageID, embed)
					if err != nil {
						fmt.Println(err)
						return
					}

					err = s.ChannelMessageDelete(m.ChannelID, m.MessageID)
					if err != nil {
						fmt.Println(err)
						return
					}
				}
			}
		}
	}
}

func HandleSpLogMessage(svc Service) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
					if i.ApplicationCommandData().Name == "history-sp" {
						if !spHistory {
							go InitResetSpLog(s)
							spHistory = true
						}

						embed := &discordgo.MessageEmbed{
							Author: &discordgo.MessageEmbedAuthor{},
							Color:  0x000000,
							Title:  "STRATEGIC POINT HISTORY",
							Fields: []*discordgo.MessageEmbedField{
								{
									Name:   "Map: ",
									Value:  "----------",
									Inline: true,
								},
								{
									Name:   "Spawn Time: ",
									Value:  "----------",
									Inline: true,
								},
								{
									Name:   "Winning Nation: ",
									Value:  "----------",
									Inline: true,
								},
							},
						}
						em, err := s.ChannelMessageSendEmbed(i.ChannelID, embed)
						if err != nil {
							fmt.Println(err)
							return
						}

						logSpChannelID = em.ChannelID
						logSpMessageID = em.ID

						err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "Remember now this is the current SP Spawn History that will work",
								Flags:   discordgo.MessageFlagsEphemeral,
							},
						})

					}
				}
			}
		}
	}
}

func HandleServerStatus() func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		perms, err := s.UserChannelPermissions(i.Member.User.ID, i.ChannelID)
		if err != nil {
			fmt.Println(err)
			return
		}

		if perms&discordgo.PermissionAdministrator != 0 {
			channel, err := s.Channel(i.ChannelID)
			if err != nil {
				return
			}

			switch i.Type {
			case discordgo.InteractionApplicationCommand:
				{
					if i.ApplicationCommandData().Name == "server-online" {
						serverStatusChannelID = i.ChannelID
						editChannel := discordgo.ChannelEdit{
							Name:     " 🟢┃game-info",
							Position: channel.Position,
						}
						_, err := s.ChannelEdit(i.ChannelID, &editChannel)
						if err != nil {
							fmt.Println(err)
							return
						}
						err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "You changed channel name to 🟢┃game-info",
								Flags:   discordgo.MessageFlagsEphemeral,
							},
						})
						if err != nil {
							fmt.Println(err)
							return
						}

						message, err := s.ChannelMessageSend(serverStatusChannelID, "<t:"+strconv.Itoa(int(time.Now().Unix()))+":R> 🟢 **"+i.ApplicationCommandData().Options[0].StringValue()+"**")
						if err != nil {
							fmt.Println(err)
							return
						}
						serverOnlineMessage.ID = message.ID
						serverOnlineMessage.Content = message.Content

						if serverOfflineMessage.ID != "" {
							if serverOfflineMessage.Content[0] != '~' {
								s.ChannelMessageEdit(serverStatusChannelID, serverOfflineMessage.ID, "~~"+serverOfflineMessage.Content+"~~")
							}
						}
						if serverMaintMessage.ID != "" {
							if serverMaintMessage.Content[0] != '~' {
								s.ChannelMessageEdit(serverStatusChannelID, serverMaintMessage.ID, "~~"+serverMaintMessage.Content+"~~")
							}
						}
					}
					if i.ApplicationCommandData().Name == "server-offline" {
						serverStatusChannelID = i.ChannelID
						editChannel := discordgo.ChannelEdit{
							Name:     "🔴┃game-info❕",
							Position: channel.Position,
						}
						_, err := s.ChannelEdit(i.ChannelID, &editChannel)
						if err != nil {
							fmt.Println(err)
							return
						}
						err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "You changed channel name to 🔴┃game-info❕",
								Flags:   discordgo.MessageFlagsEphemeral,
							},
						})
						if err != nil {
							fmt.Println(err)
							return
						}

						message, err := s.ChannelMessageSend(serverStatusChannelID, "<t:"+strconv.Itoa(int(time.Now().Unix()))+":R>🔴 **"+i.ApplicationCommandData().Options[0].StringValue()+"**")
						if err != nil {
							fmt.Println(err)
							return
						}
						serverOfflineMessage.ID = message.ID
						serverOfflineMessage.Content = message.Content

						if serverOnlineMessage.ID != "" {
							if serverOnlineMessage.Content[0] != '~' {
								s.ChannelMessageEdit(serverStatusChannelID, serverOnlineMessage.ID, "~~"+serverOnlineMessage.Content+"~~")
							}
						}
						if serverMaintMessage.ID != "" {
							if serverMaintMessage.Content[0] != '~' {
								s.ChannelMessageEdit(serverStatusChannelID, serverMaintMessage.ID, "~~"+serverMaintMessage.Content+"~~")
							}
						}

					}
					if i.ApplicationCommandData().Name == "server-maint" {
						serverStatusChannelID = i.ChannelID
						editChannel := discordgo.ChannelEdit{
							Name:     "🟠┃game-info❕",
							Position: channel.Position,
						}
						_, err := s.ChannelEdit(i.ChannelID, &editChannel)
						if err != nil {
							fmt.Println(err)
							return
						}
						err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "You changed channel name to 🟠┃game-info❕",
								Flags:   discordgo.MessageFlagsEphemeral,
							},
						})
						if err != nil {
							fmt.Println(err)
							return
						}

						message, err := s.ChannelMessageSend(serverStatusChannelID, "<t:"+strconv.Itoa(int(time.Now().Unix()))+":R> 🟠 **"+i.ApplicationCommandData().Options[0].StringValue()+"**")
						if err != nil {
							fmt.Println(err)
							return
						}
						serverMaintMessage.ID = message.ID
						serverMaintMessage.Content = message.Content

						if serverOfflineMessage.ID != "" {
							if serverOfflineMessage.Content[0] != '~' {
								s.ChannelMessageEdit(serverStatusChannelID, serverOfflineMessage.ID, "~~"+serverOfflineMessage.Content+"~~")
							}
						}
						if serverOnlineMessage.ID != "" {
							if serverOnlineMessage.Content[0] != '~' {
								s.ChannelMessageEdit(serverStatusChannelID, serverOnlineMessage.ID, "~~"+serverOnlineMessage.Content+"~~")
							}
						}
					}
				}
			}
		}
	}
}

func Register(bot *discordgo.Session, svc Service) {
	bot.AddHandler(SP(svc))
	bot.AddHandler(HandleSpNotification(svc))
	bot.AddHandler(HandleSpNotificationReactions(svc))
	bot.AddHandler(HandleSpLogMessage(svc))
	bot.AddHandler(HandleServerStatus())
}
