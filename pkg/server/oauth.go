package server

import (
	"context"
	"net/http"

	httphelper "github.com/catherinetcai/gsuite-aws-sso/pkg/http"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// LoginHandler redirects a client to the Google OAuth login page
func (s *Server) LoginHandler(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, s.oAuthSvc.GetOAuthLoginURL(), http.StatusFound)
}

// CallbackHandler handles the OAuth callback
// https://developers.google.com/identity/protocols/OAuth2WebServer
func (s *Server) CallbackHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	idToken, err := s.oAuthSvc.Exchange(context.Background(), vars["code"])
	if err != nil {
		s.logger.Error("error exchanging OAuth code", zap.Error(err))
		httphelper.JSONResponse(w, struct{}{}, http.StatusBadRequest)
		return
	}

	// TODO: This is not how this should work
	httphelper.JSONResponse(w, idToken, http.StatusOK)
}
