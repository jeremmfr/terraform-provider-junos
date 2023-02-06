package providerfwk

import "github.com/hashicorp/terraform-plugin-framework/diag"

func diagWarns(summary string, warns []error) []diag.Diagnostic {
	diags := make([]diag.Diagnostic, len(warns))
	for i, w := range warns {
		diags[i] = diag.NewWarningDiagnostic(summary, w.Error())
	}

	return diags
}
