/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("download called")
	},
}

var (
	downloadName    string
	shouldExtract   bool
	outputDirectory string
)

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringVarP(&downloadName, "name", "n", "", "Name of distribution to download")
	downloadCmd.MarkFlagRequired("name")

	downloadCmd.Flags().BoolVarP(&shouldExtract, "extract", "x", false, "Should archive after download be extracted")
	downloadCmd.Flags().StringVarP(&outputDirectory, "output", "o", "", "Output directory")
}
