package onepassword

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/1Password/connect-sdk-go/onepassword"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathItems(b *backend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "vaults/" + framework.GenericNameRegex("vault") + "/items/" + framework.GenericNameRegex("id"),

			HelpSynopsis: `
			Allows for reading and deleting items from a 1Password vault.
			`,

			HelpDescription: strings.TrimSpace(`
			Allows for reading and deleting items from a 1Password vault.
			`),

			Fields: map[string]*framework.FieldSchema{
				"id": {
					Type:        framework.TypeString,
					Description: "Specifies the id of the item.",
					Required:    true,
				},
				"vault": {
					Type:        framework.TypeString,
					Description: "Specifies the id of the vault.",
					Required:    true,
				},
				"url": {
					Type:        framework.TypeString,
					Description: "Specifies the url of the item",
				},
				"category": {
					Type:        framework.TypeString,
					Description: "Specifies the category of the item to generate; database, login, and password are currently supported",
					Required:    true,
				},
				"title": {
					Type:        framework.TypeString,
					Description: "Specifies the title of the item",
					Required:    true,
				},
				"fields": {
					Type:        framework.TypeSlice,
					Description: strings.TrimSpace(fieldsDescription),
				},
				"sections": {
					Type:        framework.TypeSlice,
					Description: strings.TrimSpace(sectionsDescription),
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.handleReadItem,
					Summary:  "Retrieve the item from the specified location",
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.handleDeleteItem,
					Summary:  "Delete an item from the specified location",
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.handleUpdateItem,
					Summary:  "Update an item at the specified location",
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.handleUpdateItem,
					Summary:  "Update an item at the specified location",
				},
			},

			ExistenceCheck: b.handleExistenceCheck,
		},
		{
			Pattern: "vaults/" + framework.GenericNameRegex("vault") + "/items/?",

			HelpSynopsis: `
			Allows for creating an item or listing items for a 1Password Vault.
			`,

			HelpDescription: strings.TrimSpace(`
			Allows for creating an item or listing items for a 1Password Vault.
			`),

			Fields: map[string]*framework.FieldSchema{
				"vault": {
					Type:        framework.TypeString,
					Description: "Specifies the id of the vault.",
					Required:    true,
				},
				"url": {
					Type:        framework.TypeString,
					Description: "Specifies the url of the item",
				},
				"category": {
					Type:        framework.TypeString,
					Description: "Specifies the category of the item to generate; database, login, and password are currently supported",
					Required:    true,
				},
				"title": {
					Type:        framework.TypeString,
					Description: "Specifies the title of the item",
					Required:    true,
				},
				"fields": {
					Type:        framework.TypeSlice,
					Description: strings.TrimSpace(fieldsDescription),
				},
				"sections": {
					Type:        framework.TypeSlice,
					Description: strings.TrimSpace(sectionsDescription),
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.handleWriteItem,
					Summary:  "Store a 1Password item at the specified vault.",
				},
				logical.ListOperation: &framework.PathOperation{
					Callback: b.handleListItems,
					Summary:  "List all items in the specified vault",
				},
			},

			ExistenceCheck: b.handleExistenceCheck,
		},
	}
}

func (b *backend) handleReadItem(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	client, err := b.OnePasswordConnectClient(req.Storage)
	if err != nil {
		return nil, errwrap.Wrapf("Unable to fetch client: {{err}}", err)
	}

	id := data.Get("id").(string)
	vaultParam := data.Get("vault").(string)
	vault, err := b.getVault(req.Storage, vaultParam)
	if err != nil {
		return nil, errwrap.Wrapf("Unable to retrieve item: {{err}}", err)
	}

	item, err := client.GetItem(id, vault)
	if err != nil {
		return nil, errwrap.Wrapf("Unable to retrieve item: {{err}}", err)
	}
	fields := item.Fields

	field_map := make(map[string]interface{})
	for i := 0; i < len(fields); i++ {
		field_map[fields[i].Label] = fields[i].Value
	}

	return &logical.Response{
		Data: field_map,
	}, nil
}

func (b *backend) handleWriteItem(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	client, err := b.OnePasswordConnectClient(req.Storage)
	if err != nil {
		return nil, fmt.Errorf("Unable to fetch client")
	}

	vaultParam := data.Get("vault").(string)
	vault, err := b.getVault(req.Storage, vaultParam)
	if err != nil {
		return nil, errwrap.Wrapf("Unable to create item: {{err}}", err)
	}

	item, err := buildItem(b, vault, data)
	if err != nil {
		return nil, errwrap.Wrapf("Unable to create item: {{err}}", err)
	}

	createdItem, err := client.CreateItem(item, vault)
	if err != nil {
		return nil, errwrap.Wrapf("Unable to create item: {{err}}", err)
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"ID":        createdItem.ID,
			"category":  createdItem.Category,
			"createdAt": createdItem.CreatedAt,
		},
	}, nil
}

