package Commands

type UploadCareAPIFilesResponse struct {
	Results []struct {
		OriginalFileUrl  string `json:"original_file_url"`
		OriginalFilename string `json:"original_filename"`
	} `json:"results"`
}
