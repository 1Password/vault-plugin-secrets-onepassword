# Contributing

Thanks for your interest in contributing to the 1Password vault-plugin-secrets-onepassword project! 🙌 We appreciate your time and effort. Here are some guidelines to help you get started.

## Building
Run the following command to build the 1Password vault-plugin-secrets-onepassword:

```sh
go build .
```

This will create the `vault-plugin-secrets-onepassword` binary.

## Testing

To run the Go tests and check test coverage run the following command:

```sh
go test -v ./... -cover
```

## Debugging

### Start a debugging session
Run the following command to start a Delve debugging session:

```sh
dlv debug . -- --debug
```

Or use your IDE debugger. You can configure editors like GoLand to start a debugging session by passing the `--debug` flag as a program argument.

## Documentation Updates

If applicable, update the [README.md](./README.md) to reflect any changes introduced by the new code.

## Sign your commits

To get your PR merged, we require you to sign your commits.

### Sign commits with `1Password`

You can also sign commits using 1Password, which lets you sign commits with biometrics without the signing key leaving the local 1Password process.

Learn how to use [1Password to sign your commits](https://developer.1password.com/docs/ssh/git-commit-signing/).


### Sign commits with `ssh-agent`

Follow the steps below to set up commit signing with `ssh-agent`:

1. Generate an SSH key and add it to ssh-agent
2. Add the SSH key to your GitHub account
3. Configure git to use your SSH key for commit signing

### Sign commits `gpg`

Follow the steps below to set up commit signing with `gpg`

1. Generate a GPG key
2. Add the GPG key to your GitHub account
3. Configure git to use your GPG key for commit signing