package tfdata

import (
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// for each struct in blocks (list of struct)
// search field with name = 'structFieldName' and with type 'types.String'
//   - if the value of this field is equal with inputValue,
//     remove element from slice and return the new slice and the element
//   - if not equal, create a new empty struct and return the slice unaltered and the new struct.
func ExtractBlockWithTFTypesString[B any]( //nolint:ireturn
	blocks []B, structFieldName, inputValue string,
) (
	[]B, B,
) {
	for i, block := range blocks {
		fieldValue := reflect.ValueOf(block).FieldByNameFunc(func(name string) bool {
			return strings.EqualFold(structFieldName, name)
		})
		if !fieldValue.IsValid() {
			continue
		}
		if tfString, ok := fieldValue.Interface().(types.String); ok {
			if tfString.ValueString() == inputValue {
				blocks = append(blocks[:i], blocks[i+1:]...)

				return blocks, block
			}
		}
	}
	e := new(B)

	return blocks, *e
}

// for each struct in blocks (list of struct)
// search field with name = 'structFieldName' and with type 'types.Int64'
//   - if the value of this field is equal with inputValue,
//     remove element from slice and return the new slice and the element
//   - if not equal, create a new empty struct and return the slice unaltered and the new struct.
func ExtractBlockWithTFTypesInt64[B any]( //nolint:ireturn
	blocks []B, structFieldName string, inputValue int64,
) (
	[]B, B,
) {
	for i, block := range blocks {
		fieldValue := reflect.ValueOf(block).FieldByNameFunc(func(name string) bool {
			return strings.EqualFold(structFieldName, name)
		})
		if !fieldValue.IsValid() {
			continue
		}
		if tfInt64, ok := fieldValue.Interface().(types.Int64); ok {
			if tfInt64.ValueInt64() == inputValue {
				blocks = append(blocks[:i], blocks[i+1:]...)

				return blocks, block
			}
		}
	}
	e := new(B)

	return blocks, *e
}
