package onepassword

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/1Password/connect-sdk-go/onepassword"
	"github.com/hashicorp/vault/sdk/logical"
)

func TestItemsPath(t *testing.T) {
	t.Parallel()

	b, reqStorage := getTestBackendWithCachedClient(t)

	item1 := generateVaultItem()
	expectedCategory := item1["category"].(string)
	id1 := testAddItem(t, b, reqStorage, item1, expectedCategory)
	testReadItem(t, b, reqStorage, item1, id1)

	item2 := generateVaultItem()
	id2 := testAddItem(t, b, reqStorage, item2, expectedCategory)

	expectedListedItems := map[string]interface{}{
		"keys": []string{id1, id2},
		"keys_info": map[string]interface{}{
			id1: item1["title"],
			id2: item2["title"],
		},
	}
	testListItems(t, b, reqStorage, expectedListedItems)
	testDeleteItem(t, b, reqStorage, id2)
	testReadNonExistentItem(t, b, reqStorage, id2)

	item3 := generateVaultItem()
	testUpdateItem(t, b, reqStorage, item3, expectedCategory, id1)
	testReadItem(t, b, reqStorage, item3, id1)
}

func testReadNonExistentItem(t *testing.T, b logical.Backend, s logical.Storage, id string) {
	_, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Path:      fmt.Sprintf("vaults/hfnjvi6aymbsnfc2xeeoheizda/items/%v", id),
		Storage:   s,
	})
	if err == nil {
		t.Errorf("Retrieving non existent item did not cause the expected error")
	}
}

func testAddItem(t *testing.T, b logical.Backend, s logical.Storage, d map[string]interface{}, expectedCategory string) string {
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "vaults/hfnjvi6aymbsnfc2xeeoheizda/items/",
		Data:      d,
		Storage:   s,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.IsError() {
		t.Fatal(resp.Error())
	}

	if resp.Data["ID"] == nil {
		t.Errorf("No ID returned on Item Creation")
	}

	if resp.Data["createdAt"] == nil {
		t.Errorf("No createdAt returned on Item Creation")
	}

	savedCategory := strings.ToLower(fmt.Sprintf("%v", resp.Data["category"]))
	if savedCategory != expectedCategory {
		t.Errorf("Expected category to be %v was %v", expectedCategory, savedCategory)
	}
	return fmt.Sprintf("%v", resp.Data["ID"])
}

func testUpdateItem(t *testing.T, b logical.Backend, s logical.Storage, d map[string]interface{}, expectedCategory, id string) string {
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.CreateOperation,
		Path:      fmt.Sprintf("vaults/hfnjvi6aymbsnfc2xeeoheizda/items/%v", id),
		Data:      d,
		Storage:   s,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.IsError() {
		t.Fatal(resp.Error())
	}

	if resp.Data["ID"] == nil {
		t.Errorf("No ID returned on Item Creation")
	}

	if resp.Data["createdAt"] == nil {
		t.Errorf("No createdAt returned on Item Creation")
	}

	savedCategory := strings.ToLower(fmt.Sprintf("%v", resp.Data["category"]))
	if savedCategory != expectedCategory {
		t.Errorf("Expected category to be %v was %v", expectedCategory, savedCategory)
	}
	return fmt.Sprintf("%v", resp.Data["ID"])
}

func testDeleteItem(t *testing.T, b logical.Backend, s logical.Storage, id string) {
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.DeleteOperation,
		Path:      fmt.Sprintf("vaults/hfnjvi6aymbsnfc2xeeoheizda/items/%v", id),
		Storage:   s,
	})

	if err != nil {
		t.Fatal(err)
	}

	if resp == nil {
		return
	}

	if resp.IsError() {
		t.Fatal(resp.Error())
	}

	if t.Failed() {
		t.FailNow()
	}
}

func testReadItem(t *testing.T, b logical.Backend, s logical.Storage, expected map[string]interface{}, id string) {
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Path:      fmt.Sprintf("vaults/hfnjvi6aymbsnfc2xeeoheizda/items/%v", id),
		Storage:   s,
	})

	if err != nil {
		t.Fatal(err)
	}

	if resp == nil && expected == nil {
		return
	}

	if resp.IsError() {
		t.Fatal(resp.Error())
	}

	fields := expected["fields"].([]map[string]interface{})
	if len(fields) != len(resp.Data) {
		t.Errorf("read data mismatch (expected %d values, got %d)", len(fields), len(resp.Data))
	}

	for _, field := range fields {
		itemField := onepassword.ItemField{}
		marshalledItem, err := json.Marshal(field)
		if err != nil {
			t.Errorf("Saved field is not valid")
		}
		err = json.Unmarshal(marshalledItem, &itemField)
		if err != nil {
			t.Errorf("Saved field is not valid")
		}
		if resp.Data[itemField.Label] != itemField.Value {
			t.Errorf("Expected data[%v] to be %v was %v", itemField.Label, itemField.Value, resp.Data[itemField.Label])
		}
	}

	if t.Failed() {
		t.FailNow()
	}
}

func testListItems(t *testing.T, b logical.Backend, s logical.Storage, expected map[string]interface{}) {
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ListOperation,
		Path:      "vaults/hfnjvi6aymbsnfc2xeeoheizda/items",
		Storage:   s,
	})

	if err != nil {
		t.Fatal(err)
	}

	if resp == nil && expected == nil {
		return
	}

	if resp.IsError() {
		t.Fatal(resp.Error())
	}

	actual := resp.Data

	if len(expected) != len(actual) {
		t.Errorf("read data mismatch (expected %d values, got %v)", len(expected), actual["keys"].([]string))
	}

	for _, expectedV := range expected["keys"].([]string) {
		found := false
		for _, actualV := range actual["keys"].([]string) {
			if actualV == expectedV {
				found = true
			}
		}

		if !found {
			t.Errorf(`expected data["keys"] = %v but was %v"`, expected["keys"].([]string), actual["keys"].([]string))
		}
	}

	if t.Failed() {
		t.FailNow()
	}
}

func generateVaultItem() map[string]interface{} {
	return map[string]interface{}{
		"category": "database",
		"title":    "Test Item",
		"fields": []map[string]interface{}{
			{
				"id":    "some_id",
				"type":  "STRING",
				"label": "some title",
				"value": "some vlaue",
			},
			{
				"id":      "username",
				"label":   "username",
				"type":    "STRING",
				"purpose": "username",
				"value":   "new_user",
			},
		},
		"sections": []map[string]interface{}{
			{
				"id":    "new_section",
				"label": "New Section",
			},
		},
	}
}
