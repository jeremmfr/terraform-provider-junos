package tfdata

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	identifier  = "identifier"
	identifier1 = "identifier_1"
	identifier2 = "identifier_2"
)

type reflectStructField struct {
	value       reflect.Value
	structField reflect.StructField
}

// for each struct in blocks (list of struct)
//
//   - if the values of fields with identifier tag are equal to the elements of identifierValues
//     (first element of identifierValues equal to field with tag value identifier or identifier_1,
//     optional second element of identifierValues equal to field with tag value identifier_2)
//
//     -> remove element from slice and return the new slice and the element
//
//   - if no one match identifierValues
//
//     -> create a new empty struct, assign identifierValues to fields with identifier tag
//     and return the slice unaltered and the new struct.
func ExtractBlock[B any](
	blocks []B, identifierValues ...attr.Value,
) (
	[]B, B,
) {
	// read blocks to search if one block match with arguments
loopBlocks:
	// FIX ME with go 1.23.O in go.mod
	// for iBlock, block := range slices.Backward(blocks) {
	for iBlock := len(blocks) - 1; iBlock >= 0; iBlock-- {
		block := blocks[iBlock]
		blockValue := reflect.ValueOf(block)
		if !blockValue.IsValid() {
			continue
		}

		blockIdentifierValues := make([]reflect.Value, len(identifierValues))
		blockIdentifierFound := make(map[string]struct{})

		// read fields to pick up values from fields with identifier tags
		for iField := range blockValue.NumField() {
			fieldValue := blockValue.Field(iField)
			if !fieldValue.IsValid() {
				continue
			}

			if blockValue.Type().Field(iField).Anonymous {
				for iEmbedField := range fieldValue.NumField() {
					embedFieldValue := fieldValue.Field(iEmbedField)
					if !embedFieldValue.IsValid() {
						continue
					}

					reflectStructField{
						value:       embedFieldValue,
						structField: fieldValue.Type().Field(iEmbedField),
					}.searchIdentifierOnStructField(
						blockIdentifierFound,
						blockIdentifierValues,
					)
				}

				continue
			}

			reflectStructField{
				value:       fieldValue,
				structField: blockValue.Type().Field(iField),
			}.searchIdentifierOnStructField(
				blockIdentifierFound,
				blockIdentifierValues,
			)
		}

		// detect mismatch between tags and arguments
		if len(blockIdentifierFound) == 0 {
			panic("no tfdata:identifier tag found on struct")
		}
		if len(blockIdentifierFound) != len(identifierValues) {
			panic(fmt.Sprintf(
				"mismatch number of tfdata:identifier tags and identifierValues: found %d tags, got %d identifierValues",
				len(blockIdentifierFound), len(identifierValues),
			))
		}

		// check if match with the arguments
		for iIdent, fieldValue := range blockIdentifierValues {
			switch attrValue := fieldValue.Interface().(type) {
			case types.String:
				if !attrValue.Equal(identifierValues[iIdent].(types.String)) {
					continue loopBlocks
				}
			case types.Int64:
				if !attrValue.Equal(identifierValues[iIdent].(types.Int64)) {
					continue loopBlocks
				}
			default:
				panic(fmt.Sprintf(
					"don't know how to compare field with type %q",
					fieldValue.Type().Name(),
				))
			}
		}

		// match so remove block from slice and return result
		blocks = append(blocks[:iBlock], blocks[iBlock+1:]...)

		return blocks, block
	}

	// no blocks match with arguments so generate new block
	newBlock := new(B)
	newBlockValue := reflect.ValueOf(newBlock).Elem()
	for iField := range newBlockValue.NumField() {
		if newBlockValue.Type().Field(iField).Anonymous {
			embedNewBlockValue := newBlockValue.Field(iField)

			for iEmbedField := range embedNewBlockValue.NumField() {
				reflectStructField{
					value:       embedNewBlockValue.Field(iEmbedField),
					structField: embedNewBlockValue.Type().Field(iEmbedField),
				}.assignValueToIdentifier(identifierValues)
			}

			continue
		}

		reflectStructField{
			value:       newBlockValue.Field(iField),
			structField: newBlockValue.Type().Field(iField),
		}.assignValueToIdentifier(identifierValues)
	}

	return blocks, *newBlock
}

func (field reflectStructField) searchIdentifierOnStructField(
	blockIdentifierFound map[string]struct{},
	blockIdentifierValues []reflect.Value,
) {
	fieldTags := tagsOfStructField(field.structField)

	switch {
	case slices.ContainsFunc(fieldTags, func(s string) bool {
		return strings.EqualFold(s, identifier)
	}),
		slices.ContainsFunc(fieldTags, func(s string) bool {
			return strings.EqualFold(s, identifier1)
		}):
		if _, ok := blockIdentifierFound[identifier1]; ok {
			panic("multiple tfdata " + identifier + " or " + identifier1 + " tags on struct")
		}

		blockIdentifierFound[identifier1] = struct{}{}
		blockIdentifierValues[0] = field.value
	case slices.ContainsFunc(fieldTags, func(s string) bool {
		return strings.EqualFold(s, identifier2)
	}):
		if _, ok := blockIdentifierFound[identifier2]; ok {
			panic("multiple tfdata " + identifier2 + " tags on struct")
		}

		blockIdentifierFound[identifier2] = struct{}{}
		blockIdentifierValues[1] = field.value
	}
}

func (field reflectStructField) assignValueToIdentifier(
	identifierValues []attr.Value,
) {
	fieldTags := tagsOfStructField(field.structField)

	switch {
	case slices.ContainsFunc(fieldTags, func(s string) bool {
		return strings.EqualFold(s, identifier)
	}),
		slices.ContainsFunc(fieldTags, func(s string) bool {
			return strings.EqualFold(s, identifier1)
		}):
		field.value.Set(reflect.ValueOf(identifierValues[0]))
	case slices.ContainsFunc(fieldTags, func(s string) bool {
		return strings.EqualFold(s, identifier2)
	}):
		field.value.Set(reflect.ValueOf(identifierValues[1]))
	}
}
