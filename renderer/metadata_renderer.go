package renderer

import (
	"encoding/json"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/gipsh/yizk/model"
	"go.uber.org/zap"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// MetadataRenderer renders the metadata of a page into a new translated image
// It needs the original image and the metadata json file as input

type MetadataRenderer struct {
	log *zap.Logger
}

func NewMetadataRenderer(log *zap.Logger) *MetadataRenderer {
	return &MetadataRenderer{
		log: log,
	}
}

func (mr *MetadataRenderer) getFontFace() font.Face {
	return basicfont.Face7x13
}

func (mr *MetadataRenderer) fontHeight() float64 {
	return 13
}

func (mr *MetadataRenderer) MeasureString(s string) (w, h float64) {
	d := &font.Drawer{
		Face: mr.getFontFace(),
	}
	a := d.MeasureString(s)
	return float64(a >> 6), mr.fontHeight()
}

func (mr *MetadataRenderer) RenderPage(filename string) error {

	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var page model.YizkPage
	err = json.Unmarshal([]byte(file), &page)
	if err != nil {
		return err
	}

	// open original image file
	img, err := ReadFile(page.Filename)
	if err != nil {
		return err
	}

	// copy content of original image into new image
	newImage := image.NewRGBA((*img).Bounds())
	draw.Draw(newImage, (*img).Bounds(), *img, image.Point{X: 0, Y: 0}, draw.Src)

	for _, block := range page.Blocks {
		rect := image.Rectangle{}

		//rect.Min = image.Point{X: block.Points[1].X, Y: block.Points[1].Y}
		//rect.Max = image.Point{X: block.Points[0].X, Y: block.Points[0].Y}

		rect.Min = image.Point{X: block.UpperLeftPoint.X, Y: block.UpperLeftPoint.Y}
		rect.Max = image.Point{X: block.BottomRightPoint.X, Y: block.BottomRightPoint.Y}

		draw.Draw(newImage, rect, *img, image.Point{X: 0, Y: 0}, draw.Src)

		//mr.addLabelWrapped(newImage, block.Points[1].X, block.Points[1].Y, block.Points[0].X, block.TranslatedText)
		mr.addLabelWrapped(newImage, block.UpperLeftPoint.X, block.UpperLeftPoint.Y, block.BottomRightPoint.X, block.TranslatedText)

	}

	// save to file
	fileNameOnly := strings.TrimSuffix(page.Filename, filepath.Ext(page.Filename))

	newFilename := fileNameOnly + ".render" + ".png"

	return WritePNGFile(newFilename, newImage)

}

func addLabel(img *image.RGBA, x, y int, label string) {
	color := color.RGBA{200, 100, 0, 255}
	point := fixed.Point26_6{fixed.I(x), fixed.I(y)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

func (mr *MetadataRenderer) addLabelWrapped(img *image.RGBA, x, y, limitX int, label string) {

	lines := wordWrap(mr, label, float64(limitX))
	for i, line := range lines {
		//mr.log.Debug("addLabelWrapped", zap.String("line", line))
		addLabel(img, x, y+i*int(mr.fontHeight()), line)
	}
}

// Write RGBA struct as png file
func WriteFile(filename string, im *image.Image) error {
	fd, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fd.Close()
	return jpeg.Encode(fd, *im, nil)
}

func WritePNGFile(filename string, im *image.RGBA) error {
	fd, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fd.Close()
	return png.Encode(fd, im)
}

func ReadFile(filename string) (im *image.Image, err error) {
	fd, err := os.Open(filename)
	if err != nil {
		return
	}
	img, err := jpeg.Decode(fd)
	if err != nil {
		return
	}
	return &img, nil
}
