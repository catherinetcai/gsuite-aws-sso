package aws

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/catherinetcai/gsuite-aws-sso/pkg/file"
	"github.com/catherinetcai/gsuite-aws-sso/pkg/logging"
	"github.com/catherinetcai/gsuite-aws-sso/pkg/role"
	"go.uber.org/zap"
	ini "gopkg.in/ini.v1"
)

const (
	defaultProfileSectionName = "default"
)

// AWS ...
type AWS struct {
	IAM  iamiface.IAMAPI
	STS  stsiface.STSAPI
	sess *session.Session
}

// New ...
func New(sess *session.Session) *AWS {
	iamSvc := iam.New(sess)
	stsSvc := sts.New(sess)

	return &AWS{
		IAM:  iamSvc,
		STS:  stsSvc,
		sess: sess,
	}
}

// CredentialLocation returns the location the AWS credentials will be stored in.
// Defaults to ~/.aws/credentials
func (a *AWS) CredentialLocation() (string, error) {
	userHome, err := file.WithUserHomeDir(".aws", "credentials")
	if err != nil {
		return "", err
	}

	return userHome, nil
}

// GetRegion gets the region associated with the calling credentials
func (a *AWS) GetRegion() string {
	return *a.sess.Config.Region
}

// GetCredential takes in a role ID and returns a set of wrapped credentials, or an error
func (a *AWS) GetCredential(roleID string) (*role.Credential, error) {
	roleCreds, err := a.AssumeRole(roleID)
	if err != nil {
		return nil, err
	}

	// TODO: Extract this from this function
	// Create an empty credential file
	credFile := ini.Empty()

	// Seed the role creds
	defaultSection, err := credFile.NewSection(defaultProfileSectionName)
	if err != nil {
		logging.Logger().Error("error creating section", zap.Error(err))
		return nil, err
	}

	defaultSection.NewKey("aws_access_key_id", *roleCreds.AccessKeyId)
	defaultSection.NewKey("aws_secret_access_key", *roleCreds.SecretAccessKey)
	defaultSection.NewKey("aws_session_token", *roleCreds.SessionToken)
	defaultSection.NewKey("region", a.GetRegion())

	var b bytes.Buffer

	if _, err := credFile.WriteTo(&b); err != nil {
		return nil, err
	}

	credLocation, err := a.CredentialLocation()
	if err != nil {
		return nil, err
	}

	return &role.Credential{
		Raw:      b.Bytes(),
		Location: credLocation,
	}, nil
}

// GetRoleARN gets the role ARN from a role name
func (a *AWS) GetRoleARN(role string) (string, error) {
	input := iam.GetRoleInput{
		RoleName: aws.String(role),
	}

	out, err := a.IAM.GetRole(&input)
	if err != nil {
		return "", err
	}

	return *out.Role.Arn, nil
}

// AssumeRole will assume a role from the name
func (a *AWS) AssumeRole(role string) (*sts.Credentials, error) {
	// TODO: Need to be able to override duration
	input := sts.AssumeRoleInput{
		RoleArn: aws.String(role),
		// TODO: Pass in email
		RoleSessionName: aws.String("default"),
	}
	out, err := a.STS.AssumeRole(&input)
	if err != nil {
		logging.Logger().Error("error assuming role", zap.Error(err))
		return nil, err
	}

	return out.Credentials, nil
}
