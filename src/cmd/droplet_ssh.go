/*
Copyright © 2024 Alejandro Cavallo <alejandro.cavallo@stoneitcloud.com>
*/
package cmd

import (
	"docmd/pkgs/droplets"

	"github.com/spf13/cobra"
)

// dropletCmd represents the droplet command
var dropletSshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "SSH into a droplet",
	Long:  `SSH into a droplet. This command takes no arguments.`,
	Run: func(cmd *cobra.Command, args []string) {
		droplets.SshDropletUi()
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
