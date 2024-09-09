/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"lib_reserve/lib"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lib_reserve",
	Short: "GZHU Library seat reserve, and more.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&lib.Un, "user", "u", "", "your student number")
	rootCmd.PersistentFlags().StringVarP(&lib.Pd, "passwd", "p", "", "your student password")
	rootCmd.PersistentFlags().StringVarP(&lib.Cookies, "cookie", "c", "", "your cookies")
	rootCmd.PersistentFlags().StringVar(&lib.Cookies_path, "cookie-path", lib.Cookies_path, "your cookies file path")
}
