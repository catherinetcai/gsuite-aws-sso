package cmd

import (
	"github.com/catherinetcai/gsuite-aws-sso/pkg/config"
	"github.com/catherinetcai/gsuite-aws-sso/pkg/http/middleware"
	"github.com/catherinetcai/gsuite-aws-sso/pkg/logging"
	"github.com/catherinetcai/gsuite-aws-sso/pkg/oauth"
	"github.com/catherinetcai/gsuite-aws-sso/pkg/server"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var logger = logging.Logger()

var serverCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the server",
	Run:   runServer,
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func runServer(cmd *cobra.Command, args []string) {
	oauthClient := oauth.NewClient(
		oauth.WithLogger(logger),
		oauth.WithConfig(config.Get().OAuth),
	)

	// TODO: Need to handle any unmatched routes
	router := mux.NewRouter()
	router.Use(middleware.NewLogging(logging.Logger()).Middleware)
	authRouter := router.PathPrefix("/auth/").Subrouter()
	authRouter.HandleFunc("/login", oauthClient.LoginHandler)
	authRouter.HandleFunc("/callback", oauthClient.CallbackHandler).Queries("code", "{code}")

	s, err := server.New(
		server.WithLogger(logger),
		server.WithPort(config.Get().Server.Port),
		server.WithRouter(router),
	)
	if err != nil {
		logging.Logger().Fatal("failed to start server", zap.Error(err))
	}

	logging.Logger().Fatal("error running", zap.Error(s.Run()))
}
