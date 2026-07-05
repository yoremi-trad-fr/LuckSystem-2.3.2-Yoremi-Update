package cmd

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"lucksystem/audio"
)

var audioMusicExtractCmd = &cobra.Command{
	Use:   "music-extract",
	Short: "Extract MUSIC.PAK to native Ogg files",
	Run: func(cmd *cobra.Command, args []string) {
		runAudioExtract("music")
	},
}

var audioVoiceExtractCmd = &cobra.Command{
	Use:   "voice-extract",
	Short: "Extract VOICE/SYSVOICE PAK to native Ogg files",
	Run: func(cmd *cobra.Command, args []string) {
		runAudioExtract("voice")
	},
}

func init() {
	audioCmd.AddCommand(audioMusicExtractCmd)
	audioCmd.AddCommand(audioVoiceExtractCmd)

	for _, c := range []*cobra.Command{audioMusicExtractCmd, audioVoiceExtractCmd} {
		c.Flags().BoolVar(&AudioMP3, "mp3", false, "also convert extracted native Ogg files to MP3")
	}
}

func runAudioExtract(kind string) {
	if PakInput == "" || PakOutput == "" {
		fmt.Println("Error: required flag(s) \"input\" and \"output\" not set")
		return
	}
	fmt.Printf("audio %s extract called\n", kind)
	summary, err := audio.ExtractPak(audio.ExtractOptions{
		PakFile:    PakInput,
		OutputDir:  PakOutput,
		Kind:       kind,
		ConvertMP3: AudioMP3,
	}, func(line string) {
		fmt.Println(line)
	})
	if err != nil {
		glog.Fatalln(err)
	}
	fmt.Println("List:", summary.ListFile)
	fmt.Println("Native:", summary.NativeDir)
	if summary.MP3Dir != "" {
		fmt.Println("MP3:", summary.MP3Dir)
	}
	printAudioSummary("Done", summary.Files, summary.Converted, summary.Skipped, summary.Errors)
}
