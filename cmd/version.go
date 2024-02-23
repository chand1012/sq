/*
Copyright © 2024 TimeSurgeLabs <chandler@timesurgelabs.com>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var commitHash string
var buildDate string
var tag string

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints version information.",
	Long:  `Prints version information.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sq")
		fmt.Println("commit hash:", commitHash)
		fmt.Println("build date:", strings.ReplaceAll(buildDate, "_", " "))
		fmt.Println("version:", tag)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
