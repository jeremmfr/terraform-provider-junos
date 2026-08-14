package tftypes

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ basetypes.StringTypable                    = StringDateType{}
	_ basetypes.StringValuableWithSemanticEquals = StringDate{}
)

// StringDateType is the attribute type of a StringDate value.
type StringDateType struct {
	basetypes.StringType
}

func (t StringDateType) String() string {
	return "tftypes.StringDateType"
}

func (t StringDateType) Equal(o attr.Type) bool {
	other, ok := o.(StringDateType)
	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t StringDateType) ValueType(_ context.Context) attr.Value {
	return StringDate{}
}

func (t StringDateType) ValueFromString(
	_ context.Context, in basetypes.StringValue,
) (basetypes.StringValuable, diag.Diagnostics) {
	return StringDate{StringValue: in}, nil
}

func (t StringDateType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	stringValue, ok := attrValue.(basetypes.StringValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type %T, expected basetypes.StringValue", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}

	return stringValuable, nil
}

// StringDate is a date with the format YYYY-MM-DD.HH:MM:SS, where the leading zero
// of the month and the day is not significant.
//
// Junos removes these leading zeros when storing the value, so the value read
// is always in the short form. Considering both forms equal avoids a permanent
// diff and keeps in state the form present in the configuration.
type StringDate struct {
	basetypes.StringValue
}

func (v StringDate) Type(_ context.Context) attr.Type {
	return StringDateType{}
}

func (v StringDate) Equal(o attr.Value) bool {
	other, ok := o.(StringDate)
	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

// StringSemanticEquals returns true when the two dates only differ by the
// leading zero of the month or the day.
func (v StringDate) StringSemanticEquals(
	_ context.Context, newValuable basetypes.StringValuable,
) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(StringDate)
	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected Value Type: "+fmt.Sprintf("%T", v)+"\n"+
				"Got Value Type: "+fmt.Sprintf("%T", newValuable),
		)

		return false, diags
	}

	return trimDateLeadingZeros(v.ValueString()) == trimDateLeadingZeros(newValue.ValueString()), diags
}

// NewStringDateNull creates a StringDate with a null value.
func NewStringDateNull() StringDate {
	return StringDate{StringValue: basetypes.NewStringNull()}
}

// NewStringDateUnknown creates a StringDate with an unknown value.
func NewStringDateUnknown() StringDate {
	return StringDate{StringValue: basetypes.NewStringUnknown()}
}

// NewStringDateValue creates a StringDate with a known value.
func NewStringDateValue(value string) StringDate {
	return StringDate{StringValue: basetypes.NewStringValue(value)}
}

// NewStringDateFromString creates a StringDate with the value of a string, keeping
// the null and unknown states.
func NewStringDateFromString(value basetypes.StringValue) StringDate {
	return StringDate{StringValue: value}
}

var dateLeadingZeros = regexp.MustCompile(`^(\d{4})-0?(\d{1,2})-0?(\d{1,2})\.(.+)$`)

// trimDateLeadingZeros removes the leading zero of the month and the day to
// compare two dates. The result is only used for comparison, never stored.
func trimDateLeadingZeros(value string) string {
	matches := dateLeadingZeros.FindStringSubmatch(value)
	if matches == nil {
		return value
	}

	return matches[1] + "-" + matches[2] + "-" + matches[3] + "." + matches[4]
}
