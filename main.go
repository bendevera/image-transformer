package main

import (
	"io"
	"net/http"
	"fmt"
	"path/filepath"
	"github.com/imagetransformer/primitive"
	"log"
)


func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `<html><body>
				<form action="/upload" method="post" enctype="multipart/form-data">
					<input type="file" name="image" />
					<button type="submit">Upload image</button>
				</body></html>`
		fmt.Fprint(w, html)
	})
	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		file, header, err := r.FormFile("image")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return 
		}
		defer file.Close()
		// todo: actually use this
		ext := filepath.Ext(header.Filename)[1:]
		switch ext {
		case "jpg":
			fallthrough
		case "jpeg":
			w.Header().Set("Content-Type", "image/jpeg")
		case "png":
			w.Header().Set("Content-Type", "image/png")
		default:
			http.Error(w, "Invalid image type", http.StatusBadRequest)
			return
		}
		out, err := primitive.Transform(file, ext, 30)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return 
		}
		io.Copy(w, out)
	})
	log.Fatal(http.ListenAndServe(":3000", mux))
}
