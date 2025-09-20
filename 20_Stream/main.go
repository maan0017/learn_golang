package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	fileServer := http.FileServer(http.Dir("./static"))

	http.HandleFunc("/stream-video", streamVideoHandler)
	http.HandleFunc("/stream-audio", streamAudioHandler)
	fmt.Println("server is live at http://localhost:8080/")
	http.Handle("/", fileServer)
	http.ListenAndServe(":8080", nil)
}

func streamVideoHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("./videos/arya4121__-05-07-2025-0001.mp4")
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	// set content type
	w.Header().Set("Content-Type", "video/mp4")
	w.WriteHeader(http.StatusOK)

	// copy chunks in stream
	buff := make([]byte, 1024*32) // 32 KB Chunks

	for {
		n, err := file.Read(buff)
		if err != nil && err != io.EOF {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
		if n == 0 {
			break
		}
		if _, writeErr := w.Write(buff[:n]); writeErr != nil {
			return // client closed connection
		}
		w.(http.Flusher).Flush() // flush buffer to client
	}

}

func streamAudioHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("./audios/_Goëtia. _ Dark Magic Music(MP3_320K).mp3")
	if err != nil {
		http.Error(w, "something went wrong", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// set content type
	w.Header().Set("Content-Type", "audio/mp3")
	w.WriteHeader(http.StatusOK)

	buff := make([]byte, 1024*32) // 32 KB Chunks

	for {

		n, err := file.Read(buff)
		if err != nil && err != io.EOF {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
		if n == 0 {
			break
		}
		if _, err := w.Write(buff[:n]); err != nil {
			return
		}

		w.(http.Flusher).Flush()
	}

}
