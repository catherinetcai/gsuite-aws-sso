package server

import (
	"net/http"

	httphelper "github.com/catherinetcai/gsuite-aws-sso/pkg/http"
)

// HealthHandler is a health check endpoint that just returns a 200 to show the
// server is now alive and handling traffic.
func (s *Server) HealthHandler(w http.ResponseWriter, req *http.Request) {
	httphelper.JSONResponse(w, struct{}{}, http.StatusOK)
}
