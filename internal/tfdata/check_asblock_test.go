package tfdata_test

import (
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestCheckBlockIsEmpty(t *testing.T) {
	t.Parallel()

	type block2 struct {
		SubStringAttr     types.String   `tfsdk:"sub_string_attr"`
		SubInt64Attr      types.Int64    `tfsdk:"sub_int64_attr"`
		SubStringListAttr []types.String `tfsdk:"sub_string_list_attr"`
	}

	type block struct {
		StringAttr     types.String   `tfsdk:"string_attr"`
		Int64Attr      types.Int64    `tfsdk:"int64_attr"`
		StringListAttr []types.String `tfsdk:"string_list_attr"`
		BlockAttr      *block2        `tfsdk:"block_attr"`
		StructAttr     *struct {
			SubBoolAttr types.Bool `tfsdk:"sub_bool_attr"`
		} `tfsdk:"struct_attr"`
	}

	type testCase struct {
		val            *block
		excludeFields  []string
		expectResponse bool
	}

	tests := map[string]testCase{
		"empty": {
			val:            new(block),
			expectResponse: true,
		},
		"all_null": {
			val: &block{
				StringAttr:     types.StringNull(),
				Int64Attr:      types.Int64Null(),
				StringListAttr: nil,
				BlockAttr:      nil,
				StructAttr:     nil,
			},
			expectResponse: true,
		},
		"StringAttr_known": {
			val: &block{
				StringAttr:     types.StringValue("value_string"),
				Int64Attr:      types.Int64Null(),
				StringListAttr: nil,
				BlockAttr:      nil,
				StructAttr:     nil,
			},
			expectResponse: false,
		},
		"Int64Attr_unknown": {
			val: &block{
				StringAttr:     types.StringNull(),
				Int64Attr:      types.Int64Unknown(),
				StringListAttr: nil,
				BlockAttr:      nil,
				StructAttr:     nil,
			},
			expectResponse: false,
		},
		"StringListAttr_set_but_empty": {
			val: &block{
				StringAttr:     types.StringNull(),
				Int64Attr:      types.Int64Null(),
				StringListAttr: make([]types.String, 0),
				BlockAttr:      nil,
				StructAttr:     nil,
			},
			expectResponse: true,
		},
		"StringListAttr_set": {
			val: &block{
				StringAttr: types.StringNull(),
				Int64Attr:  types.Int64Null(),
				StringListAttr: []types.String{
					types.StringValue("value_string_list"),
				},
				BlockAttr:  nil,
				StructAttr: nil,
			},
			expectResponse: false,
		},
		"BlockAttr_set_but_empty": {
			val: &block{
				StringAttr:     types.StringNull(),
				Int64Attr:      types.Int64Null(),
				StringListAttr: nil,
				BlockAttr: &block2{
					SubStringAttr:     types.StringNull(),
					SubInt64Attr:      types.Int64Null(),
					SubStringListAttr: nil,
				},
				StructAttr: nil,
			},
			expectResponse: false,
		},
		"BlockAttr_set": {
			val: &block{
				StringAttr:     types.StringNull(),
				Int64Attr:      types.Int64Null(),
				StringListAttr: nil,
				BlockAttr: &block2{
					SubStringAttr:     types.StringValue("value_sub_string"),
					SubInt64Attr:      types.Int64Null(),
					SubStringListAttr: nil,
				},
				StructAttr: nil,
			},
			expectResponse: false,
		},
		"StrucAttr_set_but_empty": {
			val: &block{
				StringAttr:     types.StringNull(),
				Int64Attr:      types.Int64Null(),
				StringListAttr: nil,
				BlockAttr:      nil,
				StructAttr: &struct {
					SubBoolAttr types.Bool `tfsdk:"sub_bool_attr"`
				}{
					SubBoolAttr: types.BoolNull(),
				},
			},
			expectResponse: false,
		},
		"StrucAttr_set": {
			val: &block{
				StringAttr:     types.StringNull(),
				Int64Attr:      types.Int64Null(),
				StringListAttr: nil,
				BlockAttr:      nil,
				StructAttr: &struct {
					SubBoolAttr types.Bool `tfsdk:"sub_bool_attr"`
				}{
					SubBoolAttr: types.BoolValue(true),
				},
			},
			expectResponse: false,
		},
		"StringAttr_set_but_exclude": {
			val: &block{
				StringAttr:     types.StringValue("value_string"),
				Int64Attr:      types.Int64Null(),
				StringListAttr: nil,
				BlockAttr:      nil,
				StructAttr:     nil,
			},
			excludeFields:  []string{"StringAttr"},
			expectResponse: true,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if resp := tfdata.CheckBlockIsEmpty(test.val, test.excludeFields...); resp != test.expectResponse {
				t.Errorf("the expected response %v, got %v", test.expectResponse, resp)
			}
		})
	}
}

