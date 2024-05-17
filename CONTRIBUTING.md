# Contributing

Thanks for your interest in contributing to the 1Password vault-plugin-secrets-onepassword project! 🙌 We appreciate your time and effort. Here are some guidelines to help you get started.

## Building
Run the following command to build the 1Password vault-plugin-secrets-onepassword:

```sh
make build 
```
This will build into the vault/plugins/op-connect binary.

## Testing
To demonstrate and evaluate the 1Password secrets engine, you can start HashiCorp Vault server in [development mode](https://developer.hashicorp.com/vault/docs/concepts/dev-server). The Vault server will be automatically initialized and unsealed in this configuration, and you do not need to register the plugin. 

```sh
vault server -dev -dev-root-token-id=root -dev-plugin-dir=./vault/plugins -log-level=debug
```
> **Warning:** Running Vault in development mode is useful for evaluating the plugin, but should **never** be used in production.

Connect to the Vault server in a **new** terminal to enable the secrets engine and start using it. For more information on how to enable and configure the plugin, follow the steps in the [QUICKSTART.md](./QUICKSTART.md)

To run the Go tests run the following command:

```sh
make test
```

To include the test coverage run the following command:

```sh
make test/coverage
```

## Debugging

#### Build the binary without optimizations

Run the following command to build the binary without enabling optimizations:

```sh
go build -o vault/plugins/op-connect -gcflags="all=-N -l" .
```

#### Start a debugging session
Run the following command to start a Delve debugging session:

```sh
dlv debug . -- --debug
```

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