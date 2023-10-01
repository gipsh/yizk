package translator

import (
	"context"
	"encoding/json"
	"os"

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
func (mdt *MetadataTranslator) TranslatePage(ctx context.Context, filename string, ts *TranslatorService) error {

	file, _ := os.ReadFile(filename)

	var page model.YizkPage
	err := json.Unmarshal([]byte(file), &page)
	if err != nil {
		return err
	}

	for i, block := range page.Blocks {
		translatedText, err := ts.Translate(ctx, block.Text)
		if err != nil {
			return err
		}
		page.Blocks[i].TranslatedText = translatedText
	}

	// save to file
	return ocr.WriteMetadataFile(filename, &page)

}
