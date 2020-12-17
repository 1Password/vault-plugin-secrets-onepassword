package onepassword

import (
	"context"
	"fmt"
	"strings"

	"github.com/1Password/connect-sdk-go/onepassword"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathVaults(b *backend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "vaults/?",

			HelpSynopsis: `
			Allows for listing all vaults available with configured 1Password API Token
			`,

			HelpDescription: strings.TrimSpace(`
			Allows for listing all vaults available with configured 1Password API Token
			`),

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: b.handleListVaults,
					Summary:  "List all available vaults",
				},
			},

			ExistenceCheck: b.handleExistenceCheck,
		},
	}
}

func (b *backend) handleListVaults(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	client, err := b.OnePasswordConnectClient(req.Storage)
	if err != nil {
		return nil, fmt.Errorf("Unable to fetch client")
	}

	vaults, err := client.GetVaults()
	if err != nil {
		return nil, errwrap.Wrapf("Unable to list vaults: {{err}}", err)
	}
	return generateVaultList(vaults), nil
}

func generateVaultList(vaults []onepassword.Vault) *logical.Response {
	var vault_list []string
	vault_map := make(map[string]interface{})
	for i := 0; i < len(vaults); i++ {
		key := fmt.Sprintf("%v %v", vaults[i].Name, vaults[i].ID)
		vault_list = append(vault_list, key)
		vault_map[key] = vaults[i].ID
	}

	return logical.ListResponseWithInfo(vault_list, vault_map)
}
