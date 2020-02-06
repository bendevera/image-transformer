package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"fmt"
	"path/filepath"
	"github.com/imagetransformer/primitive"
	"log"
	"errors"
	"os"
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

		out, err := primitive.Transform(file, ext, 30, primitive.WithMode(primitive.ModeRect))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return 
		}
		outFile, err := tempfile("out_", ext)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer outFile.Close()
		io.Copy(outFile, out)
		redirURL := fmt.Sprintf("/%s", outFile.Name())
		http.Redirect(w, r, redirURL, http.StatusFound)
	})
	fs := http.FileServer(http.Dir("./img/"))
	mux.Handle("/img/", http.StripPrefix("/img/", fs))
	log.Fatal(http.ListenAndServe(":3000", mux))
}

func tempfile(prefix, ext string) (*os.File, error) {
	in, err := ioutil.TempFile("./img/", prefix)
	if err != nil {
		return nil, errors.New("main: failed to create temporary input file")
	}
	defer os.Remove(in.Name())
	return os.Create(fmt.Sprintf("%s.%s", in.Name(), ext))
}