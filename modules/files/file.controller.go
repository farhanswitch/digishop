package files

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"time"

	"strings"

	"github.com/google/uuid"

	"github.com/go-chi/chi/v5"
)

type fileController struct {
	repo fileRepo
}

var controller fileController

func (f fileController) UploadFileCtrl(w http.ResponseWriter, r *http.Request) {
	validMaxSize := 8 << 20
	w.Header().Set("Content-Type", "application/json")
	r.ParseMultipartForm(8 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println("Error Retrieving the File")
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"errors":"%s"}`, err.Error())
		return
	}
	defer file.Close()
	if handler.Size > int64(validMaxSize) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"errors":"File size is too large"}`)
		return
	}
	var listValidMimeType []string = []string{"image/png", "image/jpeg", "image/jpg", "image/webp", "application/octet-stream"}
	fileContentType := handler.Header.Get("Content-Type")
	var isValidMimeType bool = false
	if slices.Contains(listValidMimeType, string(fileContentType)) {
		isValidMimeType = true
	}
	if !isValidMimeType {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"errors":"Invalid file type"}`)
		return
	}
	var fileExtension string = filepath.Ext(handler.Filename)
	if fileExtension == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"errors":"Invalid file type"}`)
		return
	}
	strUUID, err := uuid.NewV7()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"errors":"%s"}`, err.Error())
		return
	}
	var newFileName string = fmt.Sprintf("%d", time.Now().Unix()) + "--" + strUUID.String() + fileExtension
	var fileNamePatternToSave string = "uploads/" + newFileName
	var dataToSave fileData = fileData{
		ID:               strUUID.String(),
		FileName:         newFileName,
		OriginalFileName: handler.Filename,
		FileExt:          fileExtension,
		FileMimeType:     fileContentType,
		FilePath:         fileNamePatternToSave,
	}

	tempFile, err := os.Create(fileNamePatternToSave)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"errors:"%s"}`, err.Error())
		return
	}
	defer tempFile.Close()
	isError, customErr := f.repo.CreateFile(dataToSave)
	if isError {
		log.Println(customErr.Message)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"errors":"%s"}`, customErr.Message)
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"errors":"%s"}`, err.Error())
		return
	}
	tempFile.Write(fileBytes)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"File uploaded successfully", "id":"%s","fileName":"%s"}`, strUUID.String(), dataToSave.FileName)
}
func (f fileController) GetFileCtrl(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")
	if strings.Contains(filename, "..") {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	// Path ke folder uploads di root project
	filePath := filepath.Join(".", "uploads", filename)

	// Buka file
	file, err := os.Open(filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer file.Close()

	// Dapatkan MIME type berdasarkan ekstensi
	ext := filepath.Ext(filePath)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// Set header
	w.Header().Set("Content-Type", mimeType)
	w.Header().Set("Content-Disposition", "inline; filename=\""+filename+"\"")

	// Kirim file
	http.ServeFile(w, r, filePath)
}
func factoryFileController(repo fileRepo) fileController {
	if controller == (fileController{}) {
		controller = fileController{
			repo: repo,
		}
	}
	return controller
}
