package yizk

import (
	"fmt"
	"log"

	"github.com/gipsh/yizk/downloader"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Run DOWNLOAD batch processing",
	Long: `Run DOWNLOAD batch processing to download the images of a book. 
	Iterate over the pages of the book and download the images`,
	//	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)

		err := doDownload(initialPage, pages, bookName)
		if err != nil {
			log.Fatalln("The timezone string is invalid")
		}
		fmt.Println("Done")
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringVarP(&bookName, "book", "b", "", "path to the book folder with downloaded images")
	downloadCmd.Flags().IntVarP(&initialPage, "initial-page", "i", 0, "initial page number to start processing")
	downloadCmd.Flags().IntVarP(&pages, "pages", "p", 0, "number of pages to process")
	downloadCmd.MarkFlagRequired("book")
	downloadCmd.MarkFlagRequired("initial-page")
	downloadCmd.MarkFlagRequired("pages")
}

func doDownload(initialPage int, pages int, bookName string) error {

	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	downloader := downloader.NewPageDownloader(log, bookName, initialPage, pages)

	for i := initialPage; i <= initialPage+pages; i++ {
		downloader.Start()
	}

	return nil

}
