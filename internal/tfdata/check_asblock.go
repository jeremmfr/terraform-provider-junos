package tfdata

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/attr"
)

// check if block struct doesn't have either:
//   - an framework attribute with not null value
//   - a slice with at least an alement
//   - a not nil pointer.
func CheckBlockIsEmpty[B any](block B, excludeFields ...string) bool {
	v := reflect.Indirect(reflect.ValueOf(block).Elem())

	for i := 0; i < v.NumField(); i++ {
		if slices.Contains(excludeFields, v.Type().Field(i).Name) {
			continue
		}

		fieldValue := v.Field(i)
		if !fieldValue.IsValid() {
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
			v.Type().Field(i).Name, fieldValue.Type().Name(),
		))
	}

	return true
}

// check if struct has either :
//   - an framework attribute with known value (not null and not unknown)
//   - an pointer to an other struct with a known framework attribute value.
func CheckBlockHasKnownValue[B any](block B, excludeFields ...string) bool {
	v := reflect.Indirect(reflect.ValueOf(block).Elem())

	for i := 0; i < v.NumField(); i++ {
		if slices.Contains(excludeFields, v.Type().Field(i).Name) {
			continue
		}

		fieldValue := v.Field(i)
		if !fieldValue.IsValid() {
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
				if CheckBlockHasKnownValue(v.Field(i).Interface()) {
					return true
				}
			}

			continue
		}

		panic(fmt.Sprintf(
			"don't know how to determine if field %q (type: %s) is known",
			v.Type().Field(i).Name, fieldValue.Type().Name(),
		))
	}

	return false
}
