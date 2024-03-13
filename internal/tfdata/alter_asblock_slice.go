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
func ExtractBlockWithTFTypesString[B any](
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
// search field with name = 'structField1Name' and with type 'types.String'
//   - if the value of this field is equal with input1Value,
//     search again field with name = 'structField2Name' and with type 'types.String'
//   - if the value of second field is equal with input2Value,
//     remove element from slice and return the new slice and the element
//   - if structField1Name and structField2Name not equal with respective value,
//     create a new empty struct and return the slice unaltered and the new struct.
func ExtractBlockWith2TFTypesString[B any]( //nolint:ireturn
	blocks []B, structField1Name, input1Value, structField2Name, input2Value string,
) (
	[]B, B,
) {
	for i, block := range blocks {
		field1Value := reflect.ValueOf(block).FieldByNameFunc(func(name string) bool {
			return strings.EqualFold(structField1Name, name)
		})
		if !field1Value.IsValid() {
			continue
		}
		if tfString, ok := field1Value.Interface().(types.String); ok {
			if tfString.ValueString() != input1Value {
				continue
			}
		} else {
			continue
		}
		field2Value := reflect.ValueOf(block).FieldByNameFunc(func(name string) bool {
			return strings.EqualFold(structField2Name, name)
		})
		if !field2Value.IsValid() {
			continue
		}
		if tfString, ok := field2Value.Interface().(types.String); ok {
			if tfString.ValueString() == input2Value {
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
func ExtractBlockWithTFTypesInt64[B any](
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
