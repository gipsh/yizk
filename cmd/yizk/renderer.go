package yizk

import (
	"fmt"
	"log"

	"github.com/gipsh/yizk/renderer"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Run RENDER batch processing on a download folder",
	Long: `Run RENDER batch processing on a book. 
	Looks into the translated metadata json file and using the image will try to redraw the page with the translated text.`,
	//	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)

		err := doRender(initialPage, pages, bookName)
		if err != nil {
			log.Fatalln("The timezone string is invalid")
		}
		fmt.Println("Done")
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)
	renderCmd.Flags().StringVarP(&bookName, "book", "b", "", "path to the book folder with downloaded images")
	renderCmd.Flags().IntVarP(&initialPage, "initial-page", "i", 0, "initial page number to start processing")
	renderCmd.Flags().IntVarP(&pages, "pages", "p", 0, "number of pages to process")
	renderCmd.MarkFlagRequired("book")
	renderCmd.MarkFlagRequired("initial-page")
	renderCmd.MarkFlagRequired("pages")
}

func doRender(initialPage int, pages int, bookName string) error {

	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	renderService := renderer.NewMetadataRenderer(log)

	for i := initialPage; i <= initialPage+pages; i++ {
		err := renderService.RenderPage("tmp/" + bookName + "/" + fmt.Sprintf("%d.json.json", i))
		if err != nil {
			log.Info("Error rendering page", zap.Int("page", i), zap.Error(err))
		}
	}

	return nil

}