func TestCheckBlockHasKnownValue(t *testing.T) {
	t.Parallel()

	type block2 struct {
		SubStringAttr types.String `tfsdk:"sub_string_attr"`
		SubInt64Attr  types.Int64  `tfsdk:"sub_int64_attr"`
	}

	type block struct {
		StringAttr types.String `tfsdk:"string_attr"`
		Int64Attr  types.Int64  `tfsdk:"int64_attr"`
		BlockAttr  *block2      `tfsdk:"block_attr"`
		StructAttr *struct {
			SubBoolAttr types.Bool `tfsdk:"sub_bool_attr"`
		} `tfsdk:"struct_attr"`
	}

	type testCase struct {
		val            *block
		expectResponse bool
	}

	tests := map[string]testCase{
		"empty": {
			val:            new(block),
			expectResponse: false,
		},
		"all_null": {
			val: &block{
				StringAttr: types.StringNull(),
				Int64Attr:  types.Int64Null(),
				BlockAttr:  nil,
				StructAttr: nil,
			},
			expectResponse: false,
		},
		"StringAttr_known": {
			val: &block{
				StringAttr: types.StringValue("value_string"),
				Int64Attr:  types.Int64Null(),
				BlockAttr:  nil,
				StructAttr: nil,
			},
			expectResponse: true,
		},
		"Int64Attr_unknown": {
			val: &block{
				StringAttr: types.StringNull(),
				Int64Attr:  types.Int64Unknown(),
				BlockAttr:  nil,
				StructAttr: nil,
			},
			expectResponse: false,
		},
		"BlockAttr_set_empty": {
			val: &block{
				StringAttr: types.StringNull(),
				Int64Attr:  types.Int64Unknown(),
				BlockAttr:  new(block2),
				StructAttr: nil,
			},
			expectResponse: false,
		},
		"BlockAttr_set_null": {
			val: &block{
				StringAttr: types.StringNull(),
				Int64Attr:  types.Int64Unknown(),
				BlockAttr: &block2{
					SubStringAttr: types.StringNull(),
					SubInt64Attr:  types.Int64Null(),
				},
				StructAttr: nil,
			},
			expectResponse: false,
		},
		"BlockAttr_set_unknown": {
			val: &block{
				StringAttr: types.StringNull(),
				Int64Attr:  types.Int64Unknown(),
				BlockAttr: &block2{
					SubStringAttr: types.StringUnknown(),
					SubInt64Attr:  types.Int64Null(),
				},
				StructAttr: nil,
			},
			expectResponse: false,
		},
		"BlockAttr_set_known": {
			val: &block{
				StringAttr: types.StringNull(),
				Int64Attr:  types.Int64Unknown(),
				BlockAttr: &block2{
					SubStringAttr: types.StringValue("sub_value_string"),
					SubInt64Attr:  types.Int64Unknown(),
				},
				StructAttr: nil,
			},
			expectResponse: true,
		},
		"StrucAttr_set_null": {
			val: &block{
				StringAttr: types.StringNull(),
				Int64Attr:  types.Int64Unknown(),
				BlockAttr:  nil,
				StructAttr: &struct {
					SubBoolAttr types.Bool `tfsdk:"sub_bool_attr"`
				}{
					SubBoolAttr: types.BoolNull(),
				},
			},
			expectResponse: false,
		},
		"StrucAttr_set_unknown": {
			val: &block{
				StringAttr: types.StringNull(),
				Int64Attr:  types.Int64Unknown(),
				BlockAttr:  nil,
				StructAttr: &struct {
					SubBoolAttr types.Bool `tfsdk:"sub_bool_attr"`
				}{
					SubBoolAttr: types.BoolUnknown(),
				},
			},
			expectResponse: false,
		},
		"StrucAttr_set_known": {
			val: &block{
				StringAttr: types.StringNull(),
				Int64Attr:  types.Int64Unknown(),
				BlockAttr:  nil,
				StructAttr: &struct {
					SubBoolAttr types.Bool `tfsdk:"sub_bool_attr"`
				}{
					SubBoolAttr: types.BoolValue(true),
				},
			},
			expectResponse: true,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if resp := tfdata.CheckBlockHasKnownValue(test.val); resp != test.expectResponse {
				t.Errorf("the expected response %v, got %v", test.expectResponse, resp)
			}
		})
	}
}
