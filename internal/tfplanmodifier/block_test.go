package tfplanmodifier_test

import (
	"context"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/tfplanmodifier"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestBlockRemoveNull(t *testing.T) {
	t.Parallel()

	testSchema := schema.Schema{
		Blocks: map[string]schema.Block{
			"testattrobject": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{"testattr": schema.StringAttribute{}},
			},
		},
	}

	nullPlan := tfsdk.Plan{
		Schema: testSchema,
		Raw: tftypes.NewValue(
			testSchema.Type().TerraformType(context.Background()),
			nil,
		),
	}

	testConfig := func(value types.Object) tfsdk.Config {
		tfValue, err := value.ToTerraformValue(context.Background())
		if err != nil {
			panic("ToTerraformValue error: " + err.Error())
		}

		return tfsdk.Config{
			Schema: testSchema,
			Raw: tftypes.NewValue(
				testSchema.Type().TerraformType(context.Background()),
				map[string]tftypes.Value{
					"testattrobject": tfValue,
				},
			),
		}
	}

	testPlan := func(value types.Object) tfsdk.Plan {
		tfValue, err := value.ToTerraformValue(context.Background())
		if err != nil {
			panic("ToTerraformValue error: " + err.Error())
		}

		return tfsdk.Plan{
			Schema: testSchema,
			Raw: tftypes.NewValue(
				testSchema.Type().TerraformType(context.Background()),
				map[string]tftypes.Value{
					"testattrobject": tfValue,
				},
			),
		}
	}

	testCases := map[string]struct {
		request  planmodifier.ObjectRequest
		expected *planmodifier.ObjectResponse
	}{
		"Config->null Plan->unknown": {
			// resource creation with optional block attribute not set but unknown in plan
			// Reference: https://github.com/hashicorp/terraform/issues/32460
			request: planmodifier.ObjectRequest{
				Config: testConfig(types.ObjectNull(
					map[string]attr.Type{"testattr": types.StringType},
				)),
				ConfigValue: types.ObjectNull(
					map[string]attr.Type{"testattr": types.StringType},
				),
				Plan: testPlan(types.ObjectUnknown(
					map[string]attr.Type{"testattr": types.StringType},
				)),
				PlanValue: types.ObjectUnknown(
					map[string]attr.Type{"testattr": types.StringType},
				),
			},
			expected: &planmodifier.ObjectResponse{
				PlanValue: types.ObjectNull(map[string]attr.Type{"testattr": types.StringType}),
			},
		},
		"Config->null Plan->set": {
			// resource creation with optional block attribute not set but known in plan
			// Reference: https://github.com/hashicorp/terraform/issues/32460
			request: planmodifier.ObjectRequest{
				Config: testConfig(types.ObjectNull(
					map[string]attr.Type{"testattr": types.StringType},
				)),
				ConfigValue: types.ObjectNull(
					map[string]attr.Type{"testattr": types.StringType},
				),
				Plan: testPlan(types.ObjectValueMust(
					map[string]attr.Type{"testattr": types.StringType},
					map[string]attr.Value{"testattr": types.StringValue("test")},
				)),
				PlanValue: types.ObjectValueMust(
					map[string]attr.Type{"testattr": types.StringType},
					map[string]attr.Value{"testattr": types.StringValue("test")},
				),
			},
			expected: &planmodifier.ObjectResponse{
				PlanValue: types.ObjectNull(map[string]attr.Type{"testattr": types.StringType}),
			},
		},
		"Config->unknown Plan->unknown": {
			// resource creation with optional block attribute set with unknown value
			request: planmodifier.ObjectRequest{
				Config: testConfig(types.ObjectUnknown(
					map[string]attr.Type{"testattr": types.StringType},
				)),
				ConfigValue: types.ObjectUnknown(
					map[string]attr.Type{"testattr": types.StringType},
				),
				Plan: testPlan(types.ObjectUnknown(
					map[string]attr.Type{"testattr": types.StringType},
				)),
				PlanValue: types.ObjectUnknown(
					map[string]attr.Type{"testattr": types.StringType},
				),
			},
			expected: &planmodifier.ObjectResponse{
				PlanValue: types.ObjectUnknown(
					map[string]attr.Type{"testattr": types.StringType},
				),
			},
		},
		"Config->set Plan->set": {
			// resource creation with optional block attribute set with known value
			request: planmodifier.ObjectRequest{
				Config: testConfig(types.ObjectValueMust(
					map[string]attr.Type{"testattr": types.StringType},
					map[string]attr.Value{"testattr": types.StringValue("test")},
				)),
				ConfigValue: types.ObjectValueMust(
					map[string]attr.Type{"testattr": types.StringType},
					map[string]attr.Value{"testattr": types.StringValue("test")},
				),
				Plan: testPlan(types.ObjectValueMust(
					map[string]attr.Type{"testattr": types.StringType},
					map[string]attr.Value{"testattr": types.StringValue("test")},
				)),
				PlanValue: types.ObjectValueMust(
					map[string]attr.Type{"testattr": types.StringType},
					map[string]attr.Value{"testattr": types.StringValue("test")},
				),
			},
			expected: &planmodifier.ObjectResponse{
				PlanValue: types.ObjectValueMust(
					map[string]attr.Type{"testattr": types.StringType},
					map[string]attr.Value{"testattr": types.StringValue("test")},
				),
			},
		},
		"RawPlan->null": {
			// resource destroy
			request: planmodifier.ObjectRequest{
				Config: testConfig(types.ObjectNull(
					map[string]attr.Type{"testattr": types.StringType},
				)),
				ConfigValue: types.ObjectNull(
					map[string]attr.Type{"testattr": types.StringType},
				),
				Plan: nullPlan,
				PlanValue: types.ObjectNull(
					map[string]attr.Type{"testattr": types.StringType},
				),
			},
			expected: &planmodifier.ObjectResponse{
				PlanValue: types.ObjectNull(
					map[string]attr.Type{"testattr": types.StringType},
				),
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &planmodifier.ObjectResponse{
				PlanValue: testCase.request.PlanValue,
			}

			tfplanmodifier.BlockRemoveNull().PlanModifyObject(context.Background(), testCase.request, resp)

			if diff := cmp.Diff(testCase.expected, resp); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestBlockSetUnsetRequireReplace(t *testing.T) {
	t.Parallel()

	testSchema := schema.Schema{
		Blocks: map[string]schema.Block{
			"testattrobject": schema.SingleNestedBlock{
				Attributes: map[string]schema.Attribute{"testattr": schema.StringAttribute{}},
			},
		},
	}

	nullPlan := tfsdk.Plan{
		Schema: testSchema,
		Raw: tftypes.NewValue(
			testSchema.Type().TerraformType(context.Background()),
			nil,
		),
	}
	nullState := tfsdk.State{
		Schema: testSchema,
		Raw: tftypes.NewValue(
			testSchema.Type().TerraformType(context.Background()),
			nil,
		),
	}

	testConfig := func(value types.Object) tfsdk.Config {
		tfValue, err := value.ToTerraformValue(context.Background())
		if err != nil {
			panic("ToTerraformValue error: " + err.Error())
		}

		return tfsdk.Config{
			Schema: testSchema,
			Raw: tftypes.NewValue(
				testSchema.Type().TerraformType(context.Background()),
				map[string]tftypes.Value{
					"testattrobject": tfValue,
				},
			),
		}
	}

	testPlan := func(value types.Object) tfsdk.Plan {
		tfValue, err := value.ToTerraformValue(context.Background())
		if err != nil {
			panic("ToTerraformValue error: " + err.Error())
		}

		return tfsdk.Plan{
			Schema: testSchema,
			Raw: tftypes.NewValue(
				testSchema.Type().TerraformType(context.Background()),
				map[string]tftypes.Value{
					"testattrobject": tfValue,
				},
			),
		}
	}

	testState := func(value types.Object) tfsdk.State {
		tfValue, err := value.ToTerraformValue(context.Background())
		if err != nil {
			panic("ToTerraformValue error: " + err.Error())
		}

		return tfsdk.State{
			Schema: testSchema,
			Raw: tftypes.NewValue(
				testSchema.Type().TerraformType(context.Background()),
				map[string]tftypes.Value{
					"testattrobject": tfValue,
				},
			),
		}
	}

	testCases := map[string]struct {
		request  planmodifier.ObjectRequest
		expected *planmodifier.ObjectResponse
	}{
		"Config->null State->null": {
			request: planmodifier.ObjectRequest{
				Config: testConfig(types.ObjectNull(
					map[string]attr.Type{"testattr": types.StringType},
				)),
				ConfigValue: types.ObjectNull(
					map[string]attr.Type{"testattr": types.StringType},
				),
				Plan: testPlan(types.ObjectNull(
					map[string]attr.Type{"testattr": types.StringType},
				)),
				PlanValue: types.ObjectNull(
					map[string]attr.Type{"testattr": types.StringType},
				),
				State: testState(types.ObjectNull(
					map[string]attr.Type{"testattr": types.StringType},
				)),
				StateValue: types.ObjectNull(
					map[string]attr.Type{"testattr": types.StringType},
				),
			},
			expected: &planmodifier.ObjectResponse{
				PlanValue:       types.ObjectNull(map[string]attr.Type{"testattr": types.StringType}),
				RequiresReplace: false,
			},
		},
		"Config->null State->set": {
			request: planmodifier.ObjectRequest{
				Config: testConfig(types.ObjectNull(
					map[string]attr.Type{"testattr": types.StringType},
				)),
				ConfigValue: types.ObjectNull(
					map[string]attr.Type{"testattr": types.StringType},
				),
				Plan: testPlan(types.ObjectNull(
					map[string]attr.Type{"testattr": types.StringType},
				)),
				PlanValue: types.ObjectNull(
					map[string]attr.Type{"testattr": types.StringType},
				),
				State: testState(types.ObjectValueMust(
					map[string]attr.Type{"testattr": types.StringType},
					map[string]attr.Value{"testattr": types.StringValue("test")},
				)),
				StateValue: types.ObjectValueMust(
					map[string]attr.Type{"testattr": types.StringType},
					map[string]attr.Value{"testattr": types.StringValue("test")},
				),
			},
			expected: &planmodifier.ObjectResponse{
				PlanValue:       types.ObjectNull(map[string]attr.Type{"testattr": types.StringType}),
				RequiresReplace: true,
			},
		},
		"Config->unknown State->set": {
			request: planmodifier.ObjectRequest{
				Config: testConfig(types.ObjectUnknown(
					map[string]attr.Type{"testattr": types.StringType},
				)),
				ConfigValue: types.ObjectUnknown(
					map[string]attr.Type{"testattr": types.StringType},
				),
				Plan: testPlan(types.ObjectUnknown(
					map[string]attr.Type{"testattr": types.StringType},
				)),
				PlanValue: types.ObjectUnknown(
					map[string]attr.Type{"testattr": types.StringType},
				),
				State: testState(types.ObjectValueMust(
					map[string]attr.Type{"testattr": types.StringType},
					map[string]attr.Value{"testattr": types.StringValue("test")},
				)),
				StateValue: types.ObjectValueMust(
					map[string]attr.Type{"testattr": types.StringType},
					map[string]attr.Value{"testattr": types.StringValue("test")},
				),
			},
			expected: &planmodifier.ObjectResponse{
				PlanValue: types.ObjectUnknown(
					map[string]attr.Type{"testattr": types.StringType},
				),
				RequiresReplace: false,
			},
		},
		"Config->set State->null": {
			request: planmodifier.ObjectRequest{
				Config: testConfig(types.ObjectValueMust(
					map[string]attr.Type{"testattr": types.StringType},
					map[string]attr.Value{"testattr": types.StringValue("test")},
				)),
				ConfigValue: types.ObjectValueMust(
					map[string]attr.Type{"testattr": types.StringType},
					map[string]attr.Value{"testattr": types.StringValue("test")},
				),
				Plan: testPlan(types.ObjectValueMust(
					map[string]attr.Type{"testattr": types.StringType},
					map[string]attr.Value{"testattr": types.StringValue("test")},
				)),
				PlanValue: types.ObjectValueMust(
					map[string]attr.Type{"testattr": types.StringType},
					map[string]attr.Value{"testattr": types.StringValue("test")},
				),
				State: testState(types.ObjectNull(
					map[string]attr.Type{"testattr": types.StringType},
				)),
				StateValue: types.ObjectNull(
					map[string]attr.Type{"testattr": types.StringType},
				),
			},
			expected: &planmodifier.ObjectResponse{
				PlanValue: types.ObjectValueMust(
					map[string]attr.Type{"testattr": types.StringType},
					map[string]attr.Value{"testattr": types.StringValue("test")},
				),
				RequiresReplace: true,
			},
		},
		"Config->set State->set": {
			request: planmodifier.ObjectRequest{
				Config: testConfig(types.ObjectValueMust(
					map[string]attr.Type{"testattr": types.StringType},
					map[string]attr.Value{"testattr": types.StringValue("test2")},
				)),
				ConfigValue: types.ObjectValueMust(
					map[string]attr.Type{"testattr": types.StringType},
					map[string]attr.Value{"testattr": types.StringValue("test2")},
				),
				Plan: testPlan(types.ObjectValueMust(
					map[string]attr.Type{"testattr": types.StringType},
					map[string]attr.Value{"testattr": types.StringValue("test2")},
				)),
				PlanValue: types.ObjectValueMust(
					map[string]attr.Type{"testattr": types.StringType},
					map[string]attr.Value{"testattr": types.StringValue("test2")},
				),
				State: testState(types.ObjectValueMust(
					map[string]attr.Type{"testattr": types.StringType},
					map[string]attr.Value{"testattr": types.StringValue("test")},
				)),
				StateValue: types.ObjectValueMust(
					map[string]attr.Type{"testattr": types.StringType},
					map[string]attr.Value{"testattr": types.StringValue("test")},
				),
			},
			expected: &planmodifier.ObjectResponse{
				PlanValue: types.ObjectValueMust(
					map[string]attr.Type{"testattr": types.StringType},
					map[string]attr.Value{"testattr": types.StringValue("test2")},
				),
				RequiresReplace: false,
			},
		},
		"RawState->null": {
			// resource create
			request: planmodifier.ObjectRequest{
				Config: testConfig(types.ObjectValueMust(
					map[string]attr.Type{"testattr": types.StringType},
					map[string]attr.Value{"testattr": types.StringValue("test")},
				)),
				ConfigValue: types.ObjectValueMust(
					map[string]attr.Type{"testattr": types.StringType},
					map[string]attr.Value{"testattr": types.StringValue("test")},
				),
				Plan: testPlan(types.ObjectValueMust(
					map[string]attr.Type{"testattr": types.StringType},
					map[string]attr.Value{"testattr": types.StringValue("test")},
				)),
				PlanValue: types.ObjectValueMust(
					map[string]attr.Type{"testattr": types.StringType},
					map[string]attr.Value{"testattr": types.StringValue("test")},
				),
				State: nullState,
			},
			expected: &planmodifier.ObjectResponse{
				PlanValue: types.ObjectValueMust(
					map[string]attr.Type{"testattr": types.StringType},
					map[string]attr.Value{"testattr": types.StringValue("test")},
				),
				RequiresReplace: false,
			},
		},
		"RawPlan->null": {
			// resource destroy
			request: planmodifier.ObjectRequest{
				Config: testConfig(types.ObjectNull(
					map[string]attr.Type{"testattr": types.StringType},
				)),
				ConfigValue: types.ObjectNull(
					map[string]attr.Type{"testattr": types.StringType},
				),
				Plan: nullPlan,
				PlanValue: types.ObjectNull(
					map[string]attr.Type{"testattr": types.StringType},
				),
				State: testState(types.ObjectValueMust(
					map[string]attr.Type{"testattr": types.StringType},
					map[string]attr.Value{"testattr": types.StringValue("test")},
				)),
				StateValue: types.ObjectValueMust(
					map[string]attr.Type{"testattr": types.StringType},
					map[string]attr.Value{"testattr": types.StringValue("test")},
				),
			},
			expected: &planmodifier.ObjectResponse{
				PlanValue: types.ObjectNull(
					map[string]attr.Type{"testattr": types.StringType},
				),
				RequiresReplace: false,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &planmodifier.ObjectResponse{
				PlanValue: testCase.request.PlanValue,
			}

			tfplanmodifier.BlockSetUnsetRequireReplace().PlanModifyObject(context.Background(), testCase.request, resp)

			if diff := cmp.Diff(testCase.expected, resp); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
