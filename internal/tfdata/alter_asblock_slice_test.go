package tfdata_test

import (
	"fmt"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/tfdata"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestExtractBlock(t *testing.T) {
	t.Parallel()

	type block1string struct {
		Name   types.String `tfdata:"identifier"`
		Value  types.String
		Value2 types.Int64
	}

	type block2string struct {
		Name  types.String `tfdata:"identifier_1"`
		Name2 types.String `tfdata:"identifier_2"`
		Value types.String
	}

	type block1int struct {
		Name  types.Int64 `tfdata:"identifier_1"`
		Name2 types.String
		Value types.String
	}

	type blockEmbed1string struct {
		block1string
		Value3 types.String
	}

	type testCase struct {
		val1string      []block1string
		val2string      []block2string
		val1int         []block1int
		valEmbed1string []blockEmbed1string

		newValString1 types.String
		newValString2 types.String
		newValInt1    types.Int64

		expectValLength            int
		validateNewVal1string      func(block1string) error
		validateNewVal2string      func(block2string) error
		validateNewVal1int         func(block1int) error
		validateNewValEmbed1string func(blockEmbed1string) error
	}

	tests := map[string]testCase{
		"1String - Empty": {
			val1string:      []block1string{},
			newValString1:   types.StringValue("Third"),
			expectValLength: 0,
			validateNewVal1string: func(b block1string) error {
				if b.Name.ValueString() != "Third" ||
					!b.Value.IsNull() {
					return fmt.Errorf("expected newBlock has empty data: got %#v", b)
				}

				return nil
			},
		},
		"1String - In List": {
			val1string: []block1string{
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
			newValString1:   types.StringValue("Second"),
			expectValLength: 2,
			validateNewVal1string: func(b block1string) error {
				if b.Name.ValueString() != "Second" ||
					b.Value.ValueString() != "test2" {
					return fmt.Errorf("expected newBlock has second block data: got %#v", b)
				}

				return nil
			},
		},
		"1String - Not In List": {
			val1string: []block1string{
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
			newValString1:   types.StringValue("Fourth"),
			expectValLength: 3,
			validateNewVal1string: func(b block1string) error {
				if !b.Name.Equal(types.StringValue("Fourth")) ||
					!b.Value.IsNull() ||
					!b.Value2.IsNull() {
					return fmt.Errorf("expected newBlock has empty data: got %#v", b)
				}

				return nil
			},
		},
		"2String - In List": {
			val2string: []block2string{
				{
					Name:  types.StringValue("First"),
					Name2: types.StringValue("FirstBis"),
					Value: types.StringValue("test1"),
				}, {
					Name:  types.StringValue("Second"),
					Name2: types.StringValue("SecondBis"),
					Value: types.StringValue("test2"),
				}, {
					Name:  types.StringValue("Third"),
					Name2: types.StringValue("ThirdBis"),
					Value: types.StringValue("test3"),
				},
			},
			newValString1:   types.StringValue("Second"),
			newValString2:   types.StringValue("SecondBis"),
			expectValLength: 2,
			validateNewVal2string: func(b block2string) error {
				if b.Name.ValueString() != "Second" ||
					b.Name2.ValueString() != "SecondBis" ||
					b.Value.ValueString() != "test2" {
					return fmt.Errorf("expected newBlock has second block data: got %#v", b)
				}

				return nil
			},
		},
		"2String - Not In List": {
			val2string: []block2string{
				{
					Name:  types.StringValue("First"),
					Name2: types.StringValue("FirstBis"),
					Value: types.StringValue("test1"),
				}, {
					Name:  types.StringValue("Second"),
					Name2: types.StringValue("SecondBis"),
					Value: types.StringValue("test2"),
				}, {
					Name:  types.StringValue("Third"),
					Name2: types.StringValue("ThirdBis"),
					Value: types.StringValue("test3"),
				},
			},
			newValString1:   types.StringValue("Fourth"),
			newValString2:   types.StringValue("SecondBis"),
			expectValLength: 3,
			validateNewVal2string: func(b block2string) error {
				if !b.Name.Equal(types.StringValue("Fourth")) ||
					!b.Name2.Equal(types.StringValue("SecondBis")) ||
					!b.Value.IsNull() {
					return fmt.Errorf("expected newBlock has empty data: got %#v", b)
				}

				return nil
			},
		},
		"2String - Not In List Bis": {
			val2string: []block2string{
				{
					Name:  types.StringValue("First"),
					Name2: types.StringValue("FirstBis"),
					Value: types.StringValue("test1"),
				}, {
					Name:  types.StringValue("Second"),
					Name2: types.StringValue("SecondBis"),
					Value: types.StringValue("test2"),
				}, {
					Name:  types.StringValue("Third"),
					Name2: types.StringValue("ThirdBis"),
					Value: types.StringValue("test3"),
				},
			},
			newValString1:   types.StringValue("Second"),
			newValString2:   types.StringValue("Third"),
			expectValLength: 3,
			validateNewVal2string: func(b block2string) error {
				if !b.Name.Equal(types.StringValue("Second")) ||
					!b.Name2.Equal(types.StringValue("Third")) ||
					!b.Value.IsNull() {
					return fmt.Errorf("expected newBlock has empty data: got %#v", b)
				}

				return nil
			},
		},
		"1Int - In List": {
			val1int: []block1int{
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
			newValInt1:      types.Int64Value(2),
			expectValLength: 2,
			validateNewVal1int: func(b block1int) error {
				if b.Name.ValueInt64() != 2 ||
					b.Value.ValueString() != "test2" {
					return fmt.Errorf("expected newBlock has second block data: got %#v", b)
				}

				return nil
			},
		},
		"1Int - Not In List": {
			val1int: []block1int{
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
			newValInt1:      types.Int64Value(4),
			expectValLength: 3,
			validateNewVal1int: func(b block1int) error {
				if !b.Name.Equal(types.Int64Value(4)) ||
					!b.Name2.IsNull() ||
					!b.Value.IsNull() {
					return fmt.Errorf("expected newBlock has empty data: got %#v", b)
				}

				return nil
			},
		},
		"Embed1String - In List": {
			valEmbed1string: []blockEmbed1string{
				{
					block1string: block1string{
						Name:  types.StringValue("First"),
						Value: types.StringValue("test1"),
					},
					Value3: types.StringValue("test1.1"),
				}, {
					block1string: block1string{
						Name:  types.StringValue("Second"),
						Value: types.StringValue("test2"),
					},
					Value3: types.StringValue("test2.1"),
				}, {
					block1string: block1string{
						Name:  types.StringValue("Third"),
						Value: types.StringValue("test3"),
					},
					Value3: types.StringValue("test3.1"),
				},
			},
			newValString1:   types.StringValue("Second"),
			expectValLength: 2,
			validateNewValEmbed1string: func(b blockEmbed1string) error {
				if b.Name.ValueString() != "Second" ||
					b.Value.ValueString() != "test2" ||
					b.Value3.ValueString() != "test2.1" {
					return fmt.Errorf("expected newBlock has second block data: got %#v", b)
				}

				return nil
			},
		},
		"Embed1String - Not In List": {
			valEmbed1string: []blockEmbed1string{
				{
					block1string: block1string{
						Name:  types.StringValue("First"),
						Value: types.StringValue("test1"),
					},
				}, {
					block1string: block1string{
						Name:  types.StringValue("Second"),
						Value: types.StringValue("test2"),
					},
				}, {
					block1string: block1string{
						Name:  types.StringValue("Third"),
						Value: types.StringValue("test3"),
					},
				},
			},
			newValString1:   types.StringValue("Fourth"),
			expectValLength: 3,
			validateNewValEmbed1string: func(b blockEmbed1string) error {
				if !b.Name.Equal(types.StringValue("Fourth")) ||
					!b.Value.IsNull() ||
					!b.Value2.IsNull() ||
					!b.Value3.IsNull() {
					return fmt.Errorf("expected newBlock has empty data: got %#v", b)
				}

				return nil
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			switch {
			case test.val1string != nil:
				var block block1string
				blocks := test.val1string
				blocks, block = tfdata.ExtractBlock(blocks,
					test.newValString1,
				)

				if v := len(blocks); v != test.expectValLength {
					t.Errorf("the expected block length is %d, got %d", test.expectValLength, v)
				}
				if err := test.validateNewVal1string(block); err != nil {
					t.Error(err)
				}
			case test.val2string != nil:
				var block block2string
				blocks := test.val2string
				blocks, block = tfdata.ExtractBlock(blocks,
					test.newValString1,
					test.newValString2,
				)

				if v := len(blocks); v != test.expectValLength {
					t.Errorf("the expected block length is %d, got %d", test.expectValLength, v)
				}
				if err := test.validateNewVal2string(block); err != nil {
					t.Error(err)
				}
			case test.val1int != nil:
				var block block1int
				blocks := test.val1int
				blocks, block = tfdata.ExtractBlock(blocks,
					test.newValInt1,
				)

				if v := len(blocks); v != test.expectValLength {
					t.Errorf("the expected block length is %d, got %d", test.expectValLength, v)
				}
				if err := test.validateNewVal1int(block); err != nil {
					t.Error(err)
				}
			case test.valEmbed1string != nil:
				var block blockEmbed1string
				blocks := test.valEmbed1string
				blocks, block = tfdata.ExtractBlock(blocks,
					test.newValString1,
				)

				if v := len(blocks); v != test.expectValLength {
					t.Errorf("the expected block length is %d, got %d", test.expectValLength, v)
				}
				if err := test.validateNewValEmbed1string(block); err != nil {
					t.Error(err)
				}
			}
		})
	}
}

func TestAppendPotentialNewBlock(t *testing.T) {
	t.Parallel()

	type block1string struct {
		Name   types.String `tfdata:"identifier"`
		Value  types.String
		Value2 types.Int64
	}

	type block2string struct {
		Name   types.String `tfdata:"identifier_1"`
		Name2  types.String `tfdata:"identifier_2"`
		Value  types.String
		Value2 types.Int64
	}

	type block1int struct {
		Name  types.Int64 `tfdata:"identifier_1"`
		Name2 types.String
		Value types.String
	}

	type blockEmbed1string struct {
		block1string
		Value3 types.String
	}

	type testCase struct {
		val1string      []block1string
		val2string      []block2string
		val1int         []block1int
		valEmbed1string []blockEmbed1string

		newValString1 types.String
		newValString2 types.String
		newValInt1    types.Int64

		expectValLength               int
		validateLatestVal1string      func(block1string) error
		validateLatestVal2string      func(block2string) error
		validateLatestVal1int         func(block1int) error
		validateLatestValEmbed1string func(blockEmbed1string) error
	}

	tests := map[string]testCase{
		"1String - Empty": {
			val1string:      []block1string{},
			newValString1:   types.StringValue("Third"),
			expectValLength: 1,
			validateLatestVal1string: func(b block1string) error {
				if b.Name.ValueString() != "Third" ||
					!b.Value.IsNull() {
					return fmt.Errorf("expected lastest block has empty data: got %#v", b)
				}

				return nil
			},
		},
		"1String - In List": {
			val1string: []block1string{
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
			newValString1:   types.StringValue("Third"),
			expectValLength: 3,
			validateLatestVal1string: func(b block1string) error {
				if b.Name.ValueString() != "Third" ||
					b.Value.ValueString() != "test3" {
					return fmt.Errorf("expected lastest block has third block data: got %#v", b)
				}

				return nil
			},
		},
		"1String - Not In List": {
			val1string: []block1string{
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
			newValString1:   types.StringValue("Fourth"),
			expectValLength: 4,
			validateLatestVal1string: func(b block1string) error {
				if !b.Name.Equal(types.StringValue("Fourth")) ||
					!b.Value.IsNull() ||
					!b.Value2.IsNull() {
					return fmt.Errorf("expected lastest block has empty data: got %#v", b)
				}

				return nil
			},
		},
		"2String - In List": {
			val2string: []block2string{
				{
					Name:  types.StringValue("First"),
					Name2: types.StringValue("FirstBis"),
					Value: types.StringValue("test1"),
				}, {
					Name:  types.StringValue("Second"),
					Name2: types.StringValue("SecondBis"),
					Value: types.StringValue("test2"),
				}, {
					Name:  types.StringValue("Third"),
					Name2: types.StringValue("ThirdBis"),
					Value: types.StringValue("test3"),
				},
			},
			newValString1:   types.StringValue("Third"),
			newValString2:   types.StringValue("ThirdBis"),
			expectValLength: 3,
			validateLatestVal2string: func(b block2string) error {
				if b.Name.ValueString() != "Third" ||
					b.Name2.ValueString() != "ThirdBis" ||
					b.Value.ValueString() != "test3" {
					return fmt.Errorf("expected lastest block has third block data: got %#v", b)
				}

				return nil
			},
		},
		"2String - Not In List": {
			val2string: []block2string{
				{
					Name:  types.StringValue("First"),
					Name2: types.StringValue("FirstBis"),
					Value: types.StringValue("test1"),
				}, {
					Name:  types.StringValue("Second"),
					Name2: types.StringValue("SecondBis"),
					Value: types.StringValue("test2"),
				}, {
					Name:  types.StringValue("Third"),
					Name2: types.StringValue("ThirdBis"),
					Value: types.StringValue("test3"),
				},
			},
			newValString1:   types.StringValue("Fourth"),
			newValString2:   types.StringValue("ThirdBis"),
			expectValLength: 4,
			validateLatestVal2string: func(b block2string) error {
				if !b.Name.Equal(types.StringValue("Fourth")) ||
					!b.Name2.Equal(types.StringValue("ThirdBis")) ||
					!b.Value.IsNull() {
					return fmt.Errorf("expected lastest block has empty data: got %#v", b)
				}

				return nil
			},
		},
		"2String - Not In List Bis": {
			val2string: []block2string{
				{
					Name:  types.StringValue("First"),
					Name2: types.StringValue("FirstBis"),
					Value: types.StringValue("test1"),
				}, {
					Name:  types.StringValue("Second"),
					Name2: types.StringValue("SecondBis"),
					Value: types.StringValue("test2"),
				}, {
					Name:  types.StringValue("Third"),
					Name2: types.StringValue("ThirdBis"),
					Value: types.StringValue("test3"),
				},
			},
			newValString1:   types.StringValue("Third"),
			newValString2:   types.StringValue("Third"),
			expectValLength: 4,
			validateLatestVal2string: func(b block2string) error {
				if !b.Name.Equal(types.StringValue("Third")) ||
					!b.Name2.Equal(types.StringValue("Third")) ||
					!b.Value.IsNull() {
					return fmt.Errorf("expected lastest block has empty data: got %#v", b)
				}

				return nil
			},
		},
		"1Int - In List": {
			val1int: []block1int{
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
			newValInt1:      types.Int64Value(3),
			expectValLength: 3,
			validateLatestVal1int: func(b block1int) error {
				if b.Name.ValueInt64() != 3 ||
					b.Value.ValueString() != "test3" {
					return fmt.Errorf("expected lastest block has third block data: got %#v", b)
				}

				return nil
			},
		},
		"1Int - Not In List": {
			val1int: []block1int{
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
			newValInt1:      types.Int64Value(4),
			expectValLength: 4,
			validateLatestVal1int: func(b block1int) error {
				if !b.Name.Equal(types.Int64Value(4)) ||
					!b.Name2.IsNull() ||
					!b.Value.IsNull() {
					return fmt.Errorf("expected lastest block has empty data: got %#v", b)
				}

				return nil
			},
		},
		"Embed1String - In List": {
			valEmbed1string: []blockEmbed1string{
				{
					block1string: block1string{
						Name:  types.StringValue("First"),
						Value: types.StringValue("test1"),
					},
					Value3: types.StringValue("test1.1"),
				}, {
					block1string: block1string{
						Name:  types.StringValue("Second"),
						Value: types.StringValue("test2"),
					},
					Value3: types.StringValue("test2.1"),
				}, {
					block1string: block1string{
						Name:  types.StringValue("Third"),
						Value: types.StringValue("test3"),
					},
					Value3: types.StringValue("test3.1"),
				},
			},
			newValString1:   types.StringValue("Third"),
			expectValLength: 3,
			validateLatestValEmbed1string: func(b blockEmbed1string) error {
				if b.Name.ValueString() != "Third" ||
					b.Value.ValueString() != "test3" ||
					b.Value3.ValueString() != "test3.1" {
					return fmt.Errorf("expected lastest block has third block data: got %#v", b)
				}

				return nil
			},
		},
		"Embed1String - Not In List": {
			valEmbed1string: []blockEmbed1string{
				{
					block1string: block1string{
						Name:  types.StringValue("First"),
						Value: types.StringValue("test1"),
					},
				}, {
					block1string: block1string{
						Name:  types.StringValue("Second"),
						Value: types.StringValue("test2"),
					},
				}, {
					block1string: block1string{
						Name:  types.StringValue("Third"),
						Value: types.StringValue("test3"),
					},
				},
			},
			newValString1:   types.StringValue("Fourth"),
			expectValLength: 4,
			validateLatestValEmbed1string: func(b blockEmbed1string) error {
				if !b.Name.Equal(types.StringValue("Fourth")) ||
					!b.Value.IsNull() ||
					!b.Value2.IsNull() ||
					!b.Value3.IsNull() {
					return fmt.Errorf("expected lastest block has empty data: got %#v", b)
				}

				return nil
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			switch {
			case test.val1string != nil:
				blocks := test.val1string
				blocks = tfdata.AppendPotentialNewBlock(blocks,
					test.newValString1,
				)

				if v := len(blocks); v != test.expectValLength {
					t.Errorf("the expected block length is %d, got %d", test.expectValLength, v)
				}
				if err := test.validateLatestVal1string(blocks[len(blocks)-1]); err != nil {
					t.Error(err)
				}
			case test.val2string != nil:
				blocks := test.val2string
				blocks = tfdata.AppendPotentialNewBlock(blocks,
					test.newValString1,
					test.newValString2,
				)

				if v := len(blocks); v != test.expectValLength {
					t.Errorf("the expected block length is %d, got %d", test.expectValLength, v)
				}
				if err := test.validateLatestVal2string(blocks[len(blocks)-1]); err != nil {
					t.Error(err)
				}
			case test.val1int != nil:
				blocks := test.val1int
				blocks = tfdata.AppendPotentialNewBlock(blocks,
					test.newValInt1,
				)

				if v := len(blocks); v != test.expectValLength {
					t.Errorf("the expected block length is %d, got %d", test.expectValLength, v)
				}
				if err := test.validateLatestVal1int(blocks[len(blocks)-1]); err != nil {
					t.Error(err)
				}
			case test.valEmbed1string != nil:
				blocks := test.valEmbed1string
				blocks = tfdata.AppendPotentialNewBlock(blocks,
					test.newValString1,
				)

				if v := len(blocks); v != test.expectValLength {
					t.Errorf("the expected block length is %d, got %d", test.expectValLength, v)
				}
				if err := test.validateLatestValEmbed1string(blocks[len(blocks)-1]); err != nil {
					t.Error(err)
				}
			}
		})
	}
}
