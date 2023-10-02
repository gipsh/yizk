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

		if metadataFile != "" {
			fmt.Println("hola")
			err := executeRenderByMetadataFile(metadataFile)
			if err != nil {
				log.Fatalln("Error rendering file", metadataFile)
			}
			return
		}

		if metadataFolder != "" {
			err := executeRenderByMetadataFolder(metadataFolder)
			if err != nil {
				log.Fatalln("Error rendering folder", metadataFolder)
			}
			return
		}

		fmt.Println("Done")
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)
	renderCmd.Flags().StringVarP(&metadataFile, "file", "m", "", "fullpath of json metadata file. If provided, will only process this file")
	renderCmd.Flags().StringVarP(&metadataFolder, "folder", "f", "", "fullpath of folder with json metadata files. If provided, will only process this folder")
	renderCmd.MarkFlagsMutuallyExclusive("file", "folder")
}

func getRenderService() *renderer.MetadataRenderer {

	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	renderService := renderer.NewMetadataRenderer(log)

	return renderService
}

func executeRenderByMetadataFile(metadataFile string) error {

	renderService := getRenderService()

	err := renderService.RenderPage(metadataFile)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func executeRenderByMetadataFolder(metadataFolder string) error {

	renderService := getRenderService()

	err := renderService.RenderFolder(metadataFolder)
	if err != nil {
		return err
	}

	return nil
}
