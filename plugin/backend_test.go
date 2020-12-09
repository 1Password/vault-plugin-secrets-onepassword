package onepassword

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/1Password/connect-sdk-go/onepassword"
	"github.com/google/uuid"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	defaultLeaseTTLHr = 1
	maxLeaseTTLHr     = 12
)

var OpConnectClient = &TestClient{}
var Items = map[string]*onepassword.Item{}

func TestClientWithConfigSet(t *testing.T) {
	t.Parallel()

	b, reqStorage := getTestBackendNoCachedClient(t)

	config := map[string]interface{}{
		"op_connect_token": opToken,
		"op_connect_host":  host,
		"op_vault":         opVault,
	}
	testConfigUpdate(t, b, reqStorage, config)
	client, err := b.(*backend).OnePasswordConnectClient(reqStorage)
	if err != nil {
		t.Errorf("Retrieving 1Password Connect Client caused an unexpected error.")
	}
	if client == nil {
		t.Errorf("Client is unexpectantly nil.")
	}
}

func TestClientWithNoConfigSet(t *testing.T) {
	t.Parallel()

	b, reqStorage := getTestBackendNoCachedClient(t)
	client, err := b.(*backend).OnePasswordConnectClient(reqStorage)
	if err == nil {
		t.Errorf("No config set should return error when fecthing client.")
	}
	if client != nil {
		t.Errorf("Returned client should be nil on error.")
	}
}

func getTestBackendWithCachedClient(tb testing.TB) (logical.Backend, logical.Storage) {
	tb.Helper()
	config := getTestConfig()
	b, err := Factory(context.Background(), config)

	setOnePassswordConnectMocks()
	b.(*backend).configCache.Set("client", OpConnectClient, -1)

	if err != nil {
		tb.Fatal(err)
	}
	return b.(*backend), config.StorageView
}

func getTestBackendNoCachedClient(tb testing.TB) (logical.Backend, logical.Storage) {
	tb.Helper()
	config := getTestConfig()
	b, err := Factory(context.Background(), config)

	if err != nil {
		tb.Fatal(err)
	}
	return b.(*backend), config.StorageView
}

func setOnePassswordConnectMocks() {
	DoGetVaultsFunc = listVaults
	DoGetItemsFunc = listItems
	DoCreateItemFunc = createItem
	DoGetItemFunc = getItem
	DoDeleteItemFunc = deleteItem
	DoUpdateItemFunc = updateItem
}

func getTestConfig() *logical.BackendConfig {
	config := logical.TestBackendConfig()
	config.StorageView = new(logical.InmemStorage)
	config.Logger = hclog.NewNullLogger()
	config.System = &logical.StaticSystemView{
		DefaultLeaseTTLVal: defaultLeaseTTLHr * time.Hour,
		MaxLeaseTTLVal:     maxLeaseTTLHr * time.Hour,
	}
	return config
}
func listVaults() ([]onepassword.Vault, error) {
	vaults := []onepassword.Vault{
		{
			Description: "Test Vault1",
			ID:          "test1",
		},
		{
			Description: "Test Vault2",
			ID:          "test2",
		},
	}
	return vaults, nil
}

func createItem(item *onepassword.Item, vaultUUID string) (*onepassword.Item, error) {
	item.ID = uuid.New().String()
	Items[item.ID] = item
	return item, nil
}

func updateItem(item *onepassword.Item, vaultUUID string) (*onepassword.Item, error) {
	Items[item.ID] = item
	return item, nil
}

func getItem(itemUUID, vaultUUID string) (*onepassword.Item, error) {
	item, found := Items[itemUUID]
	if !found {
		return nil, fmt.Errorf("Could not retrieve item with id %v", itemUUID)
	}
	return item, nil
}

func deleteItem(item *onepassword.Item, vaultUUID string) error {
	delete(Items, item.ID)
	return nil
}

func listItems(vaultUUID string) ([]onepassword.Item, error) {
	items := []onepassword.Item{}
	for _, item := range Items {
		items = append(items, *item)
	}
	return items, nil
}
