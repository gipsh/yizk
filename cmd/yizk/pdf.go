package yizk

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/cli"
	"github.com/spf13/cobra"
)

var pdfCmd = &cobra.Command{
	Use:   "pdf",
	Short: "Build pdf from images",
	Long:  `Build pdf from images`,
	//	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		if bookName == "" {
			panic("book name is required")
		}

		buildPDF(bookName)

	},
}

func init() {
	rootCmd.AddCommand(pdfCmd)
	pdfCmd.Flags().StringVarP(&bookName, "book", "b", "", "path to the book folder with downloaded images")
	pdfCmd.MarkFlagRequired("book")

}

func buildPDF(bookName string) error {

	var pages []string

	var folder = "tmp/" + bookName
	files, err := os.ReadDir(folder)
	if err != nil {
		return err
	}

	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".render.png") {
			pages = append(pages, filepath.Join(folder, f.Name()))
		}
	}

	importCmd := cli.ImportImagesCommand(pages, "out.pdf", nil, nil)
	_, err = cli.Process(importCmd)
	if err != nil {
		return err
	}

	return nil

}
