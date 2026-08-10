package tftypes_test

import (
	"context"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestDateStringSemanticEquals(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		priorValue    tftypes.StringDate
		newValue      basetypes.StringValuable
		expectedEqual bool
		expectedError bool
	}{
		"leading zero on month and day": {
			priorValue:    tftypes.NewStringDateValue("2022-06-01.00:00:00"),
			newValue:      tftypes.NewStringDateValue("2022-6-1.00:00:00"),
			expectedEqual: true,
		},
		"leading zero on month only": {
			priorValue:    tftypes.NewStringDateValue("2022-06-11.10:09:08"),
			newValue:      tftypes.NewStringDateValue("2022-6-11.10:09:08"),
			expectedEqual: true,
		},
		"leading zero on day only": {
			priorValue:    tftypes.NewStringDateValue("2022-12-01.10:09:08"),
			newValue:      tftypes.NewStringDateValue("2022-12-1.10:09:08"),
			expectedEqual: true,
		},
		"short form in configuration": {
			priorValue:    tftypes.NewStringDateValue("2016-1-1.02:00:00"),
			newValue:      tftypes.NewStringDateValue("2016-01-01.02:00:00"),
			expectedEqual: true,
		},
		"same value": {
			priorValue:    tftypes.NewStringDateValue("2021-12-11.10:09:08"),
			newValue:      tftypes.NewStringDateValue("2021-12-11.10:09:08"),
			expectedEqual: true,
		},
		"different day": {
			priorValue:    tftypes.NewStringDateValue("2022-06-01.00:00:00"),
			newValue:      tftypes.NewStringDateValue("2022-06-02.00:00:00"),
			expectedEqual: false,
		},
		"different month": {
			priorValue:    tftypes.NewStringDateValue("2022-10-10.10:09:08"),
			newValue:      tftypes.NewStringDateValue("2022-1-10.10:09:08"),
			expectedEqual: false,
		},
		"different year": {
			priorValue:    tftypes.NewStringDateValue("2022-06-01.00:00:00"),
			newValue:      tftypes.NewStringDateValue("2023-06-01.00:00:00"),
			expectedEqual: false,
		},
		"different time": {
			priorValue:    tftypes.NewStringDateValue("2022-06-01.00:00:00"),
			newValue:      tftypes.NewStringDateValue("2022-6-1.00:00:01"),
			expectedEqual: false,
		},
		"not a date": {
			priorValue:    tftypes.NewStringDateValue("not a date"),
			newValue:      tftypes.NewStringDateValue("not a date"),
			expectedEqual: true,
		},
		"not a date and different": {
			priorValue:    tftypes.NewStringDateValue("not a date"),
			newValue:      tftypes.NewStringDateValue("not the same"),
			expectedEqual: false,
		},
		"unexpected value type": {
			priorValue:    tftypes.NewStringDateValue("2022-06-01.00:00:00"),
			newValue:      basetypes.NewStringValue("2022-06-01.00:00:00"),
			expectedEqual: false,
			expectedError: true,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			equal, diags := testCase.priorValue.StringSemanticEquals(context.Background(), testCase.newValue)

			if diags.HasError() != testCase.expectedError {
				t.Errorf("unexpected diagnostics: want error %t, got %v", testCase.expectedError, diags)
			}
			if equal != testCase.expectedEqual {
				t.Errorf("unexpected equality: want %t, got %t", testCase.expectedEqual, equal)
			}
		})
	}
}

func TestDateNullUnknown(t *testing.T) {
	t.Parallel()

	if !tftypes.NewStringDateNull().IsNull() {
		t.Errorf("expected null value")
	}
	if !tftypes.NewStringDateUnknown().IsUnknown() {
		t.Errorf("expected unknown value")
	}
	if value := tftypes.NewStringDateFromString(basetypes.NewStringNull()); !value.IsNull() {
		t.Errorf("expected null value, got %#v", value)
	}
	if value := tftypes.NewStringDateFromString(
		basetypes.NewStringValue("2022-6-1.00:00:00"),
	); value.ValueString() != "2022-6-1.00:00:00" {
		t.Errorf("unexpected value: %#v", value)
	}
}
