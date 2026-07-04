/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/golang/glog"
	"lucksystem/font"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// fontEditCmd represents the fontEdit command
var fontEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "修改或重构字体",
	Long:  `同时修改字体图像以及对应的info`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("fontEdit called")
		if len(FontTTFInput) == 0 {
			fmt.Println("Error: required flag(s) \"input_ttf\" not set")
			return
		}
		f := font.LoadLucaFontFile(FontInfoSource, FontCzSource)
		out, err := os.Create(FontOutput)
		if err != nil {
			glog.Fatalln(err)
		}
		defer out.Close()
		ttf, err := os.Open(FontTTFInput)
		if err != nil {
			glog.Fatalln(err)
		}
		defer ttf.Close()
		if FontAppend {
			FontStartIndex = -1
		}
		err = f.Import(ttf, FontStartIndex, FontRedraw, FontCharsetInput)
		if err != nil {
			glog.Fatalln(err)
		}
		metricRunes, err := fontEditMetricRunes(f, FontCharsetInput, FontRedraw)
		if err != nil {
			glog.Fatalln(err)
		}
		if FontArabicMetrics {
			arabicRunes := filterArabicMetricRunes(metricRunes)
			if len(arabicRunes) > 0 {
				metricOptions := font.MetricAdjustOptions{WOffset: -1}
				latinY, latinOK := fontEditReferenceY(f.Info, []rune{'a', 'o', 'A', 'O'})
				arabicY, arabicOK := fontEditCommonY(f.Info, arabicRunes)
				if latinOK && arabicOK {
					metricOptions.YOffset = latinY - arabicY
				}
				f.Info.AdjustMetrics(arabicRunes, metricOptions)
			}
		}
		metricOptions := font.MetricAdjustOptions{
			SetY:    cmd.Flags().Changed("metric-set-y"),
			Y:       FontMetricSetY,
			YOffset: FontMetricYOffset,
			XOffset: FontMetricXOffset,
			WOffset: FontMetricWOffset,
		}
		f.Info.AdjustMetrics(metricRunes, metricOptions)
		effectiveArabicConnectorBleed := fontEditEffectiveArabicConnectorBleed(
			FontArabicMetrics,
			FontArabicConnectorBleed,
			cmd.Flags().Changed("arabic-connector-bleed"),
		)
		if effectiveArabicConnectorBleed > 0 {
			arabicRunes := filterArabicMetricRunes(metricRunes)
			f.BleedGlyphEdges(arabicRunes, effectiveArabicConnectorBleed)
		}
		var outInfo *os.File = nil
		if len(FontInfoOutput) > 0 {
			outInfo, err = os.Create(FontInfoOutput)
			if err != nil {
				glog.Fatalln(err)
			}
		}
		err = f.Write(out, outInfo)
		if err != nil {
			glog.Fatalln(err)
		}

	},
}
var (
	FontTTFInput             string // ttf字体文件
	FontCharsetInput         string // 替换或追加的字符集
	FontRedraw               bool   // 重绘
	FontAppend               bool   // 追加到最后
	FontStartIndex           int    // 替换或者重绘的序号，从零开始
	FontArabicMetrics        bool
	FontMetricSetY           int
	FontMetricYOffset        int
	FontMetricXOffset        int
	FontMetricWOffset        int
	FontArabicConnectorBleed int
)

const defaultArabicConnectorBleed = 0

func fontEditEffectiveArabicConnectorBleed(arabicMetrics bool, bleed int, bleedChanged bool) int {
	if bleedChanged {
		return bleed
	}
	if arabicMetrics {
		return defaultArabicConnectorBleed
	}
	return bleed
}

func fontEditMetricRunes(f *font.LucaFont, charsetFile string, redraw bool) ([]rune, error) {
	if len(charsetFile) > 0 {
		data, err := os.ReadFile(charsetFile)
		if err != nil {
			return nil, err
		}
		chars := strings.TrimPrefix(string(data), "\ufeff")
		return []rune(chars), nil
	}
	if !redraw || f == nil || f.Info == nil {
		return nil, nil
	}
	chars := make([]rune, 0, len(f.Info.IndexUnicode))
	for _, char := range f.Info.IndexUnicode {
		if char == 0 {
			continue
		}
		chars = append(chars, char)
	}
	return chars, nil
}

