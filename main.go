package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
)

const (
	Port       = "8888"
	MaxUpload  = 10 << 20 // 10mb
	UploadPath = "./uploads"
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
		Port:          Port,
		MaxUploadSize: MaxUpload,
	}

	fs.Router.Post("/upload", fs.handlerUpload)
	fs.Router.Get("/files/{fileName}", fs.handlerDownload)
	fs.Router.Get("/files", fs.handlerListFiles)

	return fs
}

func ensureDir(dirName string) {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err := os.Mkdir(dirName, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	ensureDir(UploadPath)

	fileServer := newFileServer(UploadPath)
	log.Printf("Server started on localhost:%s\n", fileServer.Port)

	srv := &http.Server{
		Addr:    ":" + fileServer.Port,
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

	log.Println("Server gracefully shutdown")
}
