package server

import "net/http"

// Route wraps a route
type Route struct {
	Path        string
	HandlerFunc func(http.ResponseWriter, *http.Request)
	Method      HTTPMethod
	// Accepts pairs of query values as a slice of strings
	Queries []string
}
