package yizk

import (
	"context"
	"fmt"
	"log"

	"github.com/gipsh/yizk/translator"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var translateCmd = &cobra.Command{
	Use:   "translate",
	Short: "Run TRANSLATE batch processing",
	Long: `Run TRANSLATE batch processing on a book. 
	Looks into the download folder and for every metadata page json file and translate the content then update the file`,
	//	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)

		ctx := context.Background()

		if metadataFile != "" {
			err := executeTranslationByMetadataFile(ctx, metadataFile)
			if err != nil {
				log.Fatalln("Error rendering file", err.Error())
			}
			return
		}

		if metadataFolder != "" {
			err := executeTranslationByMetadataFolder(ctx, metadataFolder)
			if err != nil {
				log.Fatalln("Error rendering folder", err.Error())
			}
			return
		}

		fmt.Println("Done")
	},
}

func init() {
	rootCmd.AddCommand(translateCmd)
	translateCmd.Flags().StringVarP(&metadataFile, "file", "m", "", "fullpath of json metadata file. If provided, will only process this file")
	translateCmd.Flags().StringVarP(&metadataFolder, "folder", "f", "", "fullpath of folder with json metadata files. If provided, will only process this folder")
	translateCmd.MarkFlagsMutuallyExclusive("file", "folder")

}

func getTranslatorService(ctx context.Context) (*translator.MetadataTranslator, error) {

	log, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	translatorService, err := translator.NewTranslatorService(log, ctx, "creds.json")
	if err != nil {
		return nil, err
	}

	metadataTranslatorService := translator.NewMetadataTranslator(log, translatorService)

	return metadataTranslatorService, nil

}

func executeTranslationByMetadataFile(ctx context.Context, file string) error {

	metadataTranslatorService, err := getTranslatorService(ctx)
	if err != nil {
		return err
	}

	err = metadataTranslatorService.TranslatePage(ctx, file)
	if err != nil {
		return err
	}

	return nil

}

func executeTranslationByMetadataFolder(ctx context.Context, folder string) error {

	metadataTranslatorService, err := getTranslatorService(ctx)
	if err != nil {
		return err
	}

	err = metadataTranslatorService.TranslateFolder(ctx, folder)
	if err != nil {
		return err
	}

	return nil
}
