/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// dropletCmd represents the droplet command
var dropletCmd = &cobra.Command{
	Use:   "droplet",
	Short: "Interact with the account's droplets",
	Long:  `Interact with the account's droplets. You can list, ssh into, and trigger backups for droplets.`,
}

func init() {
	dropletCmd.AddCommand(dropletListCmd)
	dropletCmd.AddCommand(dropletSshCmd)
	rootCmd.AddCommand(dropletCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dropletCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dropletCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
