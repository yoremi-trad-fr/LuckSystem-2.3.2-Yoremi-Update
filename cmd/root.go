/*
Copyright © 2022 WeTor wetorx@qq.com
*/
package cmd

import (
	"flag"
	"os"
	"strconv"

	"github.com/go-restruct/restruct"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "LuckSystem",
	Version: "2.3.2-yoremi.3.31",
	Short:   "LucaSystem引擎工具集",
	Long: `LucaSystem引擎工具集
https://github.com/wetor/LuckSystem
wetor(wetorx@qq.com)`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		configureLogging()
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}

func configureLogging() {
	if Log {
		_ = flag.Set("alsologtostderr", "true")
		_ = flag.Set("log_dir", LogDir)
		_ = flag.Set("v", strconv.Itoa(LogLevel))
	}
	// Cobra parses its own pflag set, not Go's standard flag set used by glog.
	// Mark the latter as parsed without feeding it Cobra's arguments; otherwise
	// every glog call is prefixed with "ERROR: logging before flag.Parse".
	if !flag.Parsed() {
		_ = flag.CommandLine.Parse([]string{})
	}
}

var (
	Log      bool
	LogLevel int
	LogDir   string
)

func init() {
	if len(os.Args) <= 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		_ = rootCmd.Help()
		os.Exit(1)
	}

	restruct.EnableExprBeta()

	rootCmd.PersistentFlags().BoolVar(&Log, "log", true, "启用日志")
	rootCmd.PersistentFlags().IntVar(&LogLevel, "log_level", 5, "输出日志等级")
	rootCmd.PersistentFlags().StringVar(&LogDir, "log_dir", "log", "保存日志路径")
}
