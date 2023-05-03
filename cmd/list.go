/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"sort"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/soltys/gosqlitedist/internal"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		performList()
	},
}

func performList() {
	sqliteProducts := internal.MustParse()
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Extension", "Version"})
	sort.Slice(sqliteProducts, func(i, j int) bool {
		return strings.Compare(sqliteProducts[i].Name, sqliteProducts[j].Name) < 0
	})
	for _, product := range sqliteProducts {
		t.AppendRow(table.Row{
			product.Name, product.Extension, product.Version,
		})
	}
	t.Render()
}

func init() {
	rootCmd.AddCommand(listCmd)
}
