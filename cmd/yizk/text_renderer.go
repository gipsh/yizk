package yizk

import (
	"fmt"
	"log"

	"github.com/gipsh/yizk/renderer"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var textCmd = &cobra.Command{
	Use:   "text",
	Short: "Run Text RENDER batch processing on a download folder",
	Long:  `Run Text RENDER batch processing on a download folder`,
	//	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		if metadataFile != "" {
			err := executeTextRenderByMetadataFile(metadataFile)
			if err != nil {
				log.Fatalln("Error rendering file", metadataFile)
			}
			return
		}

		if metadataFolder != "" {
			var output *string
			if outputTextFile == "" {
				output = nil
			} else {
				output = &outputTextFile
			}

			err := executeTextRenderByMetadataFolder(metadataFolder, output)
			if err != nil {
				log.Fatalln("Error rendering folder", metadataFolder)
			}
			return
		}

		fmt.Println("Done")
	},
}

func init() {
	rootCmd.AddCommand(textCmd)
	textCmd.Flags().StringVarP(&metadataFile, "file", "m", "", "fullpath of json metadata file. If provided, will only process this file")
	textCmd.Flags().StringVarP(&metadataFolder, "folder", "f", "", "fullpath of folder with json metadata files. If provided, will only process this folder")
	textCmd.Flags().StringVarP(&outputTextFile, "output", "o", "", "output text file")
	textCmd.MarkFlagsMutuallyExclusive("file", "folder")
}

func getTextRenderService() renderer.Renderer {

	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	renderService := renderer.NewTextRenderer(log)

	return renderService
}

func executeTextRenderByMetadataFile(metadataFile string) error {

	renderService := getTextRenderService()

	_, err := renderService.RenderPage(metadataFile)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func executeTextRenderByMetadataFolder(metadataFolder string, outputFile *string) error {

	renderService := getTextRenderService()

	err := renderService.RenderFolder(metadataFolder, outputFile)
	if err != nil {
		return err
	}

	return nil
}
