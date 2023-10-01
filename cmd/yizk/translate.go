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

		err := doTranslation(initialPage, pages, bookName)
		if err != nil {
			log.Fatalln("The timezone string is invalid")
		}
		fmt.Println("Done")
	},
}

func init() {
	rootCmd.AddCommand(translateCmd)
	translateCmd.Flags().StringVarP(&bookName, "book", "b", "", "path to the book folder with downloaded images")
	translateCmd.Flags().IntVarP(&initialPage, "initial-page", "i", 0, "initial page number to start processing")
	translateCmd.Flags().IntVarP(&pages, "pages", "p", 0, "number of pages to process")
	translateCmd.MarkFlagRequired("book")
	translateCmd.MarkFlagRequired("initial-page")
	translateCmd.MarkFlagRequired("pages")
}

func doTranslation(initialPage int, pages int, bookName string) error {

	ctx := context.Background()
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	translatorService, err := translator.NewTranslatorService(log, ctx, "creds.json")
	if err != nil {
		return err
	}

	metadataTranslatorService := translator.NewMetadataTranslator(log, translatorService)

	for i := initialPage; i <= initialPage+pages; i++ {
		err := metadataTranslatorService.TranslatePage(ctx, "tmp/"+bookName+"/"+fmt.Sprintf("%d.json", i), translatorService)
		if err != nil {
			log.Info("Error translating page", zap.Int("page", i), zap.Error(err))
		}
	}

	return nil

}
