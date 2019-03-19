package clientcmd

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/briandowns/spinner"
	"github.com/catherinetcai/gsuite-aws-sso/pkg/client/config"
	"github.com/catherinetcai/gsuite-aws-sso/pkg/file"
	"github.com/catherinetcai/gsuite-aws-sso/pkg/logging"
	"github.com/catherinetcai/gsuite-aws-sso/pkg/shared/handlers"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var credential string

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log into AWS",
	Run:   login,
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.PersistentFlags().StringVarP(&credential, "credential", "c", defaultGCloudCredentialPath(), "Path to Google Cloud credentials file. Defaults to $HOME/.config/gcloud/application_default_credentials.json")
}

func defaultGCloudCredentialPath() string {
	path, err := file.WithUserHomeDir(".config", "gcloud", "application_default_credentials.json")
	if err != nil {
		logging.Logger().Error("error getting home path", zap.Error(err))
		return ""
	}

	return path
}

// TODO: Centralize these requests somewhere since they are shared between the server and client
func login(cmd *cobra.Command, args []string) {
	cfg, err := config.Get()
	if err != nil {
		logging.Logger().Fatal("config does not exist, try running config", zap.Error(err))
	}

	// Read GCP credentials in
	credentialFile, err := ioutil.ReadFile(cfg.GCP.CredentialFilePath)
	if err != nil {
		logging.Logger().Fatal("credentials file does not exist - you must login using GCloud", zap.Error(err))
	}

	// Before doing work, check the location of AWS file. If it exists, warn
	// TODO: Pull this into its own function
	exists, err := checkOutputCredentialsFileExist(cfg.AWS.CredentialOutputPath)
	if err != nil {
		logging.Logger().Fatal("something went wrong accessing the AWS credentials file", zap.Error(err))
	}
	if exists {
		overWritePrompt := promptui.Prompt{
			Label:     "AWS Credentials already exists. Overwrite?",
			IsConfirm: true,
		}

		result, err := overWritePrompt.Run()
		if err != nil {
			logging.Logger().Fatal("Something went wrong...")
		}

		if result == "N" {
			logging.Logger().Info("Remove AWS credentials before logging in...")
		}
	}

	logging.Logger().Info("Logging in...")
	s := spinner.New(spinner.CharSets[4], 100*time.Millisecond)
	s.Start()

	// Create request and send it
	// TODO: This client interaction should be put into a package
	req := &handlers.CredentialHandlerRequest{CredentialFile: credentialFile}
	client := &http.Client{}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		logging.Logger().Fatal("error marshaling req", zap.Error(err))
	}

	b := bytes.NewBuffer(reqBytes)

	resp, err := client.Post(cfg.Server, "application/json", b)
	if err != nil {
		logging.Logger().Fatal("error posting to credentials", zap.Error(err))
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		logging.Logger().Fatal("error reading credential response", zap.Error(err))
	}

	credentialResp := &handlers.CredentialHandlerResponse{}

	// Unmarshalling only to make sure that the response body is correct
	if err := json.Unmarshal(respBody, credentialResp); err != nil {
		logging.Logger().Fatal("error unmarshaling response", zap.Error(err))
	}

	err = writeCredentialsFile(cfg.AWS.CredentialOutputPath, credentialResp.CredentialFile)
	s.Stop()

	if err != nil {
		logging.Logger().Fatal("error writing credentials file", zap.Error(err))
	}
}

func checkOutputCredentialsFileExist(path string) (exists bool, err error) {
	_, err = os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
			return
		}
		return
	}

	return true, nil
}

func writeCredentialsFile(filePath string, body []byte) error {
	dirPath := path.Dir(filePath)
	if err := os.MkdirAll(dirPath, 0744); err != nil {
		return err
	}

	cfgFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	// Ensure file is empty
	if err := cfgFile.Truncate(0); err != nil {
		return err
	}
	cfgFile.Seek(0, 0)

	_, err = cfgFile.Write(body)
	return err
}
