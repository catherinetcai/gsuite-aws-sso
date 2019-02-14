# GSuite AWS SSO
## Introduction
Setting up SSO for Google Apps via the AWS console is pretty straightforward. Google has a very [helpful guide](https://support.google.com/a/answer/6194963) for getting set up.

Unfortunately, getting SSO for the AWS CLI isn't quite as easy. There are quite a few solutions that already exist out there, such as [aws-google-auth](https://github.com/cevoaustralia/aws-google-auth) and [saml2aws](https://github.com/Versent/saml2aws). Though they get the job done, their underlying implementation is quite brittle as they use page scraping.

This tool is a less brittle way to achieve this effect, though this method is still a bit of a hack.

## Design
There are a couple of core requirements to make the AWS SSO CLI work:
* User must login successfully into Google and prove ownership of their account within the GSuite organization
* The user must have the appropriate SSO attributes in order to successfully map them to a role
* User receives temporary set of credentials via AWS STS
* These credentials must be seeded into the `~/.aws/credentials` file in the proper format so that the user does not have to use environment variables

The problem with Google App is that there is no programmable API to enable logging in the user and mapping their role to the AWS role.

Essentially, the way to solve this problem would be an AuthN/AuthZ coordinator service that is able to 1) use Google's Auth 2.0 to validate a user owns their email 2) look up the user and the user's attributes within a GSuite organization's directory (enabled via the [Admin SDK](https://developers.google.com/admin-sdk/) 3) map a user to the appropriate AWS role via their custom attributes 4) get STS credentials for that AWS role.

This service will end up being two components. One side is the service component performing all the coordination. The other side is a CLI tool that will take the response from the server (the STS credentials) and then seed them appropriately into `~/.aws/credentials`.
