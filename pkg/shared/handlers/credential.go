package handlers

// CredentialHandlerRequest wraps in a credential
type CredentialHandlerRequest struct {
	CredentialFile []byte `json:"credential_file"`
}

// CredentialHandlerResponse returns a credential response
type CredentialHandlerResponse struct {
	CredentialFilePath string `json:"credential_file_path"`
	CredentialFile     []byte `json:"credential_file"`
}
