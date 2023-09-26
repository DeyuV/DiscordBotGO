package strategicpoint

import "time"

type UploadCareAPIFilesResponse struct {
	Results []struct {
		OriginalFileUrl  string `json:"original_file_url"`
		OriginalFilename string `json:"original_filename"`
	} `json:"results"`
}

type SPLogs struct {
	ID              string
	GuildID         string
	MapName         string
	SpawnTime       string
	WinningNation   string
	UserSpawning    string
	UserInteracting string
	SPDate          *time.Time
}

type SPLogsRefresh struct {
	ID              string
	MapName         string
	SpawnTime       string
	WinningNation   string
	UserSpawning    string
	UserInteracting string
}
