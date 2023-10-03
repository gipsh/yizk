package yizk

import (
	"context"
	"fmt"

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

		ctx := context.Background()

		if metadataFile != "" {
			executeOcrByMetadataFile(ctx, metadataFile)
			return
		}

		if metadataFolder != "" {
			executeOcrByMetadataFolder(ctx, metadataFolder)
			return
		}

		fmt.Println("Done")
	},
}

func init() {
	rootCmd.AddCommand(ocrCmd)
	ocrCmd.Flags().StringVarP(&metadataFile, "file", "m", "", "fullpath of json metadata file. If provided, will only process this file")
	ocrCmd.Flags().StringVarP(&metadataFolder, "folder", "f", "", "fullpath of folder with json metadata files. If provided, will only process this folder")

}

func getOcrService(ctx context.Context) *ocr.OcrService {
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	ocrService, err := ocr.NewOcrService(ctx, log, "creds.json")
	if err != nil {
		panic(err)
	}

	return ocrService
}

func executeOcrByMetadataFile(ctx context.Context, file string) error {

	ocrService := getOcrService(ctx)
	return ocrService.ProcessByMetadataFile(ctx, file)

}

func executeOcrByMetadataFolder(ctx context.Context, folder string) error {

	ocrService := getOcrService(ctx)
	return ocrService.ProcessByFolder(ctx, folder)

}
