/*
Copyright © 2020 Tim Drysdale <timothy.d.drysdale@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	inputPath, outputPath string
)

// addmarkboxesCmd represents the addmarkboxes command
var addmarkboxesCmd = &cobra.Command{
	Use:   "addmarkboxes",
	Short: "A brief description of your command",
	Long: `Adds acroforms marking boxes around the outside of all pages
in the input pdf, flattening the original pages to images so that
annotations can't be modified by the marker`,
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
			convertJPEGToOverlaidPDF(jpegFilename, pageFilename, formID)

			//save the pdf filename for the merge at the end
			mergePaths = append(mergePaths, pageFilename)

		}

		if outputPath == "" {
			outputPath = fmt.Sprintf("%s-mark.pdf", basename)
		}
		err = mergePdf(mergePaths, outputPath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(addmarkboxesCmd)

	// Here you will define your flags and configuration settings.

	addmarkboxesCmd.Flags().StringVarP(&inputPath, "input", "i", "", "input path")
	addmarkboxesCmd.Flags().StringVarP(&outputPath, "output", "o", "", "output path")
}
