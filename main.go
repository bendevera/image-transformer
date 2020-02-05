package main

import (
	"os"
	"io"
	"github.com/imagetransformer/primitive"
)

func main() {
	inFile, err := os.Open("/Users/bendevera/Desktop/unit2-demo/plane2.jpeg")
	if err != nil {
		panic(err)
	}
	defer inFile.Close()
	out, err := primitive.Transform(inFile, 30)
	if err != nil {
		panic(err)
	}
	os.Remove("out.png")
	outFile, err := os.Create("out.png")
	if err != nil {
		panic(err)
	}
	io.Copy(outFile, out)
	// out, err := primitive("/Users/bendevera/Desktop/unit2-demo/plane2.jpeg", "output.png", 50, triangle)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf(out)
}

// type PrimitiveMode int

// const (
// 	combo PrimitiveMode = iota
// 	triangle
// 	rect
// 	ellipse
// 	circle
// 	rotatedrect
// 	beziers
// 	rotatedellipse
// 	polygon
// )

// func primitive(inputFile, outputFile string, numShapes int, mode PrimitiveMode) (string, error) {
// 	argStr := fmt.Sprintf("-i %s -o %s -n %d -m %d", inputFile, outputFile, numShapes, mode)
// 	cmd := exec.Command("primitive", strings.Fields(argStr)...)
// 	b, err := cmd.CombinedOutput()
// 	return string(b), err
// }