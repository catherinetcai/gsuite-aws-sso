package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
)

// AWS ...
type AWS struct {
	IAM iamiface.IAMAPI
	STS stsiface.STSAPI
}

// New ...
func New(sess *session.Session) *AWS {
	iamSvc := iam.New(sess)
	stsSvc := sts.New(sess)

	return &AWS{
		IAM: iamSvc,
		STS: stsSvc,
	}
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
func (a *AWS) AssumeRole(role string) (*Credentials, error) {
	roleArn, err := a.GetRoleARN(role)
	if err != nil {
		return nil, err
	}

	// TODO: Need to be able to override duration
	input := sts.AssumeRoleInput{
		RoleArn: aws.String(roleArn),
	}
	out, err := a.STS.AssumeRole(&input)
	if err != nil {
		return nil, err
	}

	return &Credentials{
		AccessKeyId:     *out.Credentials.AccessKeyId,
		SecretAccessKey: *out.Credentials.SecretAccessKey,
	}, nil
}
