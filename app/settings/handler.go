package settings

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

type Service interface {
	UpdateChannelId(ctx context.Context, guildId, name, channelId string) error
	AddChannelId(ctx context.Context, guildId, name, channelId string) error
}

/* no use for now but may be needed later
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

					}
				}
			}
		}
	}
} */

func Register(bot *discordgo.Session, svc Service) {
	//bot.AddHandler(SetChannelId(svc))
}
