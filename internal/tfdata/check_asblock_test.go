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

	type blockWithEmbed struct {
		block
		String2Attr types.String `tfsdk:"string2_attr"`
	}

	type testCase struct {
		val            *block
		val2           *blockWithEmbed
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
		"blockEmbeded_null": {
			val2: &blockWithEmbed{
				block: block{
					StringAttr: types.StringNull(),
				},
			},
			expectResponse: true,
		},
		"blockEmbeded_unknown": {
			val2: &blockWithEmbed{
				block: block{
					StringAttr: types.StringUnknown(),
				},
			},
			expectResponse: false,
		},
		"blockEmbeded_embededSet": {
			val2: &blockWithEmbed{
				block: block{
					StringAttr: types.StringValue("value_string"),
				},
			},
			expectResponse: false,
		},
		"blockEmbeded_set": {
			val2: &blockWithEmbed{
				block:       block{},
				String2Attr: types.StringValue("value_string"),
			},
			expectResponse: false,
		},
		"blockEmbeded_embedSet_but_exclude": {
			val2: &blockWithEmbed{
				block: block{
					StringAttr: types.StringValue("value_string"),
				},
				String2Attr: types.StringNull(),
			},
			excludeFields:  []string{"StringAttr"},
			expectResponse: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			switch {
			case test.val != nil:
				if resp := tfdata.CheckBlockIsEmpty(test.val, test.excludeFields...); resp != test.expectResponse {
					t.Errorf("the expected response %v, got %v", test.expectResponse, resp)
				}
			case test.val2 != nil:
				if resp := tfdata.CheckBlockIsEmpty(test.val2, test.excludeFields...); resp != test.expectResponse {
					t.Errorf("the expected response %v, got %v", test.expectResponse, resp)
				}
			default:
				t.Error("nil vals")
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

	type blockWithEmbed struct {
		block
		String2Attr types.String `tfsdk:"string2_attr"`
	}

	type testCase struct {
		val            *block
		val2           *blockWithEmbed
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
		"blockEmbeded_null": {
			val2: &blockWithEmbed{
				block: block{
					StringAttr: types.StringNull(),
				},
			},
			expectResponse: false,
		},
		"blockEmbeded_unknown": {
			val2: &blockWithEmbed{
				block: block{
					StringAttr: types.StringUnknown(),
				},
			},
			expectResponse: false,
		},
		"blockEmbeded_embededSet": {
			val2: &blockWithEmbed{
				block: block{
					StringAttr: types.StringValue("value_string"),
				},
			},
			expectResponse: true,
		},
		"blockEmbeded_set": {
			val2: &blockWithEmbed{
				block:       block{},
				String2Attr: types.StringValue("value_string"),
			},
			expectResponse: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			switch {
			case test.val != nil:
				if resp := tfdata.CheckBlockHasKnownValue(test.val); resp != test.expectResponse {
					t.Errorf("the expected response %v, got %v", test.expectResponse, resp)
				}
			case test.val2 != nil:
				if resp := tfdata.CheckBlockHasKnownValue(test.val2); resp != test.expectResponse {
					t.Errorf("the expected response %v, got %v", test.expectResponse, resp)
				}
			default:
				t.Error("nil vals")
			}
		})
	}
}
