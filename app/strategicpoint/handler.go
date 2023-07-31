package strategicpoint

import (
	"DiscordBotGO/pkg/config"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Service interface {
	GetEmojiByName(ctx context.Context, guildId, emojiName string) string
	GetChannelIdByNameAndGuildID(ctx context.Context, guildId, name string) (string, error)
	UpdateMessageId(ctx context.Context, guildId, name, messageId string) error
	GetMessageIdByNameAndGuildID(ctx context.Context, guildId, name string) (string, error)
	AddMessageId(ctx context.Context, guildId, name, messageId string) error
	DeleteMessageId(ctx context.Context, guildId, messageId string) error // useless for now
	VerifySpId(ctx context.Context, guildId string, spId int) error

	GetImageURL(name string) string

	LogEmbed(mapValue, timeValue, nationValue string) *discordgo.MessageEmbed
	InitResetLog(session *discordgo.Session)
	UpdateLog(ctx context.Context, id int, guildId, mapName, spawnTime, winningNation, userModify string) error
	AddSPtoLog(ctx context.Context, guildId, mapName, spawnTime, winningNation, userSpawning, userInteracting string) (int, error)
	DeleteSPfromLog(ctx context.Context, id int) error
	EditeEmbeds(ctx context.Context, session *discordgo.Session, guildId string, empty bool) error

	UpdateChannelId(ctx context.Context, guildId, name, channelId string) error
	AddChannelId(ctx context.Context, guildId, name, channelId string) error
}

var (
	ANImaps          = map[string]string{"ev": "Edmont Valley", "doa": "Desert of Ardor", "cc": "Crystal Cave", "pdm": "Plain of Doleful Melody", "hrs": "Herremeze Relic Site", "ab": "Atus Beach", "gr": "Gjert Road", "sp": "Slope Port", "pmc": "Portsmouth Canyon"}
	BCUmaps          = map[string]string{"bmc": "Bach Mountain Chain", "bs": "Blackburn Site", "zb": "Zaylope Beach", "sv": "Starlite Valley", "rl": "Redline", "kb": "Kahlua Beach", "nc": "Nubarke Cave", "op": "Orina Peninsula", "dr": "Daisy Riverhead"}
	sortedANImapKeys = []string{"ev", "doa", "cc", "pdm", "hrs", "ab", "gr", "sp", "pmc"}
	sortedBCUmapKeys = []string{"bmc", "bs", "zb", "sv", "rl", "kb", "nc", "op", "dr"}
	aniMenuOption    []discordgo.SelectMenuOption
	bcuMenuOption    []discordgo.SelectMenuOption
	aniResponseData  []discordgo.MessageComponent
	bcuResponseData  []discordgo.MessageComponent
	spHistory        = false
)

func SP(svc Service) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		perms, err := s.UserChannelPermissions(i.Member.User.ID, i.ChannelID)
		if err != nil {
			fmt.Println(err)
			return
		}

		if aniMenuOption == nil {
			for _, m := range sortedANImapKeys {
				aniMenuOption = append(aniMenuOption, discordgo.SelectMenuOption{
					Label:       ANImaps[m],
					Value:       ANImaps[m],
					Description: "strategic point",
					Emoji: discordgo.ComponentEmoji{
						Name:     strings.ReplaceAll(ANImaps[m], " ", ""),
						ID:       svc.GetEmojiByName(context.Background(), i.GuildID, strings.ReplaceAll(ANImaps[m], " ", "")),
						Animated: false,
					},
					Default: false,
				})
			}

			aniResponseData = []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							CustomID:    "ani-sp",
							Placeholder: "Select ANI strategic point",
							Options:     aniMenuOption,
						},
					},
				},
			}
		}

		if bcuMenuOption == nil {
			for _, m := range sortedBCUmapKeys {
				bcuMenuOption = append(bcuMenuOption, discordgo.SelectMenuOption{
					Label:       BCUmaps[m],
					Value:       BCUmaps[m],
					Description: "strategic point",
					Emoji: discordgo.ComponentEmoji{
						Name:     strings.ReplaceAll(BCUmaps[m], " ", ""),
						ID:       svc.GetEmojiByName(context.Background(), i.GuildID, strings.ReplaceAll(BCUmaps[m], " ", "")),
						Animated: false,
					},
					Default: false,
				})
			}

			bcuResponseData = []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							CustomID:    "bcu-sp",
							Placeholder: "Select BCU strategic point",
							Options:     bcuMenuOption,
						},
					},
				},
			}
		}

		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			{
				if perms&discordgo.PermissionAdministrator != 0 {
					if i.ApplicationCommandData().Name == "setup-sp" {
						err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								CustomID: "setup-sp",
								Content:  "Command to setup SP notification used",
								Flags:    discordgo.MessageFlagsEphemeral,
							},
						})
						if err != nil {
							fmt.Println(err)
							return
						}

						err = svc.AddChannelId(context.Background(), i.GuildID, config.Strategicpoint, i.ChannelID)
						if err != nil {
							err = svc.UpdateChannelId(context.Background(), i.GuildID, config.Strategicpoint, i.ChannelID)
							if err != nil {
								fmt.Println(err)
								return
							}
						}

						if !spHistory {
							go svc.InitResetLog(s)
							spHistory = true
						}

						embed := svc.LogEmbed(config.EmptyEmbedFieldValue, config.EmptyEmbedFieldValue, config.EmptyEmbedFieldValue)

						em, err := s.ChannelMessageSendEmbed(i.ChannelID, embed)
						if err != nil {
							fmt.Println(err)
							return
						}

						err = svc.AddMessageId(context.Background(), i.GuildID, config.LogStrategicpoint, em.ID)
						if err != nil {
							err = svc.UpdateMessageId(context.Background(), i.GuildID, config.LogStrategicpoint, em.ID)
							if err != nil {
								fmt.Println(err)
								return
							}
						}

						messageComplex, err := s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
							Components: aniResponseData,
						})
						if err != nil {
							fmt.Println("Failed to create ANI menu")
							fmt.Println(err)
						}

						err = svc.AddMessageId(context.Background(), i.GuildID, config.ANIMenu, messageComplex.ID)
						if err != nil {
							err = svc.UpdateMessageId(context.Background(), i.GuildID, config.ANIMenu, messageComplex.ID)
							if err != nil {
								fmt.Println(err)
								return
							}

						}

						messageComplex, err = s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
							Components: bcuResponseData,
						})
						if err != nil {
							fmt.Println("Failed to create BCU menu")
							fmt.Println(err)
						}
						err = svc.AddMessageId(context.Background(), i.GuildID, config.BCUMenu, messageComplex.ID)
						if err != nil {
							err = svc.UpdateMessageId(context.Background(), i.GuildID, config.BCUMenu, messageComplex.ID)
							if err != nil {
								fmt.Println(err)
								return
							}

						}
					}
				}
			}
		case discordgo.InteractionMessageComponent:
			{
				aniMenuID, err := svc.GetMessageIdByNameAndGuildID(context.Background(), i.GuildID, config.ANIMenu)
				if err != nil {
					fmt.Println(err)
				}

				bcuMenuID, err := svc.GetMessageIdByNameAndGuildID(context.Background(), i.GuildID, config.BCUMenu)
				if err != nil {
					fmt.Println(err)
				}

				if i.Message.ID == aniMenuID || i.Message.ID == bcuMenuID {
					err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseModal,
						Data: &discordgo.InteractionResponseData{
							CustomID: i.MessageComponentData().Values[0],
							Title:    "SP Time",
							Components: []discordgo.MessageComponent{
								discordgo.ActionsRow{
									Components: []discordgo.MessageComponent{
										discordgo.TextInput{
											CustomID:    "TIME",
											Label:       "TIME",
											Style:       discordgo.TextInputShort,
											Placeholder: "Insert time",
											Required:    true,
											MaxLength:   2,
											MinLength:   1,
										},
									},
								},
								discordgo.ActionsRow{
									Components: []discordgo.MessageComponent{
										discordgo.TextInput{
											CustomID:  "MAP",
											Label:     "MAP",
											Style:     discordgo.TextInputShort,
											Value:     i.MessageComponentData().Values[0],
											Required:  false,
											MaxLength: 2000,
										},
									},
								},
							},
						},
					})

					if err != nil {
						fmt.Println(err)
						return
					}
				}

				// Reset interaction menu to default state
				if i.Message.ID == aniMenuID {
					_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Components: &aniResponseData})
					if err != nil {
						return
					}
				}
				if i.Message.ID == bcuMenuID {
					_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Components: &bcuResponseData})
					if err != nil {
						return
					}
				}
			}
		case discordgo.InteractionModalSubmit:
			{
				t, err := strconv.Atoi(i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value)
				if err != nil || t > 60 {
					err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: i.Member.Mention() + " Please insert a number below 61",
							Flags:   discordgo.MessageFlagsEphemeral,
						},
					})
					if err != nil {
						fmt.Println(err)
					}

					return
				}

				ani := false
				for _, m := range ANImaps {
					if m == i.ModalSubmitData().CustomID {
						ani = true
					}
				}

				var color int
				if ani {
					color = 0x00FFFF
				} else {
					color = 0xFFA500
				}

				embed := &discordgo.MessageEmbed{
					Author: &discordgo.MessageEmbedAuthor{},
					Color:  color,
					Title:  "A strategic point has been created!",
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "Map: ",
							Value:  i.ModalSubmitData().CustomID,
							Inline: true,
						},
						{
							Name:   "Time remaining: ",
							Value:  i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value + " minutes",
							Inline: true,
						},
					},
					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: svc.GetImageURL(i.ModalSubmitData().CustomID),
					},
				}

				roles, err := s.GuildRoles(i.GuildID)
				if err != nil {
					return
				}

				var mentionRole string
				for _, role := range roles {
					if role.Name == "SP Notifications" {
						mentionRole = role.Mention()
					}
				}
				wonEmoji := discordgo.Emoji{
					ID:       svc.GetEmojiByName(context.Background(), i.GuildID, "won"),
					Name:     "won",
					Animated: false,
				}

				lostEmoji := discordgo.Emoji{
					ID:       svc.GetEmojiByName(context.Background(), i.GuildID, "lost"),
					Name:     "lost",
					Animated: false,
				}

				dislikeEmoji := discordgo.Emoji{
					ID:       svc.GetEmojiByName(context.Background(), i.GuildID, "dislike"),
					Name:     "dislike",
					Animated: false,
				}

				spChannelID, err := svc.GetChannelIdByNameAndGuildID(context.Background(), i.GuildID, config.Strategicpoint)
				if err != nil {
					fmt.Println(err)
					return
				}

				_, err = s.ChannelMessageSendComplex(spChannelID, &discordgo.MessageSend{
					Content: mentionRole,
					Embed:   embed,
				})

				if err != nil {
					fmt.Println(err)
					return
				}

				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "**Click the following icons once the SP has finished:**\n" +
							wonEmoji.MessageFormat() + " - We have won this SP\n" +
							lostEmoji.MessageFormat() + " - We have lost this SP\n" +
							dislikeEmoji.MessageFormat() + " - Cancel SP (mistakes only)",
						Flags: discordgo.MessageFlagsEphemeral,
					},
				})

				if err != nil {
					fmt.Println("Failed modal submit")
					fmt.Println(err)
				}

				spForumChannelId, err := svc.GetChannelIdByNameAndGuildID(context.Background(), i.GuildID, config.SPforum)
				if err != nil {
					fmt.Println(err)
					return
				}

				_, err = s.ChannelMessageSendEmbed(spForumChannelId, &discordgo.MessageEmbed{
					Title:       "A strategic point has been created!",
					Description: "Map: " + i.ModalSubmitData().CustomID,
					Image:       &discordgo.MessageEmbedImage{URL: svc.GetImageURL(i.ModalSubmitData().CustomID)}})

				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}

