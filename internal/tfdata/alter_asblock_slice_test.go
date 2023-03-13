package tfdata_test

import (
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestExtractBlockWithTFTypesString(t *testing.T) {
	t.Parallel()

	type block struct {
		Name  types.String
		Name2 types.Int64
		Value types.String
	}

	type testCase struct {
		val             []block
		newVal          string
		fieldName       string
		expectValLength int
		validateNewVal  func(block)
	}

	tests := map[string]testCase{
		"In List": {
			val: []block{
				{
					Name:  types.StringValue("First"),
					Value: types.StringValue("test1"),
				}, {
					Name:  types.StringValue("Second"),
					Value: types.StringValue("test2"),
				}, {
					Name:  types.StringValue("Third"),
					Value: types.StringValue("test3"),
				},
			},
			fieldName:       "Name",
			newVal:          "Second",
			expectValLength: 2,
			validateNewVal: func(b block) {
				if b.Name.ValueString() != "Second" || b.Value.ValueString() != "test2" {
					t.Errorf("expected newBlock has second block data")
				}
			},
		},
		"Not In List": {
			val: []block{
				{
					Name:  types.StringValue("First"),
					Value: types.StringValue("test1"),
				}, {
					Name:  types.StringValue("Second"),
					Value: types.StringValue("test2"),
				}, {
					Name:  types.StringValue("Third"),
					Value: types.StringValue("test3"),
				},
			},
			fieldName:       "Name",
			newVal:          "Fourth",
			expectValLength: 3,
			validateNewVal: func(b block) {
				if !b.Name.IsNull() || !b.Value.IsNull() ||
					b.Name.ValueString() != "" || b.Value.ValueString() != "" {
					t.Errorf("expected newBlock has empty data")
				}
			},
		},
		"Invalid structFieldName": {
			val: []block{
				{
					Name:  types.StringValue("First"),
					Value: types.StringValue("test1"),
				}, {
					Name:  types.StringValue("Second"),
					Value: types.StringValue("test2"),
				}, {
					Name:  types.StringValue("Third"),
					Value: types.StringValue("test3"),
				},
			},
			fieldName:       "Nam",
			newVal:          "Second",
			expectValLength: 3,
			validateNewVal: func(b block) {
				if !b.Name.IsNull() || !b.Value.IsNull() ||
					b.Name.ValueString() != "" || b.Value.ValueString() != "" {
					t.Errorf("expected newBlock has empty data")
				}
			},
		},
		"Invalid type of structFieldName": {
			val: []block{
				{
					Name:  types.StringValue("First"),
					Name2: types.Int64Value(1),
					Value: types.StringValue("test1"),
				}, {
					Name:  types.StringValue("Second"),
					Name2: types.Int64Value(2),
					Value: types.StringValue("test2"),
				}, {
					Name:  types.StringValue("Third"),
					Name2: types.Int64Value(3),
					Value: types.StringValue("test3"),
				},
			},
			fieldName:       "Name2",
			newVal:          "Second",
			expectValLength: 3,
			validateNewVal: func(b block) {
				if !b.Name.IsNull() || !b.Value.IsNull() || !b.Name2.IsNull() &&
					b.Name.ValueString() != "" || b.Value.ValueString() != "" || b.Name2.ValueInt64() != 0 {
					t.Errorf("expected newBlock has empty data")
				}
			},
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var block block
			blocks := test.val
			blocks, block = tfdata.ExtractBlockWithTFTypesString(blocks, test.fieldName, test.newVal)

			if v := len(blocks); v != test.expectValLength {
				t.Errorf("the expected block length is %d, got %d", test.expectValLength, v)
			}
			test.validateNewVal(block)
		})
	}
}
