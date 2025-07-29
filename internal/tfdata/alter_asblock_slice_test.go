package tfdata_test

import (
	"reflect"
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

		expectValLength         int
		expectBlock1string      *block1string
		expectBlock2string      *block2string
		expectBlock1int         *block1int
		expectBlockEmbed1string *blockEmbed1string
	}

	tests := map[string]testCase{
		"1String - Empty": {
			val1string:      []block1string{},
			newValString1:   types.StringValue("Third"),
			expectValLength: 0,
			expectBlock1string: &block1string{
				Name: types.StringValue("Third"),
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
			expectBlock1string: &block1string{
				Name:  types.StringValue("Second"),
				Value: types.StringValue("test2"),
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
			expectBlock1string: &block1string{
				Name: types.StringValue("Fourth"),
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
			expectBlock2string: &block2string{
				Name:  types.StringValue("Second"),
				Name2: types.StringValue("SecondBis"),
				Value: types.StringValue("test2"),
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
			expectBlock2string: &block2string{
				Name:  types.StringValue("Fourth"),
				Name2: types.StringValue("SecondBis"),
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
			expectBlock2string: &block2string{
				Name:  types.StringValue("Second"),
				Name2: types.StringValue("Third"),
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
			expectBlock1int: &block1int{
				Name:  types.Int64Value(2),
				Value: types.StringValue("test2"),
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
			expectBlock1int: &block1int{
				Name: types.Int64Value(4),
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
			expectBlockEmbed1string: &blockEmbed1string{
				block1string: block1string{
					Name:  types.StringValue("Second"),
					Value: types.StringValue("test2"),
				},
				Value3: types.StringValue("test2.1"),
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
			expectBlockEmbed1string: &blockEmbed1string{
				block1string: block1string{
					Name: types.StringValue("Fourth"),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			switch {
			case test.expectBlock1string != nil:
				var block block1string
				blocks := test.val1string
				blocks, block = tfdata.ExtractBlock(blocks,
					test.newValString1,
				)

				if v := len(blocks); v != test.expectValLength {
					t.Errorf("the expected block length is %d, got %d", test.expectValLength, v)
				}
				if !reflect.DeepEqual(*test.expectBlock1string, block) {
					t.Errorf("unexpected block: want %#v, got %#v", *test.expectBlock1string, block)
				}
			case test.expectBlock2string != nil:
				var block block2string
				blocks := test.val2string
				blocks, block = tfdata.ExtractBlock(blocks,
					test.newValString1,
					test.newValString2,
				)

				if v := len(blocks); v != test.expectValLength {
					t.Errorf("the expected block length is %d, got %d", test.expectValLength, v)
				}
				if !reflect.DeepEqual(*test.expectBlock2string, block) {
					t.Errorf("unexpected block: want %#v, got %#v", *test.expectBlock2string, block)
				}
			case test.expectBlock1int != nil:
				var block block1int
				blocks := test.val1int
				blocks, block = tfdata.ExtractBlock(blocks,
					test.newValInt1,
				)

				if v := len(blocks); v != test.expectValLength {
					t.Errorf("the expected block length is %d, got %d", test.expectValLength, v)
				}
				if !reflect.DeepEqual(*test.expectBlock1int, block) {
					t.Errorf("unexpected block: want %#v, got %#v", *test.expectBlock1int, block)
				}
			case test.expectBlockEmbed1string != nil:
				var block blockEmbed1string
				blocks := test.valEmbed1string
				blocks, block = tfdata.ExtractBlock(blocks,
					test.newValString1,
				)

				if v := len(blocks); v != test.expectValLength {
					t.Errorf("the expected block length is %d, got %d", test.expectValLength, v)
				}
				if !reflect.DeepEqual(*test.expectBlockEmbed1string, block) {
					t.Errorf("unexpected block: want %#v, got %#v", *test.expectBlockEmbed1string, block)
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
		expectLatestBlock1string      *block1string
		expectLatestBlock2string      *block2string
		expectLatestBlock1int         *block1int
		expectLatestBlockEmbed1string *blockEmbed1string
	}

	tests := map[string]testCase{
		"1String - Empty": {
			val1string:      []block1string{},
			newValString1:   types.StringValue("Third"),
			expectValLength: 1,
			expectLatestBlock1string: &block1string{
				Name: types.StringValue("Third"),
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
			expectLatestBlock1string: &block1string{
				Name:  types.StringValue("Third"),
				Value: types.StringValue("test3"),
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
			expectLatestBlock1string: &block1string{
				Name: types.StringValue("Fourth"),
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
			expectLatestBlock2string: &block2string{
				Name:  types.StringValue("Third"),
				Name2: types.StringValue("ThirdBis"),
				Value: types.StringValue("test3"),
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
			expectLatestBlock2string: &block2string{
				Name:  types.StringValue("Fourth"),
				Name2: types.StringValue("ThirdBis"),
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
			expectLatestBlock2string: &block2string{
				Name:  types.StringValue("Third"),
				Name2: types.StringValue("Third"),
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
			expectLatestBlock1int: &block1int{
				Name:  types.Int64Value(3),
				Value: types.StringValue("test3"),
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
			expectLatestBlock1int: &block1int{
				Name: types.Int64Value(4),
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
			expectLatestBlockEmbed1string: &blockEmbed1string{
				block1string: block1string{
					Name:  types.StringValue("Third"),
					Value: types.StringValue("test3"),
				},
				Value3: types.StringValue("test3.1"),
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
			expectLatestBlockEmbed1string: &blockEmbed1string{
				block1string: block1string{
					Name: types.StringValue("Fourth"),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			switch {
			case test.expectLatestBlock1string != nil:
				blocks := test.val1string
				blocks = tfdata.AppendPotentialNewBlock(blocks,
					test.newValString1,
				)

				if v := len(blocks); v != test.expectValLength {
					t.Errorf("the expected block length is %d, got %d", test.expectValLength, v)
				} else if !reflect.DeepEqual(*test.expectLatestBlock1string, blocks[len(blocks)-1]) {
					t.Errorf("unexpected latest block: want %#v, got %#v", *test.expectLatestBlock1string, blocks[len(blocks)-1])
				}
			case test.expectLatestBlock2string != nil:
				blocks := test.val2string
				blocks = tfdata.AppendPotentialNewBlock(blocks,
					test.newValString1,
					test.newValString2,
				)

				if v := len(blocks); v != test.expectValLength {
					t.Errorf("the expected block length is %d, got %d", test.expectValLength, v)
				} else if !reflect.DeepEqual(*test.expectLatestBlock2string, blocks[len(blocks)-1]) {
					t.Errorf("unexpected latest block: want %#v, got %#v", *test.expectLatestBlock2string, blocks[len(blocks)-1])
				}
			case test.expectLatestBlock1int != nil:
				blocks := test.val1int
				blocks = tfdata.AppendPotentialNewBlock(blocks,
					test.newValInt1,
				)

				if v := len(blocks); v != test.expectValLength {
					t.Errorf("the expected block length is %d, got %d", test.expectValLength, v)
				} else if !reflect.DeepEqual(*test.expectLatestBlock1int, blocks[len(blocks)-1]) {
					t.Errorf("unexpected latest block: want %#v, got %#v", *test.expectLatestBlock1int, blocks[len(blocks)-1])
				}
			case test.expectLatestBlockEmbed1string != nil:
				blocks := test.valEmbed1string
				blocks = tfdata.AppendPotentialNewBlock(blocks,
					test.newValString1,
				)

				if v := len(blocks); v != test.expectValLength {
					t.Errorf("the expected block length is %d, got %d", test.expectValLength, v)
				} else if !reflect.DeepEqual(*test.expectLatestBlockEmbed1string, blocks[len(blocks)-1]) {
					t.Errorf("unexpected latest block: want %#v, got %#v", *test.expectLatestBlockEmbed1string, blocks[len(blocks)-1])
				}
			}
		})
	}
}
