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
	"html/template"
)

type Result struct {
	Path 		string
	NumShapes	int
	Mode 		string
}


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

		var ResultA Result
		a, err := genImage(file, ext, 33, primitive.ModeCombo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ResultA = setResult(ResultA, a, 33, "Combo")
		file.Seek(0, 0)
		var ResultB Result
		b, err := genImage(file, ext, 33, primitive.ModeTriangle)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return 
		}
		ResultB = setResult(ResultB, b, 33, "Triangle")
		file.Seek(0, 0)
		var ResultC Result
		c, err := genImage(file, ext, 33, primitive.ModeRect)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ResultC = setResult(ResultC, c, 33, "Rectangle")
		file.Seek(0, 0)
		var ResultD Result
		d, err := genImage(file, ext, 33, primitive.ModeEllipse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ResultD = setResult(ResultD, d, 33, "Ellipse")
		tpl := template.Must(template.ParseFiles("result.html"))
		results := []Result{ResultA, ResultB, ResultC, ResultD}
		tpl.Execute(w, results)
	})
	fs := http.FileServer(http.Dir("./img/"))
	mux.Handle("/img/", http.StripPrefix("/img/", fs))
	log.Fatal(http.ListenAndServe(":3000", mux))
}

func setResult(r Result, path string, numShapes int, mode string) (Result) {
	r.Path = path
	r.NumShapes = numShapes
	r.Mode = mode 
	return r
}

func genImage(r io.Reader, ext string, numShapes int, mode primitive.Mode) (string, error) {
	out, err := primitive.Transform(r, ext, numShapes, primitive.WithMode(mode))
	if err != nil {
		return "", err
	}
	outFile, err := tempfile("out_", ext)
	if err != nil {
		return "", err
	}
	defer outFile.Close()
	io.Copy(outFile, out)
	return outFile.Name(), nil
}

func tempfile(prefix, ext string) (*os.File, error) {
	in, err := ioutil.TempFile("./img/", prefix)
	if err != nil {
		return nil, errors.New("main: failed to create temporary input file")
	}
	defer os.Remove(in.Name())
	return os.Create(fmt.Sprintf("%s.%s", in.Name(), ext))
}