func Notification(svc Service) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
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
					mapName := m.Embeds[0].Fields[0].Value
					var spawningUser string

					for t != 0 {
						time.Sleep(1 * time.Second)
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
							fmt.Println(err)
							return
						}
					}

					err = s.ChannelMessageDelete(m.ChannelID, m.ID)
					if err != nil {
						return
					}

					var winningNationShort string
					if m.Embeds[0].Color == 0x00FFFF {
						winningNationShort = config.ANIshortName
					} else {
						winningNationShort = config.BCUshortName
					}

					_, err = svc.AddSPtoLog(context.Background(), m.GuildID, mapName, "<t:"+strconv.Itoa(int(time.Now().Add(time.Hour*time.Duration(1*-1)).Unix()))+":R>", winningNationShort, spawningUser, m.Author.Username)
					if err != nil {
						fmt.Println(err)
						return
					}
					err = svc.EditeEmbeds(context.Background(), s, m.GuildID, false)
					if err != nil {
						fmt.Println(err)
						return
					}
				}()
			}
		}
	}
}

func Reactions(svc Service) func(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	return func(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
		if m.Member.User.ID != s.State.User.ID {
			message, err := s.ChannelMessage(m.ChannelID, m.MessageID)
			if err != nil {
				fmt.Println(err)
				return
			}

			ok := false
			for _, reaction := range message.Reactions {
				if reaction.Me {
					ok = true
					break
				}
			}

			if ok {
				if m.Emoji.Name == "dislike" {
					err = s.ChannelMessageDelete(m.ChannelID, m.MessageID)
					if err != nil {
						fmt.Println(err)
						return
					}
					return
				}
				if m.Emoji.Name == "won" || m.Emoji.Name == "lost" {
					var winningNationShort string
					if m.Emoji.Name == "won" {
						winningNationShort = config.ANIshortName
					} else {
						winningNationShort = config.BCUshortName
					}

					value, _ := strconv.Atoi(strings.Split(message.Embeds[0].Fields[1].Value, " ")[0])
					value = 60 - value

					_, err = svc.AddSPtoLog(context.Background(), m.GuildID, message.Embeds[0].Fields[0].Value,
						"<t:"+strconv.Itoa(int(time.Now().Add(time.Minute*time.Duration(value*-1)).Unix()))+":R>", winningNationShort, "?", m.Member.User.Username)
					if err != nil {
						fmt.Println(err)
						return
					}

					err = svc.EditeEmbeds(context.Background(), s, m.GuildID, false)
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

func Register(bot *discordgo.Session, svc Service) {
	// SP menu + interaction response
	bot.AddHandler(SP(svc))
	bot.AddHandler(Notification(svc))
	bot.AddHandler(Reactions(svc))
}
