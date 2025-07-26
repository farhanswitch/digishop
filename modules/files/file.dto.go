package files

type fileData struct {
	ID               string `json:"id"`
	FileName         string `json:"file_name"`
	OriginalFileName string `json:"original_file_name"`
	FileExt          string `json:"file_ext"`
	FileMimeType     string `json:"file_mime_type"`
	FilePath         string `json:"file_path"`
}
