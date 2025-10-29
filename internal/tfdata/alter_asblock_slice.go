package tfdata

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
)

const (
	identifier  = "identifier"
	identifier1 = "identifier_1"
	identifier2 = "identifier_2"
)

// ExtractBlock: for each struct in blocks (list of struct)
//
//   - if the values of fields with identifier tag are equal to the elements of identifierValues
//     (first element of identifierValues equal to field with tag value identifier or identifier_1,
//     optional second element of identifierValues equal to field with tag value identifier_2)
//
//     -> remove element from slice and return the new slice and the element
//
//   - if no one match identifierValues
//
//     -> create a new empty struct, assign identifierValues to fields with identifier tag,
//     and return the slice unaltered and the new struct.
func ExtractBlock[B any](
	blocks []B, identifierValues ...attr.Value,
) (
	[]B, B,
) {
	// read blocks to search if one block match with arguments
loopBlocks:
	for iBlock, block := range slices.Backward(blocks) {
		blockValue := reflect.ValueOf(block)
		if !blockValue.IsValid() {
			continue
		}

		blockIdentifierTagsFound := make(map[string]struct{})
		blockIdentifierValues := make([]reflect.Value, len(identifierValues))

		readBlockIdentifierValues(blockValue, blockIdentifierTagsFound, blockIdentifierValues)

		// detect mismatch between tags and arguments
		if len(blockIdentifierTagsFound) == 0 {
			panic("no tfdata identifier tag found on struct")
		}
		if len(blockIdentifierTagsFound) != len(identifierValues) {
			panic(fmt.Sprintf(
				"mismatch between number of tfdata identifier tags and identifierValues:"+
					" found %d tags, got %d identifierValues",
				len(blockIdentifierTagsFound), len(identifierValues),
			))
		}

		// check if match with the arguments
		for iIdent, fieldValue := range blockIdentifierValues {
			attrValue := fieldValue.Interface().(attr.Value)
			if !attrValue.Equal(identifierValues[iIdent]) {
				continue loopBlocks
			}
		}

		// match so remove block from slice and return result
		blocks = append(blocks[:iBlock], blocks[iBlock+1:]...)

		return blocks, block
	}

	// no blocks match with arguments so generate new block
	newBlock := new(B)
	assignIdentifierValuesToBlock(identifierValues, reflect.ValueOf(newBlock).Elem())

	return blocks, *newBlock
}

// AppendPotentialNewBlock: with latest block of slice
//
//   - if the values of fields with identifier tag are equal to the elements of identifierValues
//     (first element of identifierValues equal to field with tag value identifier or identifier_1,
//     optional second element of identifierValues equal to field with tag value identifier_2)
//
//     ->  return the slice unaltered
//
//   - if not match identifierValues
//
//     -> create a new empty struct, assign identifierValues to fields with identifier tag,
//     append it to slice, and return the new slice.
func AppendPotentialNewBlock[B any](
	blocks []B, identifierValues ...attr.Value,
) []B {
	if len(blocks) == 0 {
		newBlock := new(B)
		assignIdentifierValuesToBlock(identifierValues, reflect.ValueOf(newBlock).Elem())

		return []B{*newBlock}
	}

	latestBlockValue := reflect.ValueOf(blocks[len(blocks)-1])
	lastestOK := true
	if !latestBlockValue.IsValid() {
		lastestOK = false
	} else {
		blockIdentifierTagsFound := make(map[string]struct{})
		blockIdentifierValues := make([]reflect.Value, len(identifierValues))

		readBlockIdentifierValues(latestBlockValue, blockIdentifierTagsFound, blockIdentifierValues)

		// detect mismatch between tags and arguments
		if len(blockIdentifierTagsFound) == 0 {
			panic("no tfdata identifier tag found on struct")
		}
		if len(blockIdentifierTagsFound) != len(identifierValues) {
			panic(fmt.Sprintf(
				"mismatch between number of tfdata identifier tags and identifierValues:"+
					" found %d tags, got %d identifierValues",
				len(blockIdentifierTagsFound), len(identifierValues),
			))
		}

		// check if match with the arguments
		for iIdent, fieldValue := range blockIdentifierValues {
			attrValue := fieldValue.Interface().(attr.Value)
			if !attrValue.Equal(identifierValues[iIdent]) {
				lastestOK = false

				break
			}
		}
	}

	if lastestOK {
		return blocks
	}

	newBlock := new(B)
	assignIdentifierValuesToBlock(identifierValues, reflect.ValueOf(newBlock).Elem())

	return append(blocks, *newBlock)
}

func readBlockIdentifierValues(
	blockValue reflect.Value,
	blockIdentifierTagsFound map[string]struct{},
	blockIdentifierValues []reflect.Value,
) {
	for iField := range blockValue.NumField() {
		fieldValue := blockValue.Field(iField)
		if !fieldValue.IsValid() {
			continue
		}

		if blockValue.Type().Field(iField).Anonymous {
			readBlockIdentifierValues(fieldValue, blockIdentifierTagsFound, blockIdentifierValues)

			continue
		}

		fieldTags := tagsOfStructField(blockValue.Type().Field(iField))

		switch {
		case slices.ContainsFunc(fieldTags, func(s string) bool {
			return strings.EqualFold(s, identifier)
		}),
			slices.ContainsFunc(fieldTags, func(s string) bool {
				return strings.EqualFold(s, identifier1)
			}):
			if _, ok := blockIdentifierTagsFound[identifier1]; ok {
				panic("multiple tfdata " + identifier + " or " + identifier1 + " tags on struct")
			}

			blockIdentifierTagsFound[identifier1] = struct{}{}
			blockIdentifierValues[0] = fieldValue

		case slices.ContainsFunc(fieldTags, func(s string) bool {
			return strings.EqualFold(s, identifier2)
		}):
			if _, ok := blockIdentifierTagsFound[identifier2]; ok {
				panic("multiple tfdata " + identifier2 + " tags on struct")
			}

			blockIdentifierTagsFound[identifier2] = struct{}{}
			blockIdentifierValues[1] = fieldValue
		}
	}
}

func assignIdentifierValuesToBlock(
	identifierValues []attr.Value,
	blockValue reflect.Value,
) {
	for iField := range blockValue.NumField() {
		if blockValue.Type().Field(iField).Anonymous {
			assignIdentifierValuesToBlock(identifierValues, blockValue.Field(iField))

			continue
		}

		fieldTags := tagsOfStructField(blockValue.Type().Field(iField))

		switch {
		case slices.ContainsFunc(fieldTags, func(s string) bool {
			return strings.EqualFold(s, identifier)
		}),
			slices.ContainsFunc(fieldTags, func(s string) bool {
				return strings.EqualFold(s, identifier1)
			}):
			blockValue.Field(iField).Set(
				reflect.ValueOf(identifierValues[0]),
			)

		case slices.ContainsFunc(fieldTags, func(s string) bool {
			return strings.EqualFold(s, identifier2)
		}):
			blockValue.Field(iField).Set(
				reflect.ValueOf(identifierValues[1]),
			)
		}
	}
}
