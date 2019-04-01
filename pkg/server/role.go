package server

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	httphelper "github.com/catherinetcai/gsuite-aws-sso/pkg/http"
	"github.com/catherinetcai/gsuite-aws-sso/pkg/shared/handlers"
	"go.uber.org/zap"
)

// https://stackoverflow.com/questions/13317987/authorizing-command-line-tool-to-consume-google-apis-through-oauth2-0-or-anythi

// CredentialHandler takes in a client access token, validates it, and then returns a set of credentials
func (s *Server) CredentialHandler(w http.ResponseWriter, req *http.Request) {
	request := &handlers.CredentialHandlerRequest{}
	response := &handlers.CredentialHandlerResponse{}

	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		s.logger.Error("error reading credential request", zap.Error(err))
		httphelper.JSONResponse(w, struct{}{}, http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, request); err != nil {
		s.logger.Error("error unmarshalling response", zap.Error(err))
		httphelper.JSONResponse(w, struct{}{}, http.StatusInternalServerError)
		return
	}

	tokenSource, err := s.oAuthSvc.TokenSourceFromCredentials(context.Background(), request.CredentialFile)
	if err != nil {
		s.logger.Error("error getting token source from credentials", zap.Error(err))
		httphelper.JSONResponse(w, struct{}{}, http.StatusUnauthorized)
		return
	}

	idToken, err := s.oAuthSvc.IDToken(tokenSource)
	if err != nil {
		s.logger.Error("error getting id token", zap.Error(err))
		httphelper.JSONResponse(w, struct{}{}, http.StatusUnauthorized)
		return
	}

	user, err := s.directorySvc.GetUser(idToken.Email)
	if err != nil {
		httphelper.JSONResponse(w, struct{}{}, http.StatusBadRequest)
		return
	}

	s.logger.Info("Got user", zap.Any("user", user))

	// Token is valid, therefore go and try to get the role
	cred, err := s.roleSvc.GetCredential(user.CredentialID)
	if err != nil {
		// TODO: We need to do better AWS error handling
		s.logger.Error("error getting credential for email",
			zap.String("email", idToken.Email),
			zap.Error(err))
		httphelper.JSONResponse(w, struct{}{}, http.StatusBadRequest)
		return
	}

	response.CredentialFile = cred.Raw
	response.CredentialFilePath = cred.Location

	// https://stackoverflow.com/questions/24442668/google-oauth-api-to-get-users-email-address
	// TODO: Move this into the OAuth client
	// https://developers.google.com/identity/protocols/OpenIDConnect#discovery
	// oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	// resp, err := oauthClient.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	// if err != nil {
	// 	s.logger.Error("error getting user info", zap.Error(err))
	// 	httphelper.JSONResponse(w, response, http.StatusUnauthorized)
	// 	return
	// }

	httphelper.JSONResponse(w, response, http.StatusOK)
}
