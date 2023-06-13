package onepassword

import (
	"github.com/1Password/connect-sdk-go/onepassword"
)

type TestClient struct {
	GetVaultsFunc                 func() ([]onepassword.Vault, error)
	GetVaultsByTitleFunc          func(title string) ([]onepassword.Vault, error)
	GetItemFunc                   func(uuid string, vaultUUID string) (*onepassword.Item, error)
	GetItemsFunc                  func(vaultUUID string) ([]onepassword.Item, error)
	GetItemsByTitleFunc           func(title string, vaultUUID string) ([]onepassword.Item, error)
	GetItemByTitleFunc            func(title string, vaultUUID string) (*onepassword.Item, error)
	CreateItemFunc                func(item *onepassword.Item, vaultUUID string) (*onepassword.Item, error)
	UpdateItemFunc                func(item *onepassword.Item, vaultUUID string) (*onepassword.Item, error)
	DeleteItemFunc                func(item *onepassword.Item, vaultUUID string) error
	DeleteItemByIDFunc            func(itemUUID, vaultUUID string) error
	DeleteItemByTitleFunc         func(itemName, vaultUUID string) error
	DownloadFileFunc              func(file *onepassword.File, targetDirectory string, overwrite bool) (string, error)
	GetFileFunc                   func(uuid string, itemQuery string, vaultQuery string) (*onepassword.File, error)
	GetFileContentFunc            func(file *onepassword.File) ([]byte, error)
	GetFilesFunc                  func(itemQuery string, vaultQuery string) ([]onepassword.File, error)
	GetItemByUUIDFunc             func(uuid string, vaultQuery string) (*onepassword.Item, error)
	GetVaultFunc                  func(uuid string) (*onepassword.Vault, error)
	GetVaultByTitleFunc           func(title string) (*onepassword.Vault, error)
	GetVaultByUUIDFunc            func(uuid string) (*onepassword.Vault, error)
	LoadStructFunc                func(config interface{}) error
	LoadStructFromItemFunc        func(config interface{}, itemQuery string, vaultQuery string) error
	LoadStructFromItemByTitleFunc func(config interface{}, itemTitle string, vaultQuery string) error
	LoadStructFromItemByUUIDFunc  func(config interface{}, itemUUID string, vaultQuery string) error
}

var (
	DoGetVaultsFunc                 func() ([]onepassword.Vault, error)
	DoGetVaultsByTitleFunc          func(title string) ([]onepassword.Vault, error)
	DoGetItemFunc                   func(uuid string, vaultUUID string) (*onepassword.Item, error)
	DoGetItemByTitleFunc            func(title string, vaultUUID string) (*onepassword.Item, error)
	DoGetItemsByTitleFunc           func(title string, vaultUUID string) ([]onepassword.Item, error)
	DoCreateItemFunc                func(item *onepassword.Item, vaultUUID string) (*onepassword.Item, error)
	DoDeleteItemFunc                func(item *onepassword.Item, vaultUUID string) error
	DoGetItemsFunc                  func(vaultUUID string) ([]onepassword.Item, error)
	DoUpdateItemFunc                func(item *onepassword.Item, vaultUUID string) (*onepassword.Item, error)
	DoDeleteItemByIDFunc            func(itemUUID, vaultUUID string) error
	DoDeleteItemByTitleFunc         func(itemName, vaultUUID string) error
	DoDownloadFileFunc              func(file *onepassword.File, targetDirectory string, overwrite bool) (string, error)
	DoGetFileFunc                   func(uuid string, itemQuery string, vaultQuery string) (*onepassword.File, error)
	DoGetFileContentFunc            func(file *onepassword.File) ([]byte, error)
	DoGetFilesFunc                  func(itemQuery string, vaultQuery string) ([]onepassword.File, error)
	DoGetItemByUUIDFunc             func(uuid string, vaultQuery string) (*onepassword.Item, error)
	DoGetVaultFunc                  func(uuid string) (*onepassword.Vault, error)
	DoGetVaultByTitleFunc           func(title string) (*onepassword.Vault, error)
	DoGetVaultByUUIDFunc            func(uuid string) (*onepassword.Vault, error)
	DoLoadStructFunc                func(config interface{}) error
	DoLoadStructFromItemFunc        func(config interface{}, itemQuery string, vaultQuery string) error
	DoLoadStructFromItemByTitleFunc func(config interface{}, itemTitle string, vaultQuery string) error
	DoLoadStructFromItemByUUIDFunc  func(config interface{}, itemUUID string, vaultQuery string) error
)

