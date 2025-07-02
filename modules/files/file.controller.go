package files

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type fileController struct {
}

var controller fileController

func (f fileController) UploadFileCtrl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println("Error Retrieving the File")
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
		return
	}
	defer file.Close()
	log.Printf("Uploaded File: %+v\n", handler.Filename)
	log.Printf("File Size: %+v\n", handler.Size)
	log.Printf("MIME Header: %+v\n", handler.Header)
	tempFile, err := os.CreateTemp("uploads", "upload-*.png")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
		return
	}
	defer tempFile.Close()
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
		return
	}
	tempFile.Write(fileBytes)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"File uploaded successfully"}`)
}

func factoryFileController() fileController {
	if controller == (fileController{}) {
		controller = fileController{}
	}
	return controller
}
