package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

type FileServer struct {
	UploadPath string
	Router     *chi.Mux
	Port       string
}

func newFileServer(uploadPath string) *FileServer {
	fs := &FileServer{
		UploadPath: uploadPath,
		Router:     chi.NewRouter(),
		Port:       ":8888",
	}

	fs.Router.Post("/upload", fs.uploadFileHandler)
	fs.Router.Get("/files/{fileName}", fs.downloadFileHandler)
	fs.Router.Get("/files", fs.listFilesHandler)

	return fs
}

func (fs *FileServer) uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10MB size limit
	if err != nil {
		http.Error(w, "Problem parsing file", http.StatusInternalServerError)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	filepath := filepath.Join(fs.UploadPath, handler.Filename)
	dst, err := os.Create(filepath)
	if err != nil {
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error writing the file", http.StatusInternalServerError)
	}

	w.Write([]byte("File uploaded successfully: " + handler.Filename))
}

func (fs *FileServer) downloadFileHandler(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "fileName")
	filepath := filepath.Join(fs.UploadPath, filename)

	http.ServeFile(w, r, filepath)
}

func (fs *FileServer) listFilesHandler(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(fs.UploadPath)
	if err != nil {
		http.Error(w, "Error reading directory", http.StatusInternalServerError)
		return
	}

	for _, file := range files {
		_, err := w.Write([]byte(file.Name() + "\n"))
		if err != nil {
			http.Error(w, "Problem writing file", http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	uploadPath := "./uploads"
	ensureDir(uploadPath)

	fileServer := newFileServer(uploadPath)
	log.Printf("Server started on localhost%s\n", fileServer.Port)
	log.Fatal(http.ListenAndServe(fileServer.Port, fileServer.Router))
}

func ensureDir(dirName string) {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err := os.Mkdir(dirName, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
}
