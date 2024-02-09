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
				if b.Name.ValueString() != "Second" ||
					b.Value.ValueString() != "test2" {
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
				if !b.Name.IsNull() ||
					!b.Name2.IsNull() ||
					!b.Value.IsNull() {
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
				if !b.Name.IsNull() ||
					!b.Name2.IsNull() ||
					!b.Value.IsNull() {
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
				if !b.Name.IsNull() ||
					!b.Name2.IsNull() ||
					!b.Value.IsNull() {
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

func TestExtractBlockWith2TFTypesString(t *testing.T) {
	t.Parallel()

	type block struct {
		Name    types.String
		NameBis types.String
		Name2   types.Int64
		Value   types.String
	}

	type testCase struct {
		val             []block
		field1Name      string
		newVal1         string
		field2Name      string
		newVal2         string
		expectValLength int
		validateNewVal  func(block)
	}

	tests := map[string]testCase{
		"In List": {
			val: []block{
				{
					Name:    types.StringValue("First"),
					NameBis: types.StringValue("FirstBis"),
					Value:   types.StringValue("test1"),
				}, {
					Name:    types.StringValue("Second"),
					NameBis: types.StringValue("SecondBis"),
					Value:   types.StringValue("test2"),
				}, {
					Name:    types.StringValue("Third"),
					NameBis: types.StringValue("ThirdBis"),
					Value:   types.StringValue("test3"),
				},
			},
			field1Name:      "Name",
			newVal1:         "Second",
			field2Name:      "NameBis",
			newVal2:         "SecondBis",
			expectValLength: 2,
			validateNewVal: func(b block) {
				if b.Name.ValueString() != "Second" ||
					b.NameBis.ValueString() != "SecondBis" ||
					b.Value.ValueString() != "test2" {
					t.Errorf("expected newBlock has second block data")
				}
			},
		},
		"Not In List": {
			val: []block{
				{
					Name:    types.StringValue("First"),
					NameBis: types.StringValue("FirstBis"),
					Value:   types.StringValue("test1"),
				}, {
					Name:    types.StringValue("Second"),
					NameBis: types.StringValue("SecondBis"),
					Value:   types.StringValue("test2"),
				}, {
					Name:    types.StringValue("Third"),
					NameBis: types.StringValue("ThirdBis"),
					Value:   types.StringValue("test3"),
				},
			},
			field1Name:      "Name",
			newVal1:         "Fourth",
			field2Name:      "NameBis",
			newVal2:         "SecondBis",
			expectValLength: 3,
			validateNewVal: func(b block) {
				if !b.Name.IsNull() ||
					!b.NameBis.IsNull() ||
					!b.Name2.IsNull() ||
					!b.Value.IsNull() {
					t.Errorf("expected newBlock has empty data")
				}
			},
		},
		"Not In List Bis": {
			val: []block{
				{
					Name:    types.StringValue("First"),
					NameBis: types.StringValue("FirstBis"),
					Value:   types.StringValue("test1"),
				}, {
					Name:    types.StringValue("Second"),
					NameBis: types.StringValue("SecondBis"),
					Value:   types.StringValue("test2"),
				}, {
					Name:    types.StringValue("Third"),
					NameBis: types.StringValue("ThirdBis"),
					Value:   types.StringValue("test3"),
				},
			},
			field1Name:      "Name",
			newVal1:         "Second",
			field2Name:      "NameBis",
			newVal2:         "Third",
			expectValLength: 3,
			validateNewVal: func(b block) {
				if !b.Name.IsNull() ||
					!b.NameBis.IsNull() ||
					!b.Name2.IsNull() ||
					!b.Value.IsNull() {
					t.Errorf("expected newBlock has empty data")
				}
			},
		},
		"Invalid structFieldName": {
			val: []block{
				{
					Name:    types.StringValue("First"),
					NameBis: types.StringValue("FirstBis"),
					Value:   types.StringValue("test1"),
				}, {
					Name:    types.StringValue("Second"),
					NameBis: types.StringValue("SecondBis"),
					Value:   types.StringValue("test2"),
				}, {
					Name:    types.StringValue("Third"),
					NameBis: types.StringValue("ThirdBis"),
					Value:   types.StringValue("test3"),
				},
			},
			field1Name:      "Nam",
			newVal1:         "Second",
			field2Name:      "NameBis",
			newVal2:         "SecondBis",
			expectValLength: 3,
			validateNewVal: func(b block) {
				if !b.Name.IsNull() ||
					!b.NameBis.IsNull() ||
					!b.Name2.IsNull() ||
					!b.Value.IsNull() {
					t.Errorf("expected newBlock has empty data")
				}
			},
		},
		"Invalid structFieldName Bis": {
			val: []block{
				{
					Name:    types.StringValue("First"),
					NameBis: types.StringValue("FirstBis"),
					Value:   types.StringValue("test1"),
				}, {
					Name:    types.StringValue("Second"),
					NameBis: types.StringValue("SecondBis"),
					Value:   types.StringValue("test2"),
				}, {
					Name:    types.StringValue("Third"),
					NameBis: types.StringValue("ThirdBis"),
					Value:   types.StringValue("test3"),
				},
			},
			field1Name:      "Name",
			newVal1:         "Second",
			field2Name:      "NamBis",
			newVal2:         "SecondBis",
			expectValLength: 3,
			validateNewVal: func(b block) {
				if !b.Name.IsNull() ||
					!b.NameBis.IsNull() ||
					!b.Name2.IsNull() ||
					!b.Value.IsNull() {
					t.Errorf("expected newBlock has empty data")
				}
			},
		},
		"Invalid type of structFieldName": {
			val: []block{
				{
					Name:    types.StringValue("First"),
					NameBis: types.StringValue("FirstBis"),
					Name2:   types.Int64Value(1),
					Value:   types.StringValue("test1"),
				}, {
					Name:    types.StringValue("Second"),
					NameBis: types.StringValue("SecondBis"),
					Name2:   types.Int64Value(2),
					Value:   types.StringValue("test2"),
				}, {
					Name:    types.StringValue("Third"),
					NameBis: types.StringValue("ThirdBis"),
					Name2:   types.Int64Value(3),
					Value:   types.StringValue("test3"),
				},
			},
			field1Name:      "Name2",
			newVal1:         "Second",
			field2Name:      "NameBis",
			newVal2:         "SecondBis",
			expectValLength: 3,
			validateNewVal: func(b block) {
				if !b.Name.IsNull() ||
					!b.NameBis.IsNull() ||
					!b.Name2.IsNull() ||
					!b.Value.IsNull() {
					t.Errorf("expected newBlock has empty data")
				}
			},
		},
		"Invalid type of structFieldName bis": {
			val: []block{
				{
					Name:    types.StringValue("First"),
					NameBis: types.StringValue("FirstBis"),
					Name2:   types.Int64Value(1),
					Value:   types.StringValue("test1"),
				}, {
					Name:    types.StringValue("Second"),
					NameBis: types.StringValue("SecondBis"),
					Name2:   types.Int64Value(2),
					Value:   types.StringValue("test2"),
				}, {
					Name:    types.StringValue("Third"),
					NameBis: types.StringValue("ThirdBis"),
					Name2:   types.Int64Value(3),
					Value:   types.StringValue("test3"),
				},
			},
			field1Name:      "Name",
			newVal1:         "Second",
			field2Name:      "Name2",
			newVal2:         "SecondBis",
			expectValLength: 3,
			validateNewVal: func(b block) {
				if !b.Name.IsNull() ||
					!b.NameBis.IsNull() ||
					!b.Name2.IsNull() ||
					!b.Value.IsNull() {
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
			blocks, block = tfdata.ExtractBlockWith2TFTypesString(
				blocks,
				test.field1Name, test.newVal1,
				test.field2Name, test.newVal2,
			)

			if v := len(blocks); v != test.expectValLength {
				t.Errorf("the expected block length is %d, got %d", test.expectValLength, v)
			}
			test.validateNewVal(block)
		})
	}
}

func TestExtractBlockWithTFTypesInt64(t *testing.T) {
	t.Parallel()

	type block struct {
		Name  types.Int64
		Name2 types.String
		Value types.String
	}

	type testCase struct {
		val             []block
		newVal          int64
		fieldName       string
		expectValLength int
		validateNewVal  func(block)
	}

	tests := map[string]testCase{
		"In List": {
			val: []block{
				{
					Name:  types.Int64Value(1),
					Value: types.StringValue("test1"),
				}, {
					Name:  types.Int64Value(2),
					Value: types.StringValue("test2"),
				}, {
					Name:  types.Int64Value(3),
					Value: types.StringValue("test3"),
				},
			},
			fieldName:       "Name",
			newVal:          2,
			expectValLength: 2,
			validateNewVal: func(b block) {
				if b.Name.ValueInt64() != 2 ||
					b.Value.ValueString() != "test2" {
					t.Errorf("expected newBlock has second block data")
				}
			},
		},
		"Not In List": {
			val: []block{
				{
					Name:  types.Int64Value(1),
					Value: types.StringValue("test1"),
				}, {
					Name:  types.Int64Value(2),
					Value: types.StringValue("test2"),
				}, {
					Name:  types.Int64Value(3),
					Value: types.StringValue("test3"),
				},
			},
			fieldName:       "Name",
			newVal:          4,
			expectValLength: 3,
			validateNewVal: func(b block) {
				if !b.Name.IsNull() ||
					!b.Name2.IsNull() ||
					!b.Value.IsNull() {
					t.Errorf("expected newBlock has empty data")
				}
			},
		},
		"Invalid structFieldName": {
			val: []block{
				{
					Name:  types.Int64Value(1),
					Value: types.StringValue("test1"),
				}, {
					Name:  types.Int64Value(2),
					Value: types.StringValue("test2"),
				}, {
					Name:  types.Int64Value(3),
					Value: types.StringValue("test3"),
				},
			},
			fieldName:       "Nam",
			newVal:          2,
			expectValLength: 3,
			validateNewVal: func(b block) {
				if !b.Name.IsNull() ||
					!b.Name2.IsNull() ||
					!b.Value.IsNull() {
					t.Errorf("expected newBlock has empty data")
				}
			},
		},
		"Invalid type of structFieldName": {
			val: []block{
				{
					Name:  types.Int64Value(1),
					Name2: types.StringValue("First"),
					Value: types.StringValue("test1"),
				}, {
					Name:  types.Int64Value(2),
					Name2: types.StringValue("Second"),
					Value: types.StringValue("test2"),
				}, {
					Name:  types.Int64Value(3),
					Name2: types.StringValue("Third"),
					Value: types.StringValue("test3"),
				},
			},
			fieldName:       "Name2",
			newVal:          2,
			expectValLength: 3,
			validateNewVal: func(b block) {
				if !b.Name.IsNull() ||
					!b.Name2.IsNull() ||
					!b.Value.IsNull() {
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
			blocks, block = tfdata.ExtractBlockWithTFTypesInt64(blocks, test.fieldName, test.newVal)

			if v := len(blocks); v != test.expectValLength {
				t.Errorf("the expected block length is %d, got %d", test.expectValLength, v)
			}
			test.validateNewVal(block)
		})
	}
}
