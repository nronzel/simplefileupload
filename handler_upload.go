package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func (fs *FileServer) handlerUpload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(fs.MaxUploadSize)
	if err != nil {
		log.Printf("Error parsing multipart form: %v", err)
		http.Error(w, "File too large or incorrect format", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error retrieving the file: %v", err)
		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	safeFileName := filepath.Base(handler.Filename)
	filepath := filepath.Join(fs.UploadPath, safeFileName)
	dst, err := os.Create(filepath)
	if err != nil {
		log.Printf("Error creating the file: %v", err)
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		log.Printf("Error writing the file: %v", err)
		http.Error(w, "Error writing the file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte("File uploaded successfully: " + safeFileName))
	if err != nil {
		http.Error(w, "Problem writing response body", http.StatusInternalServerError)
	}
}
