package role

// Service is an interface that implements getting credentials and
// seeds them in the right location
type Service interface {
	// GetCredential takes a role identifier and returns a wrapped credential object,
	// which is the credential in the file format, and the location of where to seed it
	GetCredential(roleID string) (*Credential, error)
}
