package tfdiag

import "github.com/hashicorp/terraform-plugin-framework/diag"

func Warns(summary string, warns []error) []diag.Diagnostic {
	diags := make([]diag.Diagnostic, len(warns))
	for i, w := range warns {
		diags[i] = diag.NewWarningDiagnostic(summary, w.Error())
	}

	return diags
}
