package main

import "net/http"

type FileServer struct {
	RootDir string
}

func NewFileServer(rootDir string) *FileServer {
	return &FileServer{
		RootDir: rootDir,
	}
}

func (fs *FileServer) StartFileServer() http.Handler {
	return http.FileServer(http.Dir(fs.RootDir))
}
