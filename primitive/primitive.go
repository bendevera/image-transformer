package primitive 

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"io"
	"io/ioutil"
	"errors"
)

type Mode int

const (
	ModeCombo Mode = iota
	ModeTriangle
	ModeRect
	ModeEllipse
	ModeCircle
	ModeRotatedRect
	ModeBeziers
	ModeRotatedEllipse
	ModePolygon
)

// option for the transform function
func WithMode(mode Mode) func() []string {
	return func() []string {
		return []string{"-m", fmt.Sprintf("%d", mode)}
	}
}

// will take the provided image and apply primitive 
// tranformation and return a reader to the resulting image
func Transform(image io.Reader, ext string, numShapes int, opts ...func() []string) (io.Reader, error) {
	var args []string
	for _, opt := range opts {
		args = append(args, opt()...)
	}
	in, err := tempfile("in_", ext)
	if err != nil {
		return nil, errors.New("primitive: failed to create temporary input file")
	}
	defer os.Remove(in.Name())
	out, err := tempfile("out_", ext)
	if err != nil {
		return nil, errors.New("primitive: failed to create temporary ouput file")
	}
	defer os.Remove(out.Name())

	// read image into in file
	_, err = io.Copy(in, image)
	if err != nil {
		return nil, errors.New("primitive: failed to copy input to temp input file")
	}
	// run primitive w/ -i in.Name() -o out.Name()
	stdCombo, err := primitive(in.Name(), out.Name(), numShapes, args...)
	if err != nil {
		return nil, errors.New("primitive: failed to run primitive on input")
	}
	fmt.Println(stdCombo)
	// read out into a reader, return reader, delete out 
	b := bytes.NewBuffer(nil)
	_, err = io.Copy(b, out)
	if err != nil {
		return nil, errors.New("primitive: failed to read primitive output into reader")
	}

	return b, nil
}

func primitive(inputFile, outputFile string, numShapes int, args ...string) (string, error) {
	argStr := fmt.Sprintf("-i %s -o %s -n %d", inputFile, outputFile, numShapes)
	args = append(strings.Fields(argStr), args...)
	cmd := exec.Command("primitive", args...)
	b, err := cmd.CombinedOutput()
	return string(b), err
}

func tempfile(prefix, ext string) (*os.File, error) {
	in, err := ioutil.TempFile("", prefix)
	if err != nil {
		return nil, errors.New("primitive: failed to create temporary input file")
	}
	defer os.Remove(in.Name())
	return os.Create(fmt.Sprintf("%s.%s", in.Name(), ext))
}