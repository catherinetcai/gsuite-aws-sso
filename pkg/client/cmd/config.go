package clientcmd

import (
	"fmt"

	"github.com/catherinetcai/gsuite-aws-sso/pkg/client/config"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	errInvalidPath = fmt.Errorf("Path cannot be blank")
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure the client CMD",
	Run:   doConfigCmd,
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func doConfigCmd(cmd *cobra.Command, args []string) {
	cfg := prompt()

	fmt.Println("Setting configs...")
	err := config.SetConfigs(cfg)
	if err != nil {
		panic(err)
	}
}

func prompt() *config.Config {
	cfg := config.Default()

	serverPrompt := promptui.Prompt{
		Label:    "Login Server",
		Default:  config.DefaultServerURL,
		Validate: validatePath,
	}
	serverPromptResult, err := serverPrompt.Run()
	if err != nil {
		panic(err)
	}
	cfg.Server = serverPromptResult

	gcloudPrompt := promptui.Prompt{
		Label:    "GCloud Credentials",
		Default:  config.DefaultGCloudCredentialPath(),
		Validate: validatePath,
	}
	gcloudPromptResult, err := gcloudPrompt.Run()
	if err != nil {
		panic(err)
	}

	cfg.GCP.CredentialFilePath = gcloudPromptResult

	awsPrompt := promptui.Prompt{
		Label:    "AWS Output Credentials",
		Default:  config.DefaultAWSOutputPath(),
		Validate: validatePath,
	}
	awsPromptResult, err := awsPrompt.Run()
	if err != nil {
		panic(err)
	}

	cfg.AWS.CredentialOutputPath = awsPromptResult

	return cfg
}

func validatePath(input string) error {
	if input == "" {
		return errInvalidPath
	}
	return nil
}
