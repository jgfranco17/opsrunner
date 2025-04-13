package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func GenerateDocs(cmd *cobra.Command, outputDir string) error {
	// Ensure output directory exists
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Failed to create docs output directory: %w", err)
	}

	// Generate the Markdown documentation
	err = doc.GenMarkdownTree(cmd, outputDir)
	if err != nil {
		return fmt.Errorf("Failed to generate markdown docs: %w", err)
	}

	fmt.Printf("📘 CLI documentation generated in: %s\n", filepath.Clean(outputDir))
	return nil
}
