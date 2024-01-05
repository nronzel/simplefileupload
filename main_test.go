package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupTestServer() (*FileServer, string) {
	tempDir, _ := os.MkdirTemp("", "testuploads")
	fileServer := newFileServer(tempDir)
	return fileServer, tempDir
}

func TestUploadFileHandler(t *testing.T) {
	fileServer, tempDir := setupTestServer()
	defer os.RemoveAll(tempDir)

	// Create a test file
	testFileContent := "test content"
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	part, _ := writer.CreateFormFile("file", "test.txt")
	_, err := io.Copy(part, strings.NewReader(testFileContent))
	if err != nil {
		t.Errorf("problem copying file: %v", err)
	}
	writer.Close()

	// Create a request with the test file
	req, _ := http.NewRequest("POST", "/upload", buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()

	// Call the handler
	fileServer.Router.ServeHTTP(rr, req)

	// Check the response and file creation
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check if file was saved correctly
	savedFile, _ := os.ReadFile(filepath.Join(fileServer.UploadPath, "test.txt"))
	if string(savedFile) != testFileContent {
		t.Errorf("file content mismatch: got %v want %v", string(savedFile), testFileContent)
	}

	// Clean up
	os.Remove(filepath.Join(fileServer.UploadPath, "test.txt"))
}

func TestDownloadFileHandler(t *testing.T) {
	fileServer, tempDir := setupTestServer()
	defer os.RemoveAll(tempDir)

	// Place a test file in the uploads directory
	testFileName := "downloadTest.txt"
	testFileContent := "download test content"
	err := os.WriteFile(filepath.Join(fileServer.UploadPath, testFileName), []byte(testFileContent), 0644)
	if err != nil {
		t.Errorf("problem testing download: %v", err)
	}

	// Create a request to download the file
	req, _ := http.NewRequest("GET", "/files/"+testFileName, nil)
	rr := httptest.NewRecorder()

	// Call the handler
	fileServer.Router.ServeHTTP(rr, req)
	// Check the response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	if rr.Body.String() != testFileContent {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), testFileContent)
	}

	// Clean up
	os.Remove(filepath.Join(fileServer.UploadPath, testFileName))
}

func TestListFilesHandler(t *testing.T) {
	fileServer, tempDir := setupTestServer()
	defer os.RemoveAll(tempDir)

	// Create a test file
	testFileName := "listTest.txt"
	err := os.WriteFile(filepath.Join(fileServer.UploadPath, testFileName), []byte("content"), 0644)
	if err != nil {
		t.Errorf("problem writing listTest.txt: %v", err)
	}

	// Create a request to list files
	req, _ := http.NewRequest("GET", "/files", nil)
	rr := httptest.NewRecorder()

	// Call the handler
	fileServer.Router.ServeHTTP(rr, req)

	// Check the response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	if !strings.Contains(rr.Body.String(), testFileName) {
		t.Errorf("handler did not return expected file: got %v want to contain %v", rr.Body.String(), testFileName)
	}

	// Clean up
	os.Remove(filepath.Join(fileServer.UploadPath, testFileName))
}
