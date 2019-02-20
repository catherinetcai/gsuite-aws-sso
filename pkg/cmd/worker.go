package cmd

import "github.com/spf13/cobra"

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log into AWS",
	Run:   login,
}

func login(cmd *cobra.Command, args []string) {}
