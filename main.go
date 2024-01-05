package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
)

type FileServer struct {
	UploadPath    string
	Router        *chi.Mux
	Port          string
	MaxUploadSize int64
}

func newFileServer(uploadPath string) *FileServer {
	fs := &FileServer{
		UploadPath:    uploadPath,
		Router:        chi.NewRouter(),
		Port:          ":8888",
		MaxUploadSize: 10 << 20, // 10MB size limit
	}

	fs.Router.Post("/upload", fs.uploadFileHandler)
	fs.Router.Get("/files/{fileName}", fs.downloadFileHandler)
	fs.Router.Get("/files", fs.listFilesHandler)

	return fs
}

func (fs *FileServer) uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(fs.MaxUploadSize)
	if err != nil {
		log.Printf("Error parsing multipart form: %v", err)
		http.Error(w, "File too large", http.StatusBadRequest)
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
		return
	}

	_, err = w.Write([]byte("File uploaded successfully: " + handler.Filename))
	if err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
	}
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

	srv := &http.Server{
		Addr:    fileServer.Port,
		Handler: fileServer.Router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server gracefully stopped")
}

func ensureDir(dirName string) {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err := os.Mkdir(dirName, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
}