// Do is the mock client's `Do` func
func (m *TestClient) GetVaults() ([]onepassword.Vault, error) {
	return DoGetVaultsFunc()
}

func (m *TestClient) GetVaultsByTitle(title string) ([]onepassword.Vault, error) {
	return DoGetVaultsByTitleFunc(title)
}

func (m *TestClient) GetItem(uuid string, vaultUUID string) (*onepassword.Item, error) {
	return DoGetItemFunc(uuid, vaultUUID)
}

func (m *TestClient) GetItemsByTitle(title, vaultUUID string) ([]onepassword.Item, error) {
	return DoGetItemsByTitleFunc(title, vaultUUID)
}

func (m *TestClient) GetItems(vaultUUID string) ([]onepassword.Item, error) {
	return DoGetItemsFunc(vaultUUID)
}

func (m *TestClient) GetItemByTitle(title string, vaultUUID string) (*onepassword.Item, error) {
	return DoGetItemByTitleFunc(title, vaultUUID)
}

func (m *TestClient) CreateItem(item *onepassword.Item, vaultUUID string) (*onepassword.Item, error) {
	return DoCreateItemFunc(item, vaultUUID)
}

func (m *TestClient) DeleteItem(item *onepassword.Item, vaultUUID string) error {
	return DoDeleteItemFunc(item, vaultUUID)
}

func (m *TestClient) DeleteItemByID(itemUUID, vaultUUID string) error {
	return DoDeleteItemByIDFunc(itemUUID, vaultUUID)
}

func (m *TestClient) DeleteItemByTitle(itemName, vaultUUID string) error {
	return DoDeleteItemByTitleFunc(itemName, vaultUUID)
}

func (m *TestClient) DownloadFile(file *onepassword.File, targetDirectory string, overwrite bool) (string, error) {
	return DoDownloadFileFunc(file, targetDirectory, overwrite)
}

func (m *TestClient) UpdateItem(item *onepassword.Item, vaultUUID string) (*onepassword.Item, error) {
	return DoUpdateItemFunc(item, vaultUUID)
}

func (m *TestClient) GetFile(uuid string, itemQuery string, vaultQuery string) (*onepassword.File, error) {
	return DoGetFileFunc(uuid, itemQuery, vaultQuery)
}

func (m *TestClient) GetFileContent(file *onepassword.File) ([]byte, error) {
	return DoGetFileContentFunc(file)
}

func (m *TestClient) GetFiles(itemQuery string, vaultQuery string) ([]onepassword.File, error) {
	return DoGetFilesFunc(itemQuery, vaultQuery)
}

func (m *TestClient) GetItemByUUID(uuid string, vaultQuery string) (*onepassword.Item, error) {
	return DoGetItemByUUIDFunc(uuid, vaultQuery)
}

func (m *TestClient) GetVault(uuid string) (*onepassword.Vault, error) {
	return DoGetVaultFunc(uuid)
}

func (m *TestClient) GetVaultByTitle(title string) (*onepassword.Vault, error) {
	return DoGetVaultByTitleFunc(title)
}

func (m *TestClient) GetVaultByUUID(uuid string) (*onepassword.Vault, error) {
	return DoGetVaultByUUIDFunc(uuid)
}

func (m *TestClient) LoadStruct(config interface{}) error {
	return DoLoadStructFunc(config)
}

func (m *TestClient) LoadStructFromItem(config interface{}, itemQuery string, vaultQuery string) error {
	return DoLoadStructFromItemFunc(config, itemQuery, vaultQuery)
}

func (m *TestClient) LoadStructFromItemByTitle(config interface{}, itemTitle string, vaultQuery string) error {
	return DoLoadStructFromItemByTitleFunc(config, itemTitle, vaultQuery)
}

func (m *TestClient) LoadStructFromItemByUUID(config interface{}, itemUUID string, vaultQuery string) error {
	return DoLoadStructFromItemByUUIDFunc(config, itemUUID, vaultQuery)
}
