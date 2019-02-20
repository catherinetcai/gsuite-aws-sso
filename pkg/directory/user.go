package directory

// User encapsulates a user's identifying email and their custom attributes
type User struct {
	Email            string
	CustomAttributes map[string]string
}
