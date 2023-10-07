package translator

import (
	"context"
	"fmt"
	"html"
	"os"
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

func (ts *TranslatorService) SetLogger(log zap.Logger) {
	ts.log = &log
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

func (ts *TranslatorService) Translate(ctx context.Context, text string) (string, error) {
	ts.log.Debug("Translating text", zap.String("text", text))

	// Preprocess text to split it by lines if needed
	input := ts.PreProcessText(text)

	var t []translate.Translation
	var err error

	// Do the translation
	if len(input) == 1 {
		t, err = ts.client.Translate(ctx, input, language.English, nil)
		if err != nil {
			return "", err
		}
	} else {
		// do multiple calls to the API to bypass the 2000 chars limit
		for _, line := range input {
			l, err := ts.client.Translate(ctx, []string{line}, language.English, nil)
			if err != nil {
				return "", err
			}
			t = append(t, l...)
		}
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

	//ts.log.Debug("Text", zap.String("text", html.UnescapeString(sb.String())))

	return html.UnescapeString(sb.String()), nil
}

// if text is grater than 2000 chars then split it by lines and translate each line
func (ts *TranslatorService) PreProcessText(text string) []string {
	var lines []string
	if len(text) > 2000 {
		// need to split by line breaks
		lines = strings.Split(text, ".")
	} else {
		lines = append(lines, text+".")
	}

	return lines
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