func (b *backend) handleListItems(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	client, err := b.OnePasswordConnectClient(req.Storage)
	if err != nil {
		return nil, fmt.Errorf("Unable to fetch client")
	}

	vaultParam := data.Get("vault").(string)
	vault, err := b.getVault(req.Storage, vaultParam)
	if err != nil {
		return nil, errwrap.Wrapf("Unable to list items: {{err}}", err)
	}

	items, err := client.GetItems(vault)
	if err != nil {
		return nil, errwrap.Wrapf("Unable to list items: {{err}}", err)
	}

	var item_list []string
	item_map := make(map[string]interface{})
	for i := 0; i < len(items); i++ {
		item_list = append(item_list, items[i].ID)
		item_map[items[i].ID] = items[i].Title
	}

	return logical.ListResponseWithInfo(item_list, item_map), nil
}

func (b *backend) handleDeleteItem(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	client, err := b.OnePasswordConnectClient(req.Storage)
	if err != nil {
		return nil, fmt.Errorf("Unable to fetch client")
	}

	id := data.Get("id").(string)
	vaultParam := data.Get("vault").(string)
	vault, err := b.getVault(req.Storage, vaultParam)
	if err != nil {
		return nil, errwrap.Wrapf("Unable to delete item: {{err}}", err)
	}

	item := onepassword.Item{
		ID: id,
		Vault: onepassword.ItemVault{
			ID: vault,
		},
	}

	err = client.DeleteItem(&item, vault)
	if err != nil {
		return nil, errwrap.Wrapf("Unable to delete item: {{err}}", err)
	}

	return nil, nil
}

func (b *backend) handleUpdateItem(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	client, err := b.OnePasswordConnectClient(req.Storage)
	if err != nil {
		return nil, fmt.Errorf("Unable to fetch client")
	}

	vaultParam := data.Get("vault").(string)
	vault, err := b.getVault(req.Storage, vaultParam)
	if err != nil {
		return nil, errwrap.Wrapf("Unable to update item: {{err}}", err)
	}

	item, err := buildItem(b, vault, data)
	if err != nil {
		return nil, errwrap.Wrapf("Unable to update item: {{err}}", err)
	}

	item.ID = data.Get("id").(string)

	createdItem, err := client.UpdateItem(item, vault)
	if err != nil {
		return nil, errwrap.Wrapf("Unable to update item: {{err}}", err)
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"ID":        createdItem.ID,
			"category":  createdItem.Category,
			"createdAt": createdItem.CreatedAt,
		},
	}, nil
}

func buildItem(b *backend, vault string, data *framework.FieldData) (*onepassword.Item, error) {
	item := onepassword.Item{
		Vault: onepassword.ItemVault{
			ID: vault,
		},
		Title: data.Get("title").(string),
		URLs: []onepassword.ItemURL{
			onepassword.ItemURL{
				Primary: true,
				URL:     data.Get("url").(string),
			},
		},
	}

	switch data.Get("category").(string) {
	case "login":
		item.Category = onepassword.Login
	case "password":
		item.Category = onepassword.Password
	case "database":
		item.Category = onepassword.Database
	}

	sections := data.Get("sections").([]interface{})
	for _, section := range sections {
		itemSection := onepassword.ItemSection{}
		marshalledSection, err := json.Marshal(section.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(marshalledSection, &itemSection)
		if err != nil {
			return nil, err
		}
		item.Sections = append(item.Sections, &itemSection)
	}

	fields := data.Get("fields").([]interface{})
	for _, field := range fields {
		itemField := onepassword.ItemField{}
		marshalledItem, err := json.Marshal(field.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(marshalledItem, &itemField)
		if err != nil {
			return nil, err
		}
		item.Fields = append(item.Fields, &itemField)
	}

	something, _ := json.Marshal(item)
	b.Logger().Debug(fmt.Sprintf("%v", string(something)))
	return &item, nil
}

const fieldsDescription = `The list of fields to create for the item. This is represented as a list of maps.
For more information on how to format fields please see the README`
const sectionsDescription = `The list of sections to create for the item. This is represented as a list of maps.
For more information on how to format fields please see the README`
