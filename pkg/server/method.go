package server

// HTTPMethod wraps all HTTP methods valid for a route
type HTTPMethod string

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods
const (
	GET     HTTPMethod = "GET"
	HEAD    HTTPMethod = "HEAD"
	POST    HTTPMethod = "POST"
	PUT     HTTPMethod = "PUT"
	DELETE  HTTPMethod = "DELETE"
	CONNECT HTTPMethod = "CONNECT"
	OPTIONS HTTPMethod = "OPTIONS"
	TRACE   HTTPMethod = "TRACE"
	PATCH   HTTPMethod = "PATCH"
)

// String converts HTTP method to string
func (h HTTPMethod) String() string {
	return string(h)
}
