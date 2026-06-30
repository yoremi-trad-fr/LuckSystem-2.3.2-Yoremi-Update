package cmd

import (
	"fmt"
	"lucksystem/movie"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

var movieExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Extract MVT movie to WebM",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("movieExport called")
		if err := movie.ExtractWebMFromMVT(MovieInput, MovieOutput); err != nil {
			glog.Fatalln(err)
		}
	},
}

func init() {
	movieCmd.AddCommand(movieExportCmd)
}
