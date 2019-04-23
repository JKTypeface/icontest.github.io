package main

import (
	"bufio"
	"bytes"
	"image"
	"image/png"
	_ "image/png"
	"log"
	"net/http"
	"os"

	flags "github.com/jessevdk/go-flags"
)

// Options holds the command line options for the cli program.
var Options struct {
	Unmask UnmaskCommand `command:"unmask" alias:"u" description:"Remove the watermark off of an IconFinder image or url."`
	Server ServerCommand `command:"server" alias:"s" description:"Run an HTTP server that provides the services of this tool."`
}

func main() {

	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	parser := flags.NewParser(&Options, flags.HelpFlag|flags.PassDoubleDash|flags.PrintErrors)

	parser.Parse()

}

func loadImage(fileName string) (image.Image, error) {
	imgFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()
	return png.Decode(imgFile)
}

func storeImage(img image.Image, fileName string) error {
	imgFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer imgFile.Close()
	return png.Encode(imgFile, img)
}

func getImage(url string) (image.Image, error) {

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)

	return png.Decode(reader)
}

func imageBytes(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := png.Encode(buf, img)
	return buf.Bytes(), err
}
