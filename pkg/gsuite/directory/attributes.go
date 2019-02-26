package directory

// Attributes wraps the custom attributes for a directory user
type Attributes struct {
	IAMRole         []IAMRole
	SessionDuration string
}

// IAMRole wraps the IAM role info set in the custom attributes of a user
type IAMRole struct {
	// Type is just "work"
	Type string
	// Value is a comma delimited list of the IAM role and the provider
	Value string
}
