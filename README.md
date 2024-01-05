# simplefileupload

## Overview

A very simple file upload server written in Go.

## Endpoints

| Method | Endpoint         | Description                       |
| ------ | ---------------- | --------------------------------- |
| POST   | /upload          | Upload a file to the server       |
| GET    | /files           | Gets a list of the uploaded files |
| GET    | /file/{fileName} | Download a specific file          |

Set up with a 10mb limit per file.

## Installation

### Prerequisites

- Go 1.20 or higher
- Chi router

### Setup

Clone the repo:

```bash
git clone https://github.com/nronzel/simplefileupload.git
```

Navigate to the clone directory:

```bash
cd simplefileupload
```

Install dependencies:

```bash
go mod tidy
```

## Usage

### Starting the Server

Run the following command in the root of the project:

```bash
go run main.go
```

> The server will start and listen on port 8888

### Uploading a File

Use a tool like `curl`, `Postman`, `Insomnia`, etc. to upload a file:

```bash
curl -X POST -F "file=@/path/to/file.txt" http://localhost:8888/upload
```

### Download a File

```bash
curl http://localhost:8888:files/file.txt -o file.txt
```

### List All Files

```bash
curl http://localhost:8888/files
```

## Testing

To run the included test suite:

```bash
go test -v
```

## Contributions

This is meant to serve as a starting point for a fileserver. Feel free to fork
this repo and add any functionality you may need or want.
