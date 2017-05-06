package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatal("err: %s", err)
	}
}

func TestProviderImplementation(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}
