package providersdk_test

import (
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/providersdk"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestProvider(t *testing.T) {
	if err := providersdk.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(_ *testing.T) {
	var _ *schema.Provider = providersdk.Provider()
}
