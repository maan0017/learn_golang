package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type UploadResult struct {
	PublicID     string `json:"public_id"`
	Version      int64  `json:"version"`
	Signature    string `json:"signature"`
	Width        int64  `json:"width"`
	Height       int64  `json:"height"`
	Format       string `json:"format"`
	ResourceType string `json:"resource_type"`
	CreatedAt    string `json:"created_at"`
	URL          string `json:"url"`        // HTTP URL
	SecureURL    string `json:"secure_url"` // HTTPS URL
	Bytes        int64  `json:"bytes"`
	Type         string `json:"type"`
}

func handler(w http.ResponseWriter, _ *http.Request) {
	// Allow CORS by setting the necessary headers
	w.Header().Set("Access-Control-Allow-Origin", "*")                       // Allows all origins
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE") // Allowed methods
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")           // Allowed headers

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("CORS is allowed."))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Allow CORS by setting the necessary headers
	w.Header().Set("Access-Control-Allow-Origin", "*")                       // Allows all origins
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE") // Allowed methods
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")           // Allowed headers

	if r.Method != "POST" {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse uploaded file
	// Parse the multipart form (maxMemory in bytes)
	file, header, err := r.FormFile("file") // "file" is the form field name
	if err != nil {
		http.Error(w, "Failed to read file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Optional: Print file info
	fmt.Printf("Received file: %s (%d bytes)\n", header.Filename, header.Size)

	// Save file locally (for demo purposes)
	// dst, err := os.Create("./uploads/" + header.Filename)
	// if err != nil {
	// 	http.Error(w, "Failed to save file: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// defer dst.Close()
	// _, err = io.Copy(dst, file)
	// if err != nil {
	// 	http.Error(w, "Failed to write file: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// Init Cloudinary
	// cld, _ := cloudinary.New()
	cld, _ := cloudinary.NewFromParams("dlnbwxce8", "498931364168295", "WfTWdMqm8Z42xImdFg_D-VvZmVU")

	// Upload
	fmt.Println("Uploading file to Cloudinary...")
	// uploadResp, err := cld.Upload.Upload(context.Background(), file, uploader.UploadParams{
	// 	PublicID: header.Filename, // optional: can generate unique ID
	// 	Folder:   "image-uploads", // optional: folder in Cloudinary
	// })
	filePath := fmt.Sprintf("uploads/%s", header.Filename)
	fmt.Println("file path: ", filePath)
	uploadResp, err := cld.Upload.Upload(context.Background(), filePath, uploader.UploadParams{
		PublicID: header.Filename, // optional: can generate unique ID
		Folder:   "image-uploads", // optional: folder in Cloudinary
	})
	if err != nil {
		http.Error(w, "Upload failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Cloudinary upload response: %+v\n", uploadResp)

	w.Header().Set("Content-Type", "application/json")
	fmt.Println("file url: ", uploadResp.SecureURL)
	// Send back uploaded file URL
	fmt.Fprintf(w, `{"url": "%s"}`, uploadResp.SecureURL)
}

func main() {
	// Ensure uploads folder exists
	os.MkdirAll("./uploads", os.ModePerm)

	http.HandleFunc("/", handler)
	http.HandleFunc("/upload", uploadHandler)
	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
