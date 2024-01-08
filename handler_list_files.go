package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func (fs *FileServer) handlerListFiles(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(fs.UploadPath)
	if err != nil {
		handleError(w, fmt.Errorf("error reading directory: %w", err))
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	for _, file := range files {
		if _, err := w.Write([]byte(file.Name() + "\n")); err != nil {
			handleError(w, fmt.Errorf("problem writing file name: %w", err))
			return
		}
	}
}

func handleError(w http.ResponseWriter, err error) {
	log.Println(err)
	http.Error(w, "An error occurred", http.StatusInternalServerError)
}
