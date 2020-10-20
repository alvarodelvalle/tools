package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	_ "github.com/spf13/viper"
	"os"
)

var rootCmd = &cobra.Command{Use: "tools"}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(ecrCmd)
	rootCmd.AddCommand(gitlabCmd)
	gitlabCmd.AddCommand(treeSearchCmd)
}
