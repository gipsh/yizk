package downloader

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"go.uber.org/zap"
)

type PageDownloader struct {
	log         *zap.Logger
	bookName    string
	initialPage int
	pages       int
}

func NewPageDownloader(log *zap.Logger, bookName string, initialPage int, pages int) *PageDownloader {
	return &PageDownloader{
		log:         log,
		bookName:    bookName,
		initialPage: initialPage,
		pages:       pages,
	}
}

func (pd *PageDownloader) Start() {

	// create tmp directory
	err := os.MkdirAll("tmp/"+pd.bookName, os.ModePerm)
	if err != nil {
		pd.log.Info("Error creating tmp directory", zap.Error(err))
		return
	}

	for i := pd.initialPage; i <= pd.initialPage+pd.pages; i++ {
		err := pd.DownloadPage(i)
		if err != nil {
			pd.log.Info("Error downloading page", zap.Int("page", i), zap.Error(err))
		}
	}
}

func (pd *PageDownloader) DownloadPage(pageIdx int) error {
	pd.log.Info("Downloading page", zap.Int("page", pageIdx))

	filepath := "tmp/" + pd.bookName + "/" + fmt.Sprintf("%d.png", pageIdx)

	// first check if file exists and skip
	if _, err := os.Stat(filepath); !errors.Is(err, os.ErrNotExist) {
		return os.ErrExist
	}

	url := fmt.Sprintf("https://images.nypl.org/index.php?id=%d&t=w", pageIdx)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)")

	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	//open a file for writing
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	return nil
}
