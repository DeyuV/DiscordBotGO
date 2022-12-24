package Commands

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type Repository interface {
	GetEmojiByName(ctx context.Context, guildId, emojiName string) string
}

func NewService(repo Repository) Service {
	return &serviceImplementation{
		repo: repo,
	}
}

type serviceImplementation struct {
	repo Repository
}

func (s *serviceImplementation) GetEmojiByName(ctx context.Context, guildId, emojiName string) string {
	return s.repo.GetEmojiByName(ctx, guildId, emojiName)
}

func (s *serviceImplementation) ToggleCommand(ctx context.Context, name string) {
	//TODO implement me
	panic("implement me")
}

func (s *serviceImplementation) GetImageURL(mapName string) string {
	requestURL := "https://api.uploadcare.com/files/"
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
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
