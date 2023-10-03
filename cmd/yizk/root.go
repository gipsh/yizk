package yizk

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Parameters
var bookName string
var initialPage int
var pages int
var metadataFile string
var metadataFolder string

var rootCmd = &cobra.Command{
	Use:   "yizk",
	Short: "run the different services of the processing pipeline",
	Long: `yizk is a CLI tool to run the different services of the processing pipeline
	- DOWNLOAD: download the images from the NYPL website
	- OCR: process the image files and generate the metadata json files
	- TRANSLATE: translate the text from metadata json files and update the metadata json files
	- RENDER: render the translated text into the image files mimicking the original text

   
One can use yizk to translate yizkor books from the NYPL website.`,
	// Run: func(cmd *cobra.Command, args []string) {

	// },
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
