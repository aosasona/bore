package main

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "bore",
	Short: "Bore is a clipboard for headless (and non-headless) environments",
	Long:  "Bore provides you with a clipboard and a clipboard manager on any machine you are working on",
}

func Execute() error {
	// TODO: add commands
	return rootCmd.Execute()
}
