/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"docmd/pkgs/droplets"
	"github.com/spf13/cobra"
)

// dropletCmd represents the droplet command
var dropletListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all droplets",
	Long:  `List all droplets in the account. This command takes no arguments.`,
	Run: func(cmd *cobra.Command, args []string) {
		droplets.ListDroplets()
	},
}

func init() {
	rootCmd.AddCommand(dropletListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dropletCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dropletCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
