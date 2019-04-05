package clientcmd

import (
	"fmt"
	"os"

	"github.com/catherinetcai/gsuite-aws-sso/pkg/config"
	"github.com/catherinetcai/gsuite-aws-sso/version"
	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use: "gsuite-aws-sso",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Surprise! We're not using Viper (as is the Cobra default)
	cobra.OnInitialize(config.Initialize)
	rootCmd.AddCommand(version.VersionCmd)
}
