package ocr

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	vision "cloud.google.com/go/vision/apiv1"
	"cloud.google.com/go/vision/v2/apiv1/visionpb"
	"google.golang.org/api/option"

	"go.uber.org/zap"
)

type OcrService struct {
	log    *zap.Logger
	client *vision.ImageAnnotatorClient
}

func NewOcrService(ctx context.Context, log *zap.Logger, credFile string) (*OcrService, error) {
	c, err := vision.NewImageAnnotatorClient(ctx, option.WithCredentialsFile(credFile))
	if err != nil {
		return nil, err
	}

	return &OcrService{
		log:    log,
		client: c,
	}, nil
}

func (ocr *OcrService) Process(ctx context.Context, fileName string, pageId string) error {
	ocr.log.Info("Processing file", zap.String("fileName", fileName))

	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	img, err := vision.NewImageFromReader(f)
	if err != nil {
		return err
	}

	ta, err := ocr.client.DetectDocumentText(ctx, img, nil)
	if err != nil {
		return err
	}

	if ta.GetPages() == nil || len(ta.GetPages()) == 0 {
		return fmt.Errorf("no pages found")
	} else {
		ocr.log.Info("OCR result", zap.Int("pages", len(ta.GetPages())))
	}

	ocr.log.Debug("OCR result", zap.Any("detectedLanguages", ta.GetPages()[0].Property.DetectedLanguages))
	ocr.log.Debug("OCR result", zap.Any("blocks", len(ta.GetPages()[0].Blocks)))

	// prepare file
	fileNameOnly := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	metadataFile := fileNameOnly + ".json"

	// save hebrew text to file
	ocr.log.Info("Saving text to file", zap.String("fileName", metadataFile))

	// save metadata to file
	md, err := ocr.generateMetadata(fileName, ta.GetPages()[0].Blocks, pageId)
	if err != nil {
		return err
	}
	WriteMetadataFile(metadataFile, md)

	if err != nil {
		return err
	}

	return nil
}

// save hebrew text to file
func (ocr *OcrService) saveToFile(fileName string, data string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}

// filter paragraphs by languge
func (ocr *OcrService) filterParagraphsByLanguage(blocks []*visionpb.Block) string {

	var sb strings.Builder

	for _, block := range blocks {
		if block.BlockType == visionpb.Block_TEXT {
			for _, paragraph := range block.Paragraphs {
				//sb.WriteString("[Paragraph]\n")
				for _, word := range paragraph.Words {
					if ocr.isValidLang(word.Property.GetDetectedLanguages()) {
						for _, symbol := range word.Symbols {
							sb.WriteString(symbol.GetText())
						}
						sb.Write([]byte(" "))
					} else {
						ocr.log.Info("Skipping word", zap.Any("lang", word.Property.GetDetectedLanguages()))
					}
				}
			}
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func (ocr *OcrService) isValidLang(langs []*visionpb.TextAnnotation_DetectedLanguage) bool {
	if langs == nil {
		return false
	}
	//	fmt.Println(langs)
	for _, lang := range langs {
		if lang != nil {
			if lang.GetLanguageCode() == "he" ||
				lang.GetLanguageCode() == "en" ||
				lang.GetLanguageCode() == "yi" {
				return true
			}
		}
	}

	return false

}