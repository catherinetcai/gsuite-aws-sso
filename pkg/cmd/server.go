package cmd

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/catherinetcai/gsuite-aws-sso/pkg/aws"
	"github.com/catherinetcai/gsuite-aws-sso/pkg/config"
	gdirectory "github.com/catherinetcai/gsuite-aws-sso/pkg/gsuite/directory"
	goauth "github.com/catherinetcai/gsuite-aws-sso/pkg/gsuite/oauth"
	"github.com/catherinetcai/gsuite-aws-sso/pkg/http/middleware"
	"github.com/catherinetcai/gsuite-aws-sso/pkg/logging"
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
	oauthClient := goauth.NewClient(
		goauth.WithLogger(logger),
		goauth.WithConfig(config.Get().OAuth),
	)

	directoryClient, err := gdirectory.NewClient(
		gdirectory.WithLogger(logger),
		gdirectory.WithImpersonationEmail(config.Get().GSuite.ImpersonationEmail),
		gdirectory.WithServiceAccountEmail(config.Get().GSuite.ServiceAccountEmail),
		// TODO: Make this flexible with both the base64 or a file path
		gdirectory.WithServiceAccountBase64EncodedFile(config.Get().GSuite.ServiceAccountBase64EncodedFile),
	)
	if err != nil {
		logging.Logger().Fatal("failed to initialize directory", zap.Error(err))
	}

	awsClient := aws.New(session.Must(session.NewSession()))
	if err != nil {
		logging.Logger().Fatal("failed to initialize aws", zap.Error(err))
	}

	// TODO: Need to handle any unmatched routes
	router := mux.NewRouter()
	router.Use(middleware.NewLogging(logging.Logger()).Middleware)

	s, err := server.New(
		server.WithLogger(logger),
		server.WithPort(config.Get().Server.Port),
		server.WithRouter(router),
		server.WithOAuth(oauthClient),
		server.WithDirectory(directoryClient),
		server.WithRole(awsClient),
	)
	if err != nil {
		logging.Logger().Fatal("failed to start server", zap.Error(err))
	}

	logging.Logger().Fatal("error running", zap.Error(s.Run()))
}
