package onepassword

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/1Password/connect-sdk-go/onepassword"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	defaultLeaseTTLHr = 1
	maxLeaseTTLHr     = 12
)

var (
	vaultId1   = "2zhyki73hfeodll4ljtgzwftwu"
	vaultName1 = "Test Vault1"
)

var (
	vaultId2   = "2zhyki73hfeodll4ljtgzwftwv"
	vaultName2 = "Test Vault2"
)

var OpConnectClient = &TestClient{}
var Items = map[string]*onepassword.Item{}
var characters = []rune("abcdefghijklmnopqrstuvwxyz123456789")

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
	DoGetItemsByTitleFunc = getItemsByTitle
	DoDeleteItemFunc = deleteItem
	DoDeleteItemByIDFunc = deleteItemByID
	DoDeleteItemByTitleFunc = deleteItemByTitle
	DoUpdateItemFunc = updateItem
	DoDownloadFileFunc = downloadFile
	DoGetFileFunc = getFile
	DoGetFileContentFunc = getFileContent
	DoGetFilesFunc = getFiles
	DoGetItemByUUIDFunc = getItemByUUID
	DoGetVaultFunc = getVault
	DoGetVaultByTitleFunc = getVaultByTitle
	DoGetVaultByUUIDFunc = getVaultByUUID
	DoLoadStructFunc = loadStruct
	DoLoadStructFromItemFunc = loadStructFromItem
	DoLoadStructFromItemByUUIDFunc = loadStructFromItemByUUID
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
			Name: vaultName1,
			ID:   vaultId1,
		},
		{
			Name: vaultName2,
			ID:   vaultId2,
		},
	}
	return vaults, nil
}

func createItem(item *onepassword.Item, vaultUUID string) (*onepassword.Item, error) {
	item.ID = RandID()
	item.CreatedAt = time.Now()
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

func getItemsByTitle(title, vaultUUID string) ([]onepassword.Item, error) {
	var itemsWithTitle []onepassword.Item
	for _, item := range Items {
		if item.Title == title {
			itemsWithTitle = append(itemsWithTitle, *item)
		}
	}
	return itemsWithTitle, nil
}

func deleteItem(item *onepassword.Item, vaultUUID string) error {
	delete(Items, item.ID)
	return nil
}

func deleteItemByID(itemUUID, vaultUUID string) error {
	return fmt.Errorf("This method is currently not supported by the test client")
}

func deleteItemByTitle(itemName, vaultUUID string) error {
	return fmt.Errorf("This method is currently not supported by the test client")
}

func downloadFile(file *onepassword.File, targetDirectory string, overwrite bool) (string, error) {
	return "", fmt.Errorf("This method is currently not supported by the test client")
}

func getFile(uuid string, itemQuery string, vaultQuery string) (*onepassword.File, error) {
	return nil, fmt.Errorf("This method is currently not supported by the test client")
}

func getFileContent(file *onepassword.File) ([]byte, error) {
	return nil, fmt.Errorf("This method is currently not supported by the test client")
}

func getFiles(itemQuery string, vaultQuery string) ([]onepassword.File, error) {
	return nil, fmt.Errorf("This method is currently not supported by the test client")
}

func getItemByUUID(uuid string, vaultQuery string) (*onepassword.Item, error) {
	return nil, fmt.Errorf("This method is currently not supported by the test client")
}

func getVault(uuid string) (*onepassword.Vault, error) {
	return nil, fmt.Errorf("This method is currently not supported by the test client")
}

func getVaultByTitle(title string) (*onepassword.Vault, error) {
	return nil, fmt.Errorf("This method is currently not supported by the test client")
}

func getVaultByUUID(uuid string) (*onepassword.Vault, error) {
	return nil, fmt.Errorf("This method is currently not supported by the test client")
}

func loadStruct(config interface{}) error {
	return fmt.Errorf("This method is currently not supported by the test client")
}

func loadStructFromItem(config interface{}, itemQuery string, vaultQuery string) error {
	return fmt.Errorf("This method is currently not supported by the test client")
}

func loadStructFromItemByUUID(config interface{}, itemUUID string, vaultQuery string) error {
	return fmt.Errorf("This method is currently not supported by the test client")
}

func listItems(vaultUUID string) ([]onepassword.Item, error) {
	items := []onepassword.Item{}
	for _, item := range Items {
		items = append(items, *item)
	}
	return items, nil
}

func RandID() string {
	b := make([]rune, 26)
	for i := range b {
		b[i] = characters[rand.Intn(len(characters))]
	}
	return string(b)
}
