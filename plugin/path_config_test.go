package onepassword

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
)

const (
	opToken = "test_token"
	opVault = "hfnjvi6aymbsnfc2xeeoheizda"
	host    = "localhost:8080"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	b, reqStorage := getTestBackendNoCachedClient(t)

	testConfigRead(t, b, reqStorage, nil)

	config := map[string]interface{}{
		"op_connect_token": opToken,
		"op_connect_host":  host,
		"op_vault":         opVault,
	}
	testConfigUpdate(t, b, reqStorage, config)

	testConfigRead(t, b, reqStorage, config)

	updatedHost := "updatedHost"
	testConfigUpdate(t, b, reqStorage, map[string]interface{}{
		"op_connect_host": updatedHost,
	})

	config["op_connect_host"] = updatedHost
	testConfigRead(t, b, reqStorage, config)
}

func testConfigUpdate(t *testing.T, b logical.Backend, s logical.Storage, d map[string]interface{}) {
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Data:      d,
		Storage:   s,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.IsError() {
		t.Fatal(resp.Error())
	}
}

func testConfigRead(t *testing.T, b logical.Backend, s logical.Storage, expected map[string]interface{}) {
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "config",
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

	if len(expected) != len(resp.Data) {
		t.Errorf("read data mismatch (expected %d values, got %d)", len(expected), len(resp.Data))
	}

	for k, expectedV := range expected {
		actualV, ok := resp.Data[k]

		if !ok {
			t.Errorf(`expected data["%s"] = %v but was not included in read output"`, k, expectedV)
		} else if expectedV != actualV {
			t.Errorf(`expected data["%s"] = %v, instead got %v"`, k, expectedV, actualV)
		}
	}

	if t.Failed() {
		t.FailNow()
	}
}
