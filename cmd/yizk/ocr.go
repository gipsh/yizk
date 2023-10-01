package yizk

import (
	"context"
	"fmt"
	"log"

	"github.com/gipsh/yizk/ocr"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var ocrCmd = &cobra.Command{
	Use:   "ocr",
	Short: "Run OCR batch processing on a download folder",
	Long: `Run OCR batch processing on a book. 
	Looks into the download folder and for every page tries to OCR the text and save it as a json file.`,
	//	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)

		err := doOCR(initialPage, pages, bookName)
		if err != nil {
			log.Fatalln("The timezone string is invalid")
		}
		fmt.Println("Done")
	},
}

var bookName string
var initialPage int
var pages int

func init() {
	rootCmd.AddCommand(ocrCmd)
	ocrCmd.Flags().StringVarP(&bookName, "book", "b", "", "path to the book folder with downloaded images")
	ocrCmd.Flags().IntVarP(&initialPage, "initial-page", "i", 0, "initial page number to start processing")
	ocrCmd.Flags().IntVarP(&pages, "pages", "p", 0, "number of pages to process")
	ocrCmd.MarkFlagRequired("book")
	ocrCmd.MarkFlagRequired("initial-page")
	ocrCmd.MarkFlagRequired("pages")
}

func doOCR(initialPage int, pages int, bookName string) error {

	ctx := context.Background()
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	ocrService, err := ocr.NewOcrService(ctx, log, "creds.json")
	if err != nil {
		panic(err)
	}

	for i := initialPage; i <= initialPage+pages; i++ {
		err := ocrService.Process(ctx, "tmp/"+bookName+"/"+fmt.Sprintf("%d.png", i), fmt.Sprintf("%d", i))
		if err != nil {
			log.Info("Error ocr'ing page", zap.Int("page", i), zap.Error(err))
		}
	}

	return nil

}