func filterArabicMetricRunes(chars []rune) []rune {
	out := make([]rune, 0, len(chars))
	seen := make(map[rune]bool, len(chars))
	for _, char := range chars {
		if seen[char] || !isArabicMetricRune(char) {
			continue
		}
		seen[char] = true
		out = append(out, char)
	}
	return out
}

func isArabicMetricRune(char rune) bool {
	return (char >= 0x0600 && char <= 0x06FF) ||
		(char >= 0x0750 && char <= 0x077F) ||
		(char >= 0x0870 && char <= 0x089F) ||
		(char >= 0x08A0 && char <= 0x08FF) ||
		(char >= 0xFB50 && char <= 0xFDFF) ||
		(char >= 0xFE70 && char <= 0xFEFF)
}

func fontEditReferenceY(info *font.Info, candidates []rune) (int, bool) {
	if info == nil {
		return 0, false
	}
	for _, char := range candidates {
		if char < 0 || int(char) >= len(info.UnicodeIndex) {
			continue
		}
		index := info.UnicodeIndex[int(char)]
		if index == 0 && char != ' ' {
			continue
		}
		if int(index) >= len(info.DrawSize) {
			continue
		}
		return int(int8(info.DrawSize[int(index)].Y)), true
	}
	return 0, false
}

func fontEditCommonY(info *font.Info, chars []rune) (int, bool) {
	if info == nil {
		return 0, false
	}
	counts := make(map[int]int)
	order := make([]int, 0)
	seen := make(map[uint16]bool, len(chars))
	for _, char := range chars {
		if char < 0 || int(char) >= len(info.UnicodeIndex) {
			continue
		}
		index := info.UnicodeIndex[int(char)]
		if index == 0 && char != ' ' {
			continue
		}
		if int(index) >= len(info.DrawSize) || seen[index] {
			continue
		}
		seen[index] = true
		y := int(int8(info.DrawSize[int(index)].Y))
		if counts[y] == 0 {
			order = append(order, y)
		}
		counts[y]++
	}
	if len(order) == 0 {
		return 0, false
	}
	best := order[0]
	for _, y := range order[1:] {
		if counts[y] > counts[best] {
			best = y
		}
	}
	return best, true
}

func init() {
	fontCmd.AddCommand(fontEditCmd)
	fontEditCmd.Flags().StringVarP(&FontInfoOutput, "output_info", "O", "", "修改后字体info保存位置")

	fontEditCmd.Flags().StringVarP(&FontTTFInput, "input_ttf", "f", "", "绘制字符使用的TTF字体")
	fontEditCmd.Flags().StringVarP(&FontCharsetInput, "input_charset", "c", "", "增加或替换的字符集文本文件")
	fontEditCmd.Flags().BoolVarP(&FontAppend, "append", "a", false, "字符集绘制并添加到原字体最后")

	fontEditCmd.Flags().IntVarP(&FontStartIndex, "index", "i", 0, "字符集绘制并添加到的位置，从0开始")
	fontEditCmd.Flags().BoolVarP(&FontRedraw, "redraw", "r", false, "重绘原字体图片")
	fontEditCmd.Flags().BoolVar(&FontArabicMetrics, "arabic-metrics", false, "apply Arabic presentation-form metrics: shift toward Latin baseline and reduce advance by 1px")
	fontEditCmd.Flags().IntVar(&FontMetricSetY, "metric-set-y", 0, "set signed draw_y for edited glyphs")
	fontEditCmd.Flags().IntVar(&FontMetricYOffset, "metric-y-offset", 0, "add signed offset to draw_y for edited glyphs")
	fontEditCmd.Flags().IntVar(&FontMetricXOffset, "metric-x-offset", 0, "add signed offset to draw_x for edited glyphs")
	fontEditCmd.Flags().IntVar(&FontMetricWOffset, "metric-w-offset", 0, "add signed offset to character advance/usize_w for edited glyphs")
	fontEditCmd.Flags().IntVar(&FontArabicConnectorBleed, "arabic-connector-bleed", 0, "experimentally extend Arabic connector pixels by N px to reduce connector gaps")
	fontEditCmd.MarkFlagsMutuallyExclusive("append", "index")
}
