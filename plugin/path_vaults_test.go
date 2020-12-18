package onepassword

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
)

func TestListVaults(t *testing.T) {
	t.Parallel()

	b, reqStorage := getTestBackendWithCachedClient(t)

	key1 := fmt.Sprintf("%s %s", vaultName1, vaultId1)
	key2 := fmt.Sprintf("%s %s", vaultName2, vaultId2)

	expected := map[string]interface{}{
		"keys": []string{key1, key2},
		"keys_info": map[string]interface{}{
			key1: vaultId1,
			key2: vaultId2,
		},
	}
	testListVaults(t, b, reqStorage, expected)
}

func testListVaults(t *testing.T, b logical.Backend, s logical.Storage, expected map[string]interface{}) {
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ListOperation,
		Path:      "vaults",
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
