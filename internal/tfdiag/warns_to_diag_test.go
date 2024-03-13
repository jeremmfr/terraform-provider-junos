package tfdiag_test

import (
	"errors"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func TestWarns(t *testing.T) {
	t.Parallel()

	err := errors.New("test")
	diags := tfdiag.Warns("Test", []error{err, err})

	if v := len(diags); v != 2 {
		t.Errorf("got unexpected diag length: %d", v)
	} else {
		if v := diags[0].Severity(); v != diag.SeverityWarning {
			t.Errorf("got unexpected severity: %s", v)
		}
		if v := diags[1].Severity(); v != diag.SeverityWarning {
			t.Errorf("got unexpected severity: %s", v)
		}
		if v := diags[0].Summary(); v != "Test" {
			t.Errorf("got unexpected summary: %s", v)
		}
		if v := diags[1].Summary(); v != "Test" {
			t.Errorf("got unexpected summary: %s", v)
		}
		if v := diags[0].Detail(); v != "test" {
			t.Errorf("got unexpected detail: %s", v)
		}
		if v := diags[1].Detail(); v != "test" {
			t.Errorf("got unexpected detail: %s", v)
		}
	}
}
