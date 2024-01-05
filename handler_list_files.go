package main

import (
	"log"
	"net/http"
	"os"
)

func (fs *FileServer) listFilesHandler(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(fs.UploadPath)
	if err != nil {
		log.Printf("Error reading directory: %v", err)
		http.Error(w, "Error reading directory", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	for _, file := range files {
		if _, err := w.Write([]byte(file.Name() + "\n")); err != nil {
			log.Printf("Problem writing file name: %v", err)
			http.Error(w, "Problem writing file", http.StatusInternalServerError)
			return
		}
	}
}
