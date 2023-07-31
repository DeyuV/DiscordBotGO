package strategicpoint

import (
	"DiscordBotGO/pkg/config"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Repository interface {
	GetEmojiByName(ctx context.Context, guildId, emojiName string) string
	GetChannelIdByNameAndGuildID(ctx context.Context, guildId, name string) (string, error)
	UpdateMessageId(ctx context.Context, guildId, name, messageId string) error
	GetMessageIdByNameAndGuildID(ctx context.Context, guildId, name string) (string, error)
	AddMessageId(ctx context.Context, guildId, name, messageId string) error
	DeleteMessageId(ctx context.Context, guildId, messageId string) error
	GetGuildId(ctx context.Context, id int) (string, error)

	AddSP(ctx context.Context, guildId, mapName, spawnTime, winningNation, userSpawning, userInteracting string) (int, error)
	DeleteSP(ctx context.Context, id int) error
	UpdateSPmap(ctx context.Context, id int, mapName string) error
	UpdateSPspawnTime(ctx context.Context, id int, spawnTime string) error
	UpdateSPwinningNation(ctx context.Context, id int, winningNation string) error
	UpdateSPmodified(ctx context.Context, id int, modified string) error
	GetAllSPLogsByGuild(ctx context.Context, guildId string) ([]SPLogs, error)
	GetSPbyGuildAndId(ctx context.Context, guildId string, spId int) error

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

func (s *serviceImplementation) UpdateLog(ctx context.Context, id int, guildId, mapName, spawnTime, winningNation, userModify string) error {
	spGuildId, err := s.repo.GetGuildId(ctx, id)
	if err != nil {
		return err
	}

	if spGuildId != guildId {
		return errors.New("wrong sp id")
	}

	if mapName != "?" {
		err = s.repo.UpdateSPmap(ctx, id, mapName)
		if err != nil {
			return err
		}
	}

	if spawnTime != "?" {
		err = s.repo.UpdateSPspawnTime(ctx, id, spawnTime)
		if err != nil {
			return err
		}
	}

	if winningNation != "?" {
		err = s.repo.UpdateSPwinningNation(ctx, id, winningNation)
		if err != nil {
			return err
		}
	}

	err = s.repo.UpdateSPmodified(ctx, id, userModify)
	if err != nil {
		return err
	}

	return nil
}

func (s *serviceImplementation) RefreshLog(ctx context.Context, guildId string) (*discordgo.MessageEmbed, error) {
	spLogs, err := s.repo.GetAllSPLogsByGuild(ctx, guildId)
	if err != nil {
		return nil, err
	}

	var concatSpLogs SPLogsRefresh
	for _, sp := range spLogs {
		if sp.SPDate.Day() == time.Now().Day() && sp.SPDate.Month() == time.Now().Month() {
			concatSpLogs.ID += strconv.Itoa(sp.ID) + "\n"
			concatSpLogs.MapName += sp.MapName + "\n"
			concatSpLogs.SpawnTime += sp.SpawnTime + "\n"

			if sp.WinningNation == "ani" {
				concatSpLogs.WinningNation += "<:" + sp.WinningNation + ":" + s.GetEmojiByName(context.Background(), guildId, sp.WinningNation) + ">" + config.ANIlongName + "\n"
			}

			if sp.WinningNation == "bcu" {
				concatSpLogs.WinningNation += "<:" + sp.WinningNation + ":" + s.GetEmojiByName(context.Background(), guildId, sp.WinningNation) + ">" + config.BCUlongName + "\n"
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

	membersLogSpChannelID, err := s.GetChannelIdByNameAndGuildID(context.Background(), guildId, config.LogStrategicpoint)
	if err != nil {
		return err
	}

	membersLogSpMessageID, err := s.GetMessageIdByNameAndGuildID(context.Background(), guildId, config.LogStrategicpoint)
	if err != nil {
		return err
	}

	_, err = session.ChannelMessageEditEmbed(membersLogSpChannelID, membersLogSpMessageID, membersEmbed)
	if err != nil {
		return err
	}

	return nil
}

func (s *serviceImplementation) AddSPtoLog(ctx context.Context, guildId, mapName, spawnTime, winningNation, userSpawning, userInteracting string) (int, error) {
	return s.repo.AddSP(ctx, guildId, mapName, spawnTime, winningNation, userSpawning, userInteracting)
}

func (s *serviceImplementation) DeleteSPfromLog(ctx context.Context, id int) error {
	return s.repo.DeleteSP(ctx, id)
}

func (s *serviceImplementation) UpdateMessageId(ctx context.Context, guildId, name, messageId string) error {
	return s.repo.UpdateMessageId(ctx, guildId, name, messageId)
}

func (s *serviceImplementation) GetMessageIdByNameAndGuildID(ctx context.Context, guildId, name string) (string, error) {
	return s.repo.GetMessageIdByNameAndGuildID(ctx, guildId, name)
}

func (s *serviceImplementation) GetEmojiByName(ctx context.Context, guildId, emojiName string) string {
	return s.repo.GetEmojiByName(ctx, guildId, emojiName)
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

func (s *serviceImplementation) VerifySpId(ctx context.Context, guildId string, spId int) error {
	return s.repo.GetSPbyGuildAndId(ctx, guildId, spId)
}

func (s *serviceImplementation) UpdateChannelId(ctx context.Context, guildId, name, channelId string) error {
	return s.repo.UpdateChannelId(ctx, guildId, name, channelId)
}

func (s *serviceImplementation) AddChannelId(ctx context.Context, guildId, name, channelId string) error {
	return s.repo.AddChannelId(ctx, guildId, name, channelId)
}
