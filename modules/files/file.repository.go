package files

import (
	"digishop/connections"
	"net/http"

	custom_errors "digishop/utilities/errors"
)

type fileRepo struct{}

var repo fileRepo

func (f fileRepo) CreateFile(file fileData) (bool, custom_errors.CustomError) {
	_, err := connections.DbMySQL().Exec("INSERT INTO files(id, original_filename, filename, path, extension, mime_type) VALUES(?,?,?,?,?,?)", file.ID, file.OriginalFileName, file.FileName, file.FilePath, file.FileExt, file.FileMimeType)
	if err != nil {
		customErr := custom_errors.CustomError{
			Code:          http.StatusInternalServerError,
			MessageToSend: "Internal Server Error",
			Message:       err.Error(),
		}
		return true, customErr
	}
	return false, custom_errors.CustomError{}
}

func factoryFileRepo() fileRepo {
	if repo == (fileRepo{}) {
		repo = fileRepo{}
	}
	return repo
}
