package cmd

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"lucksystem/audio"
)

var audioConvertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert audio folder native Ogg <-> MP3",
	Run: func(cmd *cobra.Command, args []string) {
		if PakInput == "" || PakOutput == "" {
			fmt.Println("Error: required flag(s) \"input\" and \"output\" not set")
			return
		}
		fmt.Println("audio convert called")
		summary, err := audio.ConvertPath(audio.ConvertOptions{
			InputPath: PakInput,
			OutputDir: PakOutput,
			Direction: AudioDirection,
			Overwrite: true,
		}, func(line string) {
			fmt.Println(line)
		})
		if err != nil {
			glog.Fatalln(err)
		}
		printAudioSummary("Done", 0, summary.Converted, summary.Skipped, summary.Errors)
	},
}

func init() {
	audioCmd.AddCommand(audioConvertCmd)
	audioConvertCmd.Flags().StringVar(&AudioDirection, "to", "mp3", "conversion target: mp3 or native")
}
