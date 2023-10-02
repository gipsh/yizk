package downloader

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gipsh/yizk/model"
	"github.com/gipsh/yizk/ocr"
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

	for i, idx := pd.initialPage, 1; i <= pd.initialPage+pd.pages; i, idx = i+1, idx+1 {
		err := pd.DownloadPage(i, idx)
		if err != nil {
			pd.log.Info("Error downloading page", zap.Int("page", i), zap.Error(err))
		}
	}
}

func (pd *PageDownloader) DownloadPage(pageIdx int, pageNum int) error {
	pd.log.Info("Downloading page", zap.Int("page", pageIdx))

	filepath := "tmp/" + pd.bookName + "/" + fmt.Sprintf("%d.jpg", pageIdx)

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

	// create metadata file
	return pd.InitMetadataFile(pageIdx, pageNum, filepath)
}

func (pd *PageDownloader) InitMetadataFile(pageIdx int, pageNumber int, imageFile string) error {

	filepath := "tmp/" + pd.bookName + "/" + fmt.Sprintf("%d.json", pageIdx)
	metadata := model.YizkPage{
		Order:    pageNumber,
		Id:       fmt.Sprintf("%d", pageIdx),
		Filename: imageFile,
	}

	return ocr.WriteMetadataFile(filepath, &metadata)

}
