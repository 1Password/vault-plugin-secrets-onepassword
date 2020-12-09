package onepassword

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type Config struct {
	OPToken string
	OPVault string
	Host    string
}

func pathConfig(b *backend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "config",
			Fields: map[string]*framework.FieldSchema{
				"op_connect_token": {
					Type:        framework.TypeString,
					Description: `The 1Password Connect API Token`,
				},
				"op_vault": {
					Type:        framework.TypeString,
					Description: "The id of the default 1Password vault to access",
				},
				"op_connect_host": {
					Type:        framework.TypeString,
					Description: "The host address for the 1Password Connect API",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathConfigRead,
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.pathConfigWrite,
				},
			},

			HelpSynopsis:    pathConfigHelpSyn,
			HelpDescription: pathConfigHelpDesc,
		},
	}
}

func (b *backend) pathConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	cfg, err := getConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"op_connect_host":  cfg.Host,
			"op_vault":         cfg.OPVault,
			"op_connect_token": cfg.OPToken,
		},
	}, nil
}

func (b *backend) pathConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := getConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		config = &Config{}
	}

	opToken, opTokenFound := data.GetOk("op_connect_token")
	if opTokenFound {
		config.OPToken = opToken.(string)
	}

	host, hostFound := data.GetOk("op_connect_host")
	if hostFound {
		config.Host = host.(string)
	}

	if hostFound || opTokenFound {
		b.configCache.Delete(clientCache)
	}

	opVault, ok := data.GetOk("op_vault")
	if ok {
		config.OPVault = opVault.(string)
		b.configCache.Delete(defaultVaultCache)
	}

	entry, err := logical.StorageEntryJSON("config", config)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

func getConfig(ctx context.Context, s logical.Storage) (*Config, error) {
	rawConfig, err := s.Get(ctx, "config")
	if err != nil {
		return nil, err
	}
	if rawConfig == nil {
		return nil, nil
	}

	var config Config
	if err := rawConfig.DecodeJSON(&config); err != nil {
		return nil, err
	}

	return &config, err
}

const pathConfigHelpSyn = `
Configure the 1Password Connect backend.
`

const pathConfigHelpDesc = `
The 1Password Connect backend requires an API token for accessing 1Password Connect.
This endpoint is used to set the API token as well as the host address in which to
access 1Password Connect and the default 1Password vault to interact with.
`
