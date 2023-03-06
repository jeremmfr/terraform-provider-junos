package tfplanmodifier_test

import (
	"context"
	"testing"

	"github.com/jeremmfr/terraform-provider-junos/internal/tfplanmodifier"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestStringDefault(t *testing.T) {
	t.Parallel()

	testSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"testattr": schema.StringAttribute{},
		},
	}

	nullPlan := tfsdk.Plan{
		Schema: testSchema,
		Raw: tftypes.NewValue(
			testSchema.Type().TerraformType(context.Background()),
			nil,
		),
	}

	testConfig := func(value types.String) tfsdk.Config {
		tfValue, err := value.ToTerraformValue(context.Background())
		if err != nil {
			panic("ToTerraformValue error: " + err.Error())
		}

		return tfsdk.Config{
			Schema: testSchema,
			Raw: tftypes.NewValue(
				testSchema.Type().TerraformType(context.Background()),
				map[string]tftypes.Value{
					"testattr": tfValue,
				},
			),
		}
	}

	testPlan := func(value types.String) tfsdk.Plan {
		tfValue, err := value.ToTerraformValue(context.Background())
		if err != nil {
			panic("ToTerraformValue error: " + err.Error())
		}

		return tfsdk.Plan{
			Schema: testSchema,
			Raw: tftypes.NewValue(
				testSchema.Type().TerraformType(context.Background()),
				map[string]tftypes.Value{
					"testattr": tfValue,
				},
			),
		}
	}

	testCases := map[string]struct {
		request      planmodifier.StringRequest
		expected     *planmodifier.StringResponse
		defaultValue string
	}{
		"Config->null Plan->unknown": {
			// resource creation with optional/computed attribute not set
			request: planmodifier.StringRequest{
				Config:      testConfig(types.StringNull()),
				ConfigValue: types.StringNull(),
				Plan:        testPlan(types.StringUnknown()),
				PlanValue:   types.StringUnknown(),
			},
			expected: &planmodifier.StringResponse{
				PlanValue: types.StringValue("default"),
			},
			defaultValue: "default",
		},
		"Config->unknown Plan->unknown": {
			// resource creation with optional/computed attribute set with unknown value
			request: planmodifier.StringRequest{
				Config:      testConfig(types.StringUnknown()),
				ConfigValue: types.StringUnknown(),
				Plan:        testPlan(types.StringUnknown()),
				PlanValue:   types.StringUnknown(),
			},
			expected: &planmodifier.StringResponse{
				PlanValue: types.StringUnknown(),
			},
			defaultValue: "default",
		},
		"Config->set Plan->set": {
			// resource creation with optional/computed attribute set with known value
			request: planmodifier.StringRequest{
				Config:      testConfig(types.StringValue("set")),
				ConfigValue: types.StringValue("test"),
				Plan:        testPlan(types.StringValue("test")),
				PlanValue:   types.StringValue("test"),
			},
			expected: &planmodifier.StringResponse{
				PlanValue: types.StringValue("test"),
			},
			defaultValue: "default",
		},
		"RawPlan->null": {
			// resource destroy
			request: planmodifier.StringRequest{
				Config:      testConfig(types.StringNull()),
				ConfigValue: types.StringNull(),
				Plan:        nullPlan,
				PlanValue:   types.StringNull(),
			},
			expected: &planmodifier.StringResponse{
				PlanValue: types.StringNull(),
			},
			defaultValue: "default",
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &planmodifier.StringResponse{
				PlanValue: testCase.request.PlanValue,
			}

			tfplanmodifier.StringDefault(testCase.defaultValue).PlanModifyString(context.Background(), testCase.request, resp)

			if diff := cmp.Diff(testCase.expected, resp); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
