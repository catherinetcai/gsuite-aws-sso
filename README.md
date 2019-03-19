# GSuite AWS SSO
[![CircleCI](https://circleci.com/gh/catherinetcai/gsuite-aws-sso.svg?style=svg)](https://circleci.com/gh/catherinetcai/gsuite-aws-sso)

## Introduction
Setting up SSO for Google Apps via the AWS console is pretty straightforward. Google has a very [helpful guide](https://support.google.com/a/answer/6194963) for getting set up.

Unfortunately, getting SSO for the AWS CLI isn't quite as easy. There are quite a few solutions that already exist out there, such as [aws-google-auth](https://github.com/cevoaustralia/aws-google-auth) and [saml2aws](https://github.com/Versent/saml2aws). Though they get the job done, their underlying implementation is quite brittle as they use page scraping.

This tool is a less brittle way to achieve this effect, though this method is still a bit of a hack.

## Building
The [Makefile](Makefile) has invocations for building the client and server binaries.

```bash
make # Will build client and server binaries and place them in release/ folder
```

## Running
### Client
The client can be run via the following:

```bash
make client-login
# OR
./client login
```

### Server
The server can be run via the following:

```bash
make run
# OR
./server run
```

## Design
There are a couple of core requirements to make the AWS SSO CLI work:
* User must login successfully into Google and prove ownership of their account within the GSuite organization
* The user must have the appropriate SSO attributes in order to successfully map them to a role
* User receives temporary set of credentials via AWS STS
* These credentials must be seeded into the `~/.aws/credentials` file in the proper format so that the user does not have to use environment variables

The problem with Google App is that there is no programmable API to enable logging in the user and mapping their role to the AWS role.

Essentially, the way to solve this problem would be an AuthN/AuthZ coordinator service that is able to 1) use Google's Auth 2.0 to validate a user owns their email 2) look up the user and the user's attributes within a GSuite organization's directory (enabled via the [Admin SDK](https://developers.google.com/admin-sdk/) 3) map a user to the appropriate AWS role via their custom attributes 4) get STS credentials for that AWS role.

This service will end up being two components. One side is the service component performing all the coordination. The other side is a CLI tool that will take the response from the server (the STS credentials) and then seed them appropriately into `~/.aws/credentials`.

## Flow
### Client
The client relies on default GCloud Auth credentials in order to work.

```bash
# Must point CLOUDSDK_PYTHON to valid Python 3
export CLOUDSDK_PYTHON=$HOME/.pyenv/shims/python3

# Get Google Cloud login credentials
gcloud auth application-default login
```

Then, log into the Google account with the Wurl account.

GCloud credentials will get seeded to `/Users/<user>/.config/gcloud/application_default_credentials.json`.

### Server
#### GSuite
The server must be provisioned with Client ID, service account (with domain-wide delegation) credentials, scopes, and a GSuite Admin user to impersonate.

The server will use the credential file in order to log into a service account. The Directory API requires an Admin user, so it will impersonate the Admin user in order to do work (in this case, be able to get user directory info).

#### AWS
The server must also be provisioned with AWS credentials that are able to assume the roles that are available via GSuite.
