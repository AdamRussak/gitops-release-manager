/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"gitops-release-manager/pkg/core"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var tversion string
var (
	o       FlagsOptions
	version = tversion
	rootCmd = &cobra.Command{
		Use:   "gitops-version",
		Short: rootLongDisc,
		Long:  rootLongDisc,
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&core.Verbosity, "verbose", "v", false, "verbose logging")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
