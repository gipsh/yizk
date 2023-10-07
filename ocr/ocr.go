package ocr

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	vision "cloud.google.com/go/vision/apiv1"
	"cloud.google.com/go/vision/v2/apiv1/visionpb"
	"github.com/gipsh/yizk/model"
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

func (ocr *OcrService) ProcessByFolder(ctx context.Context, folder string) error {
	ocr.log.Info("Processing folder", zap.String("folder", folder))

	files, err := os.ReadDir(folder)
	if err != nil {
		return err
	}

	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".json") {
			ocr.log.Info("Processing file", zap.String("fileName", f.Name()))
			err = ocr.ProcessByMetadataFile(ctx, filepath.Join(folder, f.Name()))
			if err != nil {
				ocr.log.Error("Error processing file", zap.String("fileName", f.Name()), zap.Error(err))
			}
		}
	}

	return nil
}

func (ocr *OcrService) ProcessByMetadataFile(ctx context.Context, fileName string) error {
	page, err := model.ReadMetadata(fileName)
	if err != nil {
		return err
	}

	ocr.log.Info("Processing file", zap.String("fileName", fileName))

	return ocr.Process(ctx, page.Filename, page.Id, page.Order)
}

func (ocr *OcrService) Process(ctx context.Context, fileName string, pageId string, pageNum int) error {
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
	md, err := ocr.generateMetadata(fileName, ta.GetPages()[0].Blocks, pageId, pageNum)
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

func (ocr *OcrService) isValidLang(langs []*visionpb.TextAnnotation_DetectedLanguage) bool {
	if langs == nil {
		return true // no language detected
	}

	for _, lang := range langs {
		if lang != nil {
			if lang.GetLanguageCode() == "he" ||
				lang.GetLanguageCode() == "en" ||
				lang.GetLanguageCode() == "iw" ||
				lang.GetLanguageCode() == "yi" {
				return true
			}
		}
	}

	return false

}
