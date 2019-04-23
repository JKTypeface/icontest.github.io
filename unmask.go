package main

import (
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"log"
	"math"
	"net/url"
	"path"
)

// UnmaskCommand holds the option for removing the watermark off of an icon.
type UnmaskCommand struct {
	Args struct {
		Image string `positional-arg-name:"<url/path>"`
	} `positional-args:"1" required:"1"`

	Output string `short:"o" description:"Output file to save the looted image to."`
}

// Execute runs the unmask command which removes the watermark off of an icon.
func (cmd *UnmaskCommand) Execute(args []string) error {

	if cmd.Output == "" {
		fileName := path.Base(cmd.Args.Image)
		cmd.Output = fileName
	}

	var img image.Image
	var err error

	if isURL(cmd.Args.Image) {

		img, err = getImage(cmd.Args.Image)

	} else {

		img, err = loadImage(cmd.Args.Image)

	}

	if err != nil {
		return err
	}

	img = Unmask(img)

	err = storeImage(img, cmd.Output)

	return err
}

// Unmask removes the watermark from an IconFinder.com premium image.
func Unmask(img image.Image) image.Image {
	size := img.Bounds().Size()

	cimg, _ := img.(draw.Image)

	nodes := getNodesForSize(size.Y)

	for i, node := range nodes {

		if node == 0 {
			continue
		}

		x := i % size.X
		y := i / size.Y

		pixel, conv := splitConvert(img.At(x, y))

		var newPixel color.Color

		if transparent(conv) {

			newPixel = color.RGBA{0, 0, 0, 0}

		} else {

			if (x > 1 && nodes[i-1] == 0) ||
				(y < size.Y-1 && nodes[i+1] == 0) {

				newPixel = rec(pixel, 0.015)

			} else if (x > 2 && nodes[i-2] == 0) ||
				(y < size.Y-2 && nodes[i+2] == 0) {

				newPixel = rec(pixel, 0.065)

			} else {

				newPixel = rec(pixel, 0.073)

			}

		}

		cimg.Set(x, y, newPixel)
	}

	return cimg.(image.Image)
}

func transparent(pixel color.NRGBA) bool {
	r := pixel.R
	g := pixel.G
	b := pixel.B
	a := pixel.A

	return a != 255 &&
		66 < r && r < 70 &&
		66 < g && g < 70 &&
		62 < b && b < 66
}

func rec(pixel color.RGBA, mult float64) color.Color {
	r := (float64(pixel.R) - 68*mult) / (1 - mult)
	g := (float64(pixel.G) - 68*mult) / (1 - mult)
	b := (float64(pixel.B) - 64*mult) / (1 - mult)
	a := pixel.A

	ret := color.NRGBA{
		R: uint8(math.Min(math.Max(r, 0), 255)),
		G: uint8(math.Min(math.Max(g, 0), 255)),
		B: uint8(math.Min(math.Max(b, 0), 255)),
		A: uint8(a),
	}

	return ret
}

func getNodesForSize(size int) []int {
	if size == 48 {
		return getNodes(48, 48, 5, 5, 2)
	}
	if size == 128 {
		return getNodes(128, 128, 13, 7, 3)
	}
	if size == 256 {
		return getNodes(256, 256, 13, 7, 3)
	}
	if size == 512 {
		return getNodes(512, 512, 13, 7, 3)
	}
	log.Fatalf("Unsupported size: %v", size)
	return nil
}

func getNodes(height, width, interval, datalen, start int) []int {
	var nodes []int
	rawLine := repeat(interval, 0)
	rawLine = append(rawLine, repeat(datalen, 1)...)
	packLen := interval + datalen
	raw := repeat(int(width/packLen)+2, rawLine...)
	start = packLen - start
	for i := 0; i < height; i++ {
		if start > packLen {
			start = start % packLen
		}

		nodes = append(nodes, raw[start:start+width]...)

		start++
	}
	return nodes
}

func repeat(times int, values ...int) []int {
	var arr []int
	for i := 0; i < times; i++ {
		arr = append(arr, values...)
	}
	return arr
}

func isURL(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	return err == nil
}

func splitConvert(pixel color.Color) (rgba color.RGBA, nrgba color.NRGBA) {

	switch pixel.(type) {

	case color.NRGBA:
		nrgba = pixel.(color.NRGBA)
		rgba = color.RGBAModel.Convert(pixel).(color.RGBA)

	case color.RGBA:
		rgba = pixel.(color.RGBA)
		nrgba = color.NRGBAModel.Convert(pixel).(color.NRGBA)
	}

	return
}
