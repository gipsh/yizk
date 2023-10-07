package renderer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gipsh/yizk/model"
	"go.uber.org/zap"
)

// TextRenderer renders the metadata of a page into a new translated image
// It needs the original image and the metadata json file as input

type TextRenderer struct {
	log *zap.Logger
}

func NewTextRenderer(log *zap.Logger) Renderer {
	return &TextRenderer{
		log: log,
	}
}

func (mr *TextRenderer) RenderFolder(folder string, outputFilename *string) error {

	var textFile string
	if outputFilename != nil {
		textFile = *outputFilename
	} else {
		x := filepath.Base(folder)
		textFile = fmt.Sprintf("book-%s.txt", x)
	}

	mr.log.Info("Saving file", zap.String("fileName", textFile))

	outputFile, err := os.Create(textFile)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	files, err := os.ReadDir(folder)
	if err != nil {
		mr.log.Error("Error reading folder", zap.String("folder", folder), zap.Error(err))
		return err
	}

	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".json") {
			mr.log.Info("Processing file", zap.String("fileName", f.Name()))
			pageText, err := mr.RenderPage(filepath.Join(folder, f.Name()))
			if err != nil {
				mr.log.Error("Error rendering file", zap.String("fileName", f.Name()), zap.Error(err))
			}

			_, err = outputFile.WriteString(pageText)
			if err != nil {
				mr.log.Error("Error writing to file", zap.String("fileName", textFile), zap.Error(err))
			}

		}
	}

	mr.log.Info("Saving file", zap.String("fileName", textFile))
	mr.log.Info("Done")
	return nil
}

func (mr *TextRenderer) RenderPage(filename string) (string, error) {

	file, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	var page model.YizkPage
	err = json.Unmarshal([]byte(file), &page)
	if err != nil {
		return "", err
	}

	var sb strings.Builder

	sb.Write([]byte("--\n"))
	for _, block := range page.Blocks {

		sb.Write([]byte(block.TranslatedText))
		// add a double new line after each block
		sb.Write([]byte("\n\n"))

	}

	// add page number
	sb.Write([]byte(fmt.Sprintf("Page %s\n--\n", page.PageNumber)))
	sb.Write([]byte("\n\n"))

	// save to file
	fileNameOnly := strings.TrimSuffix(page.Filename, filepath.Ext(page.Filename))

	newFilename := fileNameOnly + ".txt"

	mr.log.Info("Saving file", zap.String("fileName", newFilename))
	err = WriteTextFile(newFilename, sb.String())
	if err != nil {
		return "", err
	}

	return sb.String(), nil

}

func WriteTextFile(filename string, text string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(text)
	if err != nil {
		return err
	}
	return nil
}
