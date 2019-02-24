package aws

// Credentials is representative of AWS credentials
type Credentials struct {
	Account         string `json:"account"`
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
}
