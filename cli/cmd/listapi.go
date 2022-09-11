package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

var ListAPICmd = &cobra.Command{
	Use:   "list-api",
	Short: "List all available APIs",
	Run: func(cmd *cobra.Command, args []string) {
		apis := generateValidArgs()

		sort.Strings(apis)

		for i, api := range apis {
			fmt.Printf("%d. %s\n", i+1, api)
		}
	},
}
