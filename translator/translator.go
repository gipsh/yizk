package translator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cloud.google.com/go/translate"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

type TranslatorService struct {
	log    *zap.Logger
	client *translate.Client
}

func NewTranslatorService(log *zap.Logger, ctx context.Context, credFile string) (*TranslatorService, error) {

	client, err := translate.NewClient(ctx, option.WithCredentialsFile(credFile))
	if err != nil {
		return nil, err
	}

	return &TranslatorService{
		log:    log,
		client: client,
	}, nil
}

func (ts *TranslatorService) Process(ctx context.Context, fileName string) error {

	file, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	var text string

	if len(file) > 2000 {
		// need to split by line breaks
		lines := strings.Split(string(file), "\n")
		//grouping lines by chunks of 2000 chars
		ts.log.Info("Splitting text", zap.Int("lines", len(lines)))
		var sb strings.Builder
		var outsb strings.Builder
		for _, line := range lines {
			if len(line) > 0 {
				sb.Write([]byte(line))
				sb.Write([]byte("\r\n"))
				if sb.Len() > 1600 {
					translatedLine, err := ts.Translate(ctx, line)
					if err != nil {
						return err
					}
					outsb.Write([]byte(translatedLine))
					outsb.Write([]byte("\r\n"))
					sb.Reset()
				}
			}
		}

		if sb.Len() > 0 {
			translatedLine, err := ts.Translate(ctx, sb.String())
			if err != nil {
				return err
			}
			outsb.Write([]byte(translatedLine))
			outsb.Write([]byte("\r\n"))
		}

		text = outsb.String()
	} else {
		text, err = ts.Translate(ctx, string(file))
		if err != nil {
			return err
		}

	}

	fileNameOnly := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	newFilename := fileNameOnly + ".translated" + ".txt"
	// save hebrew text to file
	ts.log.Info("Saving text to file", zap.String("fileName", newFilename))
	err = ts.saveToFile(newFilename, text)
	if err != nil {
		return err
	}

	return nil

}

func (ts *TranslatorService) Translate(ctx context.Context, text string) (string, error) {
	ts.log.Debug("Translating text", zap.String("text", text))

	// Use the client.
	t, err := ts.client.Translate(ctx, []string{text}, language.English, nil)
	if err != nil {
		return "", err
	}

	if len(t) == 0 {
		return "", fmt.Errorf("no translations found")
	}

	ts.log.Info("Translated text", zap.Int("translations", len(t)))
	ts.log.Info("Detected lang", zap.String("text", t[0].Source.String()))

	// Close the client when finished.
	if err := ts.client.Close(); err != nil {
		return "", err
	}

	var sb strings.Builder

	for _, translation := range t {
		ts.log.Debug("Translated text", zap.String("text", translation.Text))
		sb.Write([]byte(translation.Text))
		sb.Write([]byte("\n"))
	}

	return sb.String(), nil
}

func (ts *TranslatorService) TranslateWithQuota(ctx context.Context, text string) (string, error) {
	ts.log.Debug("Translating text", zap.String("text", text))

	// Use the client.
	t, err := ts.client.Translate(ctx, []string{text}, language.English, nil)
	if err != nil {
		return "", err
	}

	if len(t) == 0 {
		return "", fmt.Errorf("no translations found")
	}

	ts.log.Info("Translated text", zap.Int("translations", len(t)))
	ts.log.Info("Detected lang", zap.String("text", t[0].Source.String()))

	// Close the client when finished.
	if err := ts.client.Close(); err != nil {
		return "", err
	}

	var sb strings.Builder

	for _, translation := range t {
		ts.log.Debug("Translated text", zap.String("text", translation.Text))
		sb.Write([]byte(translation.Text))
		sb.Write([]byte("\n"))
	}

	return sb.String(), nil
}

func (ts *TranslatorService) saveToFile(fileName string, data string) error {
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
