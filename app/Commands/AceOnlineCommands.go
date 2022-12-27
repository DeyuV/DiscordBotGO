package Commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	ANImaps         = []string{"Edmont Valley", "Desert of Ardor", "Crystal Cave", "Plain of Doleful Melody", "Herremeze Relic Site", "Atus Beach", "Gjert Road", "Slope Port", "Portsmouth Canyon"}
	BCUmaps         = []string{"Bach Mountain Chain", "Blackburn Site", "Zaylope Beach", "Starlite Valley", "Redline", "Kahlua Beach", "Nubarke Cave", "Orina Peninsula", "Daisy Riverhead"}
	aniMenuOption   []discordgo.SelectMenuOption
	bcuMenuOption   []discordgo.SelectMenuOption
	aniResponseData []discordgo.MessageComponent
	bcuResponseData []discordgo.MessageComponent
	aniMenuID       string
	bcuMenuID       string
	spChannelID     string
)

func SP(svc Service) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		perms, err := s.UserChannelPermissions(i.Member.User.ID, i.ChannelID)
		if err != nil {
			fmt.Println(err)
			return
		}

		if perms&discordgo.PermissionAdministrator != 0 {
			if aniMenuOption == nil {
				for _, m := range ANImaps {
					aniMenuOption = append(aniMenuOption, discordgo.SelectMenuOption{
						Label:       m,
						Value:       m,
						Description: "strategic point",
						Emoji: discordgo.ComponentEmoji{
							Name:     strings.ReplaceAll(m, " ", ""),
							ID:       svc.GetEmojiByName(context.Background(), i.GuildID, strings.ReplaceAll(m, " ", "")),
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
				for _, m := range BCUmaps {
					bcuMenuOption = append(bcuMenuOption, discordgo.SelectMenuOption{
						Label:       m,
						Value:       m,
						Description: "strategic point",
						Emoji: discordgo.ComponentEmoji{
							Name:     strings.ReplaceAll(m, " ", ""),
							ID:       svc.GetEmojiByName(context.Background(), i.GuildID, strings.ReplaceAll(m, " ", "")),
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
					if i.ApplicationCommandData().Name == "ani-sp" {
						spChannelID = i.ChannelID
						err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								CustomID: "ani-sp",
								Content:  "You used command to spawn ANI menu",
								Flags:    discordgo.MessageFlagsEphemeral,
							},
						})
						if err != nil {
							fmt.Println(err)
							return
						}

						messageComplex, err := s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
							Components: aniResponseData,
						})
						if err != nil {
							fmt.Println("Failed to create ANI menu")
							fmt.Println(err)
						}
						aniMenuID = messageComplex.ID
					}

					if i.ApplicationCommandData().Name == "bcu-sp" {
						spChannelID = i.ChannelID
						err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								CustomID: "bcu-sp",
								Content:  "You used command to spawn BCU menu",
								Flags:    discordgo.MessageFlagsEphemeral,
							},
						})
						if err != nil {
							fmt.Println(err)
							return
						}
						messageComplex, err := s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
							Components: bcuResponseData,
						})
						if err != nil {
							fmt.Println("Failed to create BCU menu")
							fmt.Println(err)
						}
						bcuMenuID = messageComplex.ID
					}
				}
			case discordgo.InteractionMessageComponent:
				{
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
								Content: "Please insert a number bellow 61",
							},
						})
						return
					}

					/* if t > 60 {
						err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "Please insert a number bellow 61",
							},
						})
						if err != nil {
							fmt.Println(err)
						}
						return
					} */

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

					err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							CustomID: "SP Notification",
							Content:  mentionRole,
							Embeds:   []*discordgo.MessageEmbed{embed},
						},
					})

					if err != nil {
						fmt.Println("Failed modal submit")
						fmt.Println(err)
					}
				}
			}
		}
	}
}
