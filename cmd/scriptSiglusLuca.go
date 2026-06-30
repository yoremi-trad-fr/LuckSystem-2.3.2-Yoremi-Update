package cmd

import (
	"fmt"

	"lucksystem/siglusluca"

	"github.com/spf13/cobra"
)

var (
	siglusLucaLucaDir      string
	siglusLucaSiglusDir    string
	siglusLucaOutputDir    string
	siglusLucaHDOutput     string
	siglusLucaReviewOutput string
	siglusLucaTargetCol    int
	siglusLucaMinScore     float64
)

var scriptSiglusLucaCmd = &cobra.Command{
	Use:   "siglus-luca",
	Short: "Import Siglus translation text into decompiled Luca scripts",
	Long: `Import Siglus translation text into decompiled Luca scripts.

The Luca script folder remains the master structure. Matching Siglus
translation lines replace the selected Luca quoted string, while Luca-only lines
and merged/split lines are exported to TSV reports for manual review.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		summary, err := siglusluca.Run(siglusluca.Options{
			LucaDir:      siglusLucaLucaDir,
			SiglusDir:    siglusLucaSiglusDir,
			OutputDir:    siglusLucaOutputDir,
			HDOutput:     siglusLucaHDOutput,
			ReviewOutput: siglusLucaReviewOutput,
			TargetCol:    siglusLucaTargetCol,
			MinScore:     siglusLucaMinScore,
		})
		if err != nil {
			return err
		}

		fmt.Printf("Siglus -> Luca bridge complete\n")
		fmt.Printf("  files processed: %d\n", summary.FilesProcessed)
		fmt.Printf("  files copied unchanged: %d\n", summary.FilesCopied)
		fmt.Printf("  imported lines: %d\n", summary.Imported)
		fmt.Printf("  HD candidate lines: %d\n", summary.HDCandidates)
		fmt.Printf("  review rows: %d\n", summary.ReviewRows)
		fmt.Printf("  low-confidence imported/reviewed lines: %d\n", summary.LowConfidence)
		fmt.Printf("  output scripts: %s\n", siglusLucaOutputDir)
		if siglusLucaHDOutput != "" {
			fmt.Printf("  HD candidates: %s\n", siglusLucaHDOutput)
		}
		if siglusLucaReviewOutput != "" {
			fmt.Printf("  review report: %s\n", siglusLucaReviewOutput)
		}
		return nil
	},
}

func init() {
	scriptCmd.AddCommand(scriptSiglusLucaCmd)

	scriptSiglusLucaCmd.Flags().StringVar(&siglusLucaLucaDir, "luca", "", "decompiled Luca scripts directory")
	scriptSiglusLucaCmd.Flags().StringVar(&siglusLucaSiglusDir, "siglus", "", "Siglus Full directory containing .ss.txt files")
	scriptSiglusLucaCmd.Flags().StringVarP(&siglusLucaOutputDir, "output", "o", "", "output directory for patched Luca scripts")
	scriptSiglusLucaCmd.Flags().StringVar(&siglusLucaHDOutput, "hd-output", "", "TSV output for Luca-only HD candidate lines")
	scriptSiglusLucaCmd.Flags().StringVar(&siglusLucaReviewOutput, "review-output", "", "TSV output for low-confidence and split/merged lines")
	scriptSiglusLucaCmd.Flags().IntVar(&siglusLucaTargetCol, "target-col", 2, "1-based quoted string column to replace in Luca scripts")
	scriptSiglusLucaCmd.Flags().Float64Var(&siglusLucaMinScore, "min-score", 0, "minimum alignment score to import a matched line; 0 imports every aligned line")

	scriptSiglusLucaCmd.MarkFlagRequired("luca")
	scriptSiglusLucaCmd.MarkFlagRequired("siglus")
	scriptSiglusLucaCmd.MarkFlagRequired("output")
}
