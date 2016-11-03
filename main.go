package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/fatih/color"
	"image"
	"image/draw"
	"image/png"
	"os"
	"text/template"
)

func init() {
	log.SetLevel(log.WarnLevel)
}

const UsageTmpl = `{{ . }} - convert text to qrcode
	{{ . }} text
`

func main() {

	if len(os.Args) != 2 {
		usage := template.Must(template.New("usage").Parse(UsageTmpl))
		usage.Execute(os.Stderr, os.Args[0])
		return
	}
	content := os.Args[1]

	black := color.New(color.BgBlack).SprintFunc()
	white := color.New(color.BgWhite).SprintFunc()
	pixel := "  "
	width := 300
	height := width
	filename := "qrcode.png"
	log.WithFields(
		log.Fields{
			"content": content,
		},
	).Info("Create QRcode image(.png)")

	qrcode, err := qr.Encode(content, qr.H, qr.Auto)
	if err != nil {
		log.WithError(err).Warn("encode qrcode failure")
		return
	}

	log.WithFields(
		log.Fields{
			"width":  width,
			"height": height,
		},
	).Info("scale qrcode")

	// output to console
	rect := qrcode.Bounds()
	min := rect.Min
	max := rect.Max

	qrimage := image.NewRGBA(image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{max.X + 2, max.Y + 2},
	})

	draw.Draw(qrimage, qrimage.Bounds(), qrcode, image.Point{-1, -1}, draw.Src)

	rect = qrimage.Bounds()
	min = rect.Min
	max = rect.Max

	for y := min.Y; y < max.Y; y++ {
		for x := min.X; x < max.X; x++ {
			r, g, b, a := qrimage.At(x, y).RGBA()
			if r == 0 && g == 0 && b == 0 && a == 0xffff {
				fmt.Print(black(pixel))
			} else {
				fmt.Print(white(pixel))
			}
		}
		fmt.Println()
	}

	qrcode, err = barcode.Scale(qrcode, width, height)
	if err != nil {
		log.WithError(err).Warn("scale qrcode failure")
	}

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, os.FileMode(0660))
	if err != nil {
		log.WithError(err).WithFields(
			log.Fields{
				"filename": filename,
			},
		).Warn("open file failure")
		return
	}
	defer file.Close()

	err = png.Encode(file, qrcode)
	if err != nil {
		log.WithError(err).Warn("write qrcode to file(.png)")
		return
	}

}
