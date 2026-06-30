package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var movieCmd = &cobra.Command{
	Use:   "movie",
	Short: "LucaSystem MVT movies",
	Long:  "LucaSystem MVT movie files containing embedded WebM data",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("movie called")
	},
}

var (
	MovieInput  string
	MovieOutput string
)

func init() {
	rootCmd.AddCommand(movieCmd)

	movieCmd.PersistentFlags().StringVarP(&MovieInput, "input", "i", "", "input MVT file")
	movieCmd.PersistentFlags().StringVarP(&MovieOutput, "output", "o", "", "output WebM file")

	movieCmd.MarkPersistentFlagRequired("input")
	movieCmd.MarkPersistentFlagRequired("output")
	movieCmd.MarkFlagsRequiredTogether("output", "input")
}
