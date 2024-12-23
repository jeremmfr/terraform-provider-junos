package tfdata

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
)

const skipIsEmpty = "skip_isempty"

// check if block struct doesn't have either:
//   - an framework attribute with not null value
//   - a slice with at least an element
//   - a not nil pointer
//
// fields with tag tfdata:"skip_isempty" are excluded.
func CheckBlockIsEmpty[B any](block B) bool {
	blockValue := reflect.Indirect(reflect.ValueOf(block).Elem())

	return checkBlockValueIsEmpty(blockValue)
}

func checkBlockValueIsEmpty(blockValue reflect.Value) bool {
	for iField := range blockValue.NumField() {
		if slices.ContainsFunc(tagsOfStructField(blockValue.Type().Field(iField)), func(s string) bool {
			return strings.EqualFold(s, skipIsEmpty)
		}) {
			continue
		}

		fieldValue := blockValue.Field(iField)
		if !fieldValue.IsValid() {
			continue
		}

		if blockValue.Type().Field(iField).Anonymous {
			embedBlockValue := reflect.Indirect(reflect.NewAt(fieldValue.Type(), fieldValue.Addr().UnsafePointer()).Elem())

			if !checkBlockValueIsEmpty(embedBlockValue) {
				return false
			}

			continue
		}

		if attrValue, ok := fieldValue.Interface().(attr.Value); ok {
			if !attrValue.IsNull() {
				return false
			}

			continue
		}
		if fieldValue.Type().Kind() == reflect.Slice {
			if fieldValue.Len() != 0 {
				return false
			}

			continue
		}
		if fieldValue.Type().Kind() == reflect.Pointer {
			if !fieldValue.IsNil() {
				return false
			}

			continue
		}

		panic(fmt.Sprintf(
			"don't know how to determine if field %q (type: %s) is empty",
			blockValue.Type().Field(iField).Name, fieldValue.Type().Name(),
		))
	}

	return true
}

// check if struct has either :
//   - an framework attribute with known value (not null and not unknown)
//   - an pointer to an other struct with a known framework attribute value.
func CheckBlockHasKnownValue[B any](block B, excludeFields ...string) bool {
	blockValue := reflect.Indirect(reflect.ValueOf(block).Elem())

	return checkBlockValueHasKnownValue(blockValue, excludeFields...)
}

func checkBlockValueHasKnownValue(blockValue reflect.Value, excludeFields ...string) bool {
	for iField := range blockValue.NumField() {
		if slices.Contains(excludeFields, blockValue.Type().Field(iField).Name) {
			continue
		}

		fieldValue := blockValue.Field(iField)
		if !fieldValue.IsValid() {
			continue
		}

		if blockValue.Type().Field(iField).Anonymous {
			embedBlockValue := reflect.Indirect(reflect.NewAt(fieldValue.Type(), fieldValue.Addr().UnsafePointer()).Elem())

			if checkBlockValueHasKnownValue(embedBlockValue, excludeFields...) {
				return true
			}

			continue
		}

		if attrValue, ok := fieldValue.Interface().(attr.Value); ok {
			if !attrValue.IsNull() && !attrValue.IsUnknown() {
				return true
			}

			continue
		}

		if fieldValue.Type().Kind() == reflect.Pointer {
			if !fieldValue.IsNil() {
				if CheckBlockHasKnownValue(blockValue.Field(iField).Interface()) {
					return true
				}
			}

			continue
		}

		panic(fmt.Sprintf(
			"don't know how to determine if field %q (type: %s) is known",
			blockValue.Type().Field(iField).Name, fieldValue.Type().Name(),
		))
	}

	return false
}
