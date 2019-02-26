package cmd

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/catherinetcai/gsuite-aws-sso/pkg/logging"
	"github.com/davecgh/go-spew/spew"
	homedir "github.com/mitchellh/go-homedir"
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
	home, err := homedir.Dir()
	if err != nil {
		logging.Logger().Error("error getting home path", zap.Error(err))
		return ""
	}

	return filepath.Join(home, ".config", "gcloud", "application_default_credentials.json")
}

// TODO: Centralize these requests somewhere since they are shared between the server and client

// CredentialHandlerRequest wraps in a credential
type CredentialHandlerRequest struct {
	CredentialFile []byte `json:"credential_file"`
}

// CredentialHandlerResponse returns a credential request
type CredentialHandlerResponse struct {
	CredentialFilePath string `json:"credential_file_path"`
	CredentialFile     []byte `json:"credential_file"`
}

func login(cmd *cobra.Command, args []string) {
	credentialFile, err := ioutil.ReadFile(credential)
	if err != nil {
		logging.Logger().Fatal("credentials file does not exist - you must login using GCloud", zap.Error(err))
	}

	req := &CredentialHandlerRequest{CredentialFile: credentialFile}
	client := &http.Client{}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		logging.Logger().Fatal("error marshaling req", zap.Error(err))
	}

	b := bytes.NewBuffer(reqBytes)

	// TODO: Don't hardcode this
	resp, err := client.Post("http://localhost:3030/credentials", "application/json", b)
	if err != nil {
		logging.Logger().Fatal("error posting to credentials", zap.Error(err))
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		logging.Logger().Fatal("error reading credential response", zap.Error(err))
	}

	response := &CredentialHandlerResponse{}
	if err := json.Unmarshal(respBody, response); err != nil {
		logging.Logger().Fatal("error unmarshaling response", zap.Error(err))
	}

	// TODO: Formalize this, obviously
	spew.Dump(response)
}
