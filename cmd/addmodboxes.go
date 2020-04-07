package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// addmodboxesCmd represents the addmodboxes command
var addmodboxesCmd = &cobra.Command{
	Use:   "addmodboxes",
	Short: "A brief description of your command",
	Long: `Adds acroforms moderating boxes around the outside of all pages
in the input pdf, flattening the original pages to images so that
annotations can't be modified by the moderator`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(os.Args) < 2 {
			fmt.Printf("Requires one argument: input_path\n")
			fmt.Printf("Usage: gradex-coverpage.exe input.pdf\n")
			os.Exit(0)
		}

		suffix := filepath.Ext(inputPath)

		// sanity check
		if suffix != ".pdf" {
			fmt.Printf("Error: input path must be a .pdf\n")
			os.Exit(1)
		}

		// need page count to find the jpeg files again later
		numPages, err := countPages(inputPath)

		// render to images
		jpegPath := "./jpg"
		err = ensureDir(jpegPath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		basename := strings.TrimSuffix(inputPath, suffix)
		jpegFileOption := fmt.Sprintf("%s/%s%%04d.jpg", jpegPath, basename)

		err = convertPDFToJPEGs(inputPath, jpegPath, jpegFileOption)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		// convert images to individual pdfs, with form overlay

		pagePath := "./pdf"
		err = ensureDir(pagePath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		pageFileOption := fmt.Sprintf("%s/%s%%04d.pdf", pagePath, basename)
		formNameOption := fmt.Sprintf("%s%%04d", basename)

		mergePaths := []string{}

		// gs starts indexing at 1
		for imgIdx := 1; imgIdx <= numPages; imgIdx = imgIdx + 1 {

			// construct image name
			jpegFilename := fmt.Sprintf(jpegFileOption, imgIdx)
			pageFilename := fmt.Sprintf(pageFileOption, imgIdx)
			formID := fmt.Sprintf(formNameOption, imgIdx)

			// do the overlay
			convertJPEGToOverlaidPDFMod(jpegFilename, pageFilename, formID)

			//save the pdf filename for the merge at the end
			mergePaths = append(mergePaths, pageFilename)

		}

		if outputPath == "" {
			outputPath = fmt.Sprintf("%s-mod.pdf", basename)
		}
		err = mergePdf(mergePaths, outputPath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(addmodboxesCmd)

	// Here you will define your flags and configuration settings.

	addmodboxesCmd.Flags().StringVarP(&inputPath, "input", "i", "", "input path")
	addmodboxesCmd.Flags().StringVarP(&outputPath, "output", "o", "", "output path")
}
