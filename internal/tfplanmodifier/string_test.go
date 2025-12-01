package tfplanmodifier_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/tfplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestStringRemoveBlankLines(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		request  planmodifier.StringRequest
		expected *planmodifier.StringResponse
	}{
		"null value": {
			request: planmodifier.StringRequest{
				PlanValue: types.StringNull(),
			},
			expected: &planmodifier.StringResponse{
				PlanValue: types.StringNull(),
			},
		},
		"unknown value": {
			request: planmodifier.StringRequest{
				PlanValue: types.StringUnknown(),
			},
			expected: &planmodifier.StringResponse{
				PlanValue: types.StringUnknown(),
			},
		},
		"single line without newline": {
			request: planmodifier.StringRequest{
				PlanValue: types.StringValue("single line"),
			},
			expected: &planmodifier.StringResponse{
				PlanValue: types.StringValue("single line"),
			},
		},
		"multiline with blank lines": {
			request: planmodifier.StringRequest{
				PlanValue: types.StringValue("line1\n\nline2\n\nline3\n"),
			},
			expected: &planmodifier.StringResponse{
				PlanValue: types.StringValue("line1\nline2\nline3\n"),
			},
		},
		"multiline with only spaces": {
			request: planmodifier.StringRequest{
				PlanValue: types.StringValue("line1\n     \nline2\n  \nline3\n"),
			},
			expected: &planmodifier.StringResponse{
				PlanValue: types.StringValue("line1\nline2\nline3\n"),
			},
		},
		"multiline with tabs and spaces": {
			request: planmodifier.StringRequest{
				PlanValue: types.StringValue("line1\n\t  \t\nline2\n  \t  \nline3\n"),
			},
			expected: &planmodifier.StringResponse{
				PlanValue: types.StringValue("line1\nline2\nline3\n"),
			},
		},
		"all blank lines": {
			request: planmodifier.StringRequest{
				PlanValue: types.StringValue("\n\n   \n\t\n"),
			},
			expected: &planmodifier.StringResponse{
				PlanValue: types.StringValue(""),
			},
		},
		"no blank lines": {
			request: planmodifier.StringRequest{
				PlanValue: types.StringValue("line1\nline2\nline3\n"),
			},
			expected: &planmodifier.StringResponse{
				PlanValue: types.StringValue("line1\nline2\nline3\n"),
			},
		},
		"preserve indentation": {
			request: planmodifier.StringRequest{
				PlanValue: types.StringValue("  line1\n\n    line2\n"),
			},
			expected: &planmodifier.StringResponse{
				PlanValue: types.StringValue("  line1\n    line2\n"),
			},
		},
		"empty string": {
			request: planmodifier.StringRequest{
				PlanValue: types.StringValue(""),
			},
			expected: &planmodifier.StringResponse{
				PlanValue: types.StringValue(""),
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &planmodifier.StringResponse{
				PlanValue: testCase.request.PlanValue,
			}

			tfplanmodifier.StringRemoveBlankLines().PlanModifyString(context.Background(), testCase.request, resp)

			if !reflect.DeepEqual(testCase.expected, resp) {
				t.Errorf("unexpected StringResponse: want %#v, got %#v", testCase.expected, resp)
			}
		})
	}
}
