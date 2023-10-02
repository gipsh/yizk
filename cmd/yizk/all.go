package yizk

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Run OCR, TRANSLATE and REDNER processing.",
	Long:  `Run OCR, TRANSLATE and REDNER processing.`,
	//	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		ctx := context.Background()

		if metadataFile != "" {
			executeOcrByMetadataFile(ctx, metadataFile)
			executeTranslationByMetadataFile(ctx, metadataFile)
			executeRenderByMetadataFile(metadataFile)
			return
		}

		if metadataFolder != "" {
			executeOcrByMetadataFolder(ctx, metadataFolder)
			executeTranslationByMetadataFolder(ctx, metadataFolder)
			executeRenderByMetadataFolder(metadataFolder)
			return
		}

		fmt.Println("Done")
	},
}

func init() {
	rootCmd.AddCommand(allCmd)
	allCmd.Flags().StringVarP(&metadataFile, "file", "m", "", "fullpath of json metadata file. If provided, will only process this file")
	allCmd.Flags().StringVarP(&metadataFolder, "folder", "f", "", "fullpath of folder with json metadata files. If provided, will only process this folder")
}
