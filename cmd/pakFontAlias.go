package cmd

import (
	"bytes"
	"fmt"
	"os"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"lucksystem/charset"
	"lucksystem/font"
	"lucksystem/pak"
)

var pakFontAliasCmd = &cobra.Command{
	Use:   "font-alias",
	Short: "copy one internal font-size entry into another while preserving target geometry",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pakFontAlias called")
		if PakSource == "" || PakOutput == "" || PakFontAliasFrom == "" || PakFontAliasTo == "" {
			fmt.Println("Error: source PAK, output PAK, from-name, and to-name are required")
			return
		}

		p := pak.LoadPak(PakSource, charset.Charset(Charset))
		source, err := p.Get(PakFontAliasFrom)
		if err != nil {
			glog.Fatalln(err)
		}
		target, err := p.Get(PakFontAliasTo)
		if err != nil {
			glog.Fatalln(err)
		}
		aliasData, err := font.BuildFontAliasEntry(PakFontAliasFrom, PakFontAliasTo, source.Data, target.Data)
		if err != nil {
			glog.Fatalln(err)
		}
		if err := p.Set(PakFontAliasTo, bytes.NewReader(aliasData)); err != nil {
			glog.Fatalln(err)
		}

		out, err := os.Create(PakOutput)
		if err != nil {
			glog.Fatalln(err)
		}
		defer out.Close()
		if err := p.Write(out); err != nil {
			glog.Fatalln(err)
		}
	},
}

var (
	PakFontAliasFrom string
	PakFontAliasTo   string
)

func init() {
	pakCmd.AddCommand(pakFontAliasCmd)
	pakFontAliasCmd.Flags().StringVar(&PakFontAliasFrom, "from-name", "", "source internal font entry name, for example info30 or 明朝30")
	pakFontAliasCmd.Flags().StringVar(&PakFontAliasTo, "to-name", "", "target internal font entry name, for example info32 or 明朝32")
}
