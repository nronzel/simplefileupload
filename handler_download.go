package main

import (
	"net/http"
	"path"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

func (fs *FileServer) handlerDownload(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "fileName")
	filepath := filepath.Join(fs.UploadPath, path.Base(filename))

	http.ServeFile(w, r, filepath)
}
