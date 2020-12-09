package onepassword

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/patrickmn/go-cache"
	"github.com/1Password/connect-sdk-go/connect"
)

const (
	envDefaultVaultVar = "OP_VAULT"
	clientCache        = "client"
	defaultVaultCache  = "vault"

	cacheCleanup    = 30 * time.Minute
	cacheExpiration = 30 * time.Minute
)

// Factory returns a new backend as logical.Backend.
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend() *backend {
	var b = &backend{
		store:       make(map[string][]byte),
		configCache: cache.New(cacheExpiration, cacheCleanup),
	}

	b.Backend = &framework.Backend{
		BackendType: logical.TypeLogical,
		Help:        strings.TrimSpace(opHelp),
		PathsSpecial: &logical.Paths{
			SealWrapStorage: []string{
				"config",
			},
		},

		Paths: framework.PathAppend(
			pathConfig(b),
			pathItems(b),
			pathVaults(b),
		),
		Secrets: []*framework.Secret{},
	}

	return b
}

type backend struct {
	*framework.Backend

	store map[string][]byte

	configCache *cache.Cache
}

func (b *backend) OnePasswordConnectClient(s logical.Storage) (connect.Client, error) {

	cachedClient, found := b.configCache.Get(clientCache)
	if found {
		b.Logger().Debug("Using cached client")
		return cachedClient.(connect.Client), nil
	}

	b.Logger().Debug("Creating new client")
	config, err := getConfig(context.Background(), s)
	if err != nil {
		return nil, errwrap.Wrapf("Error retrieving config for client: {{err}}", err)
	}

	if config == nil {
		return nil, fmt.Errorf("No config set for op backend.")
	}

	client := connect.NewClient(config.Host, config.OPToken)
	b.configCache.Set(clientCache, client, -1)

	return client, nil
}

func (b *backend) getVault(s logical.Storage, vaultParam string) (string, error) {
	defaultVault, envVaultFound := os.LookupEnv(envDefaultVaultVar)
	if vaultParam != "" {
		return vaultParam, nil
	}

	if envVaultFound {
		return defaultVault, nil
	}

	cachedDefaultVault, found := b.configCache.Get(defaultVaultCache)
	if found {
		b.Logger().Debug("Using cached default vault")
		return cachedDefaultVault.(string), nil
	}

	config, err := getConfig(context.Background(), s)
	if err == nil && config != nil && config.OPVault != "" {
		b.configCache.Set(defaultVaultCache, config.OPVault, -1)
		return config.OPVault, nil
	}

	return "", fmt.Errorf("No vault has been specified")
}

func (b *backend) handleExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	out, err := req.Storage.Get(ctx, req.Path)
	if err != nil {
		return false, errwrap.Wrapf("existence check failed: {{err}}", err)
	}

	return out != nil, nil
}

const opHelp = `
The OP backend is a secrets backend that allows for the retreival of items from 1Password using 1Password Connect.
`
