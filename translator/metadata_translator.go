package translator

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/gipsh/yizk/model"
	"github.com/gipsh/yizk/ocr"
	"go.uber.org/zap"
)

type MetadataTranslator struct {
	log *zap.Logger
	ts  *TranslatorService
}

func NewMetadataTranslator(log *zap.Logger, ts *TranslatorService) *MetadataTranslator {
	return &MetadataTranslator{
		log: log,
		ts:  ts,
	}
}

// Read a metadata OcrPage json file and translate the text
func (mdt *MetadataTranslator) TranslatePage(ctx context.Context, filename string) error {

	file, _ := os.ReadFile(filename)

	var page model.YizkPage
	err := json.Unmarshal([]byte(file), &page)
	if err != nil {
		return err
	}

	for i, block := range page.Blocks {
		translatedText, err := mdt.ts.Translate(ctx, block.Text)
		if err != nil {
			return err
		}
		page.Blocks[i].TranslatedText = translatedText
	}

	// save to file
	return ocr.WriteMetadataFile(filename, &page)

}

func (mdt *MetadataTranslator) TranslateFolder(ctx context.Context, folder string) error {

	files, err := os.ReadDir(folder)
	if err != nil {
		mdt.log.Error("Error reading folder", zap.String("folder", folder), zap.Error(err))
		return err
	}

	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".json") {
			mdt.log.Info("Processing file", zap.String("fileName", filepath.Join(folder, f.Name())))
			err = mdt.TranslatePage(ctx, filepath.Join(folder, f.Name()))
			if err != nil {
				mdt.log.Error("Error translating file", zap.String("fileName", f.Name()), zap.Error(err))
			}
		}
	}

	return nil

}
