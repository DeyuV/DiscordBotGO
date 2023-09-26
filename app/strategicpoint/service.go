package strategicpoint

import (
	"DiscordBotGO/pkg/aceonline"
	"DiscordBotGO/pkg/config"
	"DiscordBotGO/pkg/emoji"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Repository interface {
	GetChannelIdByNameAndGuildID(ctx context.Context, guildId, name string) (string, error)
	UpdateMessageId(ctx context.Context, guildId, name, messageId string) error
	GetMessageIdByNameAndGuildID(ctx context.Context, guildId, name string) (string, error)
	AddMessageId(ctx context.Context, guildId, name, messageId string) error
	DeleteMessageId(ctx context.Context, guildId, messageId string) error
	GetGuildId(ctx context.Context, id int) (string, error)

	AddSP(ctx context.Context, id, guildId, userSpawning string) error
	DeleteSP(ctx context.Context, id string) error
	UpdateSP(ctx context.Context, id, mapName, spawntime, winningNation, userInteracting string) error
	GetAllSPLogsByGuild(ctx context.Context, guildId string) ([]SPLogs, error)
	GetSPbyGuildAndId(ctx context.Context, guildId, spId string) error

	UpdateChannelId(ctx context.Context, guildId, name, channelId string) error
	AddChannelId(ctx context.Context, guildId, name, channelId string) error
}

func NewService(repo Repository) Service {
	return &serviceImplementation{
		repo: repo,
	}
}

type serviceImplementation struct {
	repo Repository
}

func (s *serviceImplementation) LogEmbed(mapValue, timeValue, nationValue string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  0x000000,
		Title:  "STRATEGIC POINT HISTORY",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Map: ",
				Value:  mapValue,
				Inline: true,
			},
			{
				Name:   "Spawn Time: ",
				Value:  timeValue,
				Inline: true,
			},
			{
				Name:   "Winning Nation: ",
				Value:  nationValue,
				Inline: true,
			},
		},
	}
}

func (s *serviceImplementation) InitResetLog(session *discordgo.Session) {
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
		for _, guild := range session.State.Guilds {
			s.EditeEmbeds(context.Background(), session, guild.ID, true)
		}
	}
}

func (s *serviceImplementation) GetImageURL(mapName string) string {
	req, err := http.NewRequest(http.MethodGet, config.UploadCareRequestURL, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/vnd.uploadcare-v0.7+json")
	req.Header.Add("Authorization", "Uploadcare.Simple "+os.Getenv("PUBLICAPIKEY")+":"+os.Getenv("PRIVATEAPIKEY"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	var result UploadCareAPIFilesResponse
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		fmt.Println("failed to decode")
		return ""
	}

	for _, name := range result.Results {
		if name.OriginalFilename == strings.ReplaceAll(mapName, " ", "")+".jpeg" {
			return name.OriginalFileUrl
		}
	}
	return ""
}

func (s *serviceImplementation) RefreshLog(ctx context.Context, guildId string) (*discordgo.MessageEmbed, error) {
	spLogs, err := s.repo.GetAllSPLogsByGuild(ctx, guildId)
	if err != nil {
		return nil, err
	}

	var concatSpLogs SPLogsRefresh
	for _, sp := range spLogs {
		if sp.SPDate.Day() == time.Now().Day() && sp.SPDate.Month() == time.Now().Month() {
			concatSpLogs.MapName += sp.MapName + "\n"
			concatSpLogs.SpawnTime += sp.SpawnTime + "\n"

			if sp.WinningNation == "ani" {
				concatSpLogs.WinningNation += emoji.ANI + aceonline.ANIlongName + "\n"
			}

			if sp.WinningNation == "bcu" {
				concatSpLogs.WinningNation += emoji.BCU + aceonline.BCUlongName + "\n"
			}
		}
	}

	if concatSpLogs.MapName == "" {
		concatSpLogs.ID = config.EmptyEmbedFieldValue
		concatSpLogs.MapName = config.EmptyEmbedFieldValue
		concatSpLogs.SpawnTime = config.EmptyEmbedFieldValue
		concatSpLogs.WinningNation = config.EmptyEmbedFieldValue
	}

	membersEmbed := s.LogEmbed(concatSpLogs.MapName, concatSpLogs.SpawnTime, concatSpLogs.WinningNation)

	return membersEmbed, nil
}

func (s *serviceImplementation) EditeEmbeds(ctx context.Context, session *discordgo.Session, guildId string, empty bool) error {
	var membersEmbed *discordgo.MessageEmbed
	var err error
	if !empty {
		membersEmbed, err = s.RefreshLog(context.Background(), guildId)
		if err != nil {
			return err
		}
	} else {
		membersEmbed = s.LogEmbed(config.EmptyEmbedFieldValue, config.EmptyEmbedFieldValue, config.EmptyEmbedFieldValue)
	}

	membersLogSpChannelID, err := s.GetChannelIdByNameAndGuildID(context.Background(), guildId, aceonline.LogStrategicpoint)
	if err != nil {
		return err
	}

	membersLogSpMessageID, err := s.GetMessageIdByNameAndGuildID(context.Background(), guildId, aceonline.LogStrategicpoint)
	if err != nil {
		return err
	}

	_, err = session.ChannelMessageEditEmbed(membersLogSpChannelID, membersLogSpMessageID, membersEmbed)
	if err != nil {
		return err
	}

	return nil
}

func (s *serviceImplementation) AddSP(ctx context.Context, id, guildId, userSpawning string) error {
	return s.repo.AddSP(ctx, id, guildId, userSpawning)
}

func (s *serviceImplementation) UpdateSP(ctx context.Context, id, mapName, spawntime, winningNation, userInteracting string) error {
	return s.repo.UpdateSP(ctx, id, mapName, spawntime, winningNation, userInteracting)
}

func (s *serviceImplementation) DeleteSPfromLog(ctx context.Context, id string) error {
	return s.repo.DeleteSP(ctx, id)
}

func (s *serviceImplementation) UpdateMessageId(ctx context.Context, guildId, name, messageId string) error {
	return s.repo.UpdateMessageId(ctx, guildId, name, messageId)
}

func (s *serviceImplementation) GetMessageIdByNameAndGuildID(ctx context.Context, guildId, name string) (string, error) {
	return s.repo.GetMessageIdByNameAndGuildID(ctx, guildId, name)
}

func (s *serviceImplementation) GetChannelIdByNameAndGuildID(ctx context.Context, guildId, name string) (string, error) {
	return s.repo.GetChannelIdByNameAndGuildID(ctx, guildId, name)
}

func (s *serviceImplementation) AddMessageId(ctx context.Context, guildId, name, messageId string) error {
	return s.repo.AddMessageId(ctx, guildId, name, messageId)
}

func (s *serviceImplementation) DeleteMessageId(ctx context.Context, guildId, messageId string) error {
	return s.repo.DeleteMessageId(ctx, guildId, messageId)
}

func (s *serviceImplementation) VerifySpId(ctx context.Context, guildId, spId string) error {
	return s.repo.GetSPbyGuildAndId(ctx, guildId, spId)
}

func (s *serviceImplementation) UpdateChannelId(ctx context.Context, guildId, name, channelId string) error {
	return s.repo.UpdateChannelId(ctx, guildId, name, channelId)
}

func (s *serviceImplementation) AddChannelId(ctx context.Context, guildId, name, channelId string) error {
	return s.repo.AddChannelId(ctx, guildId, name, channelId)
}
