package tfvalidator

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type StringNumberRangeValidator struct {
	min  int
	max  int
	name string
}

func StringNumberRange(min, max int) StringNumberRangeValidator {
	return StringNumberRangeValidator{
		min: min,
		max: max,
	}
}

func (v StringNumberRangeValidator) WithNameInError(name string) StringNumberRangeValidator {
	v.name = name

	return v
}

func (v StringNumberRangeValidator) Description(_ context.Context) string {
	if v.name != "" {
		return fmt.Sprintf("Must be a valid %s or %s range between %d and %d.", v.name, v.name, v.min, v.max)
	}

	return fmt.Sprintf("Must be a valid number or number range between %d and %d.", v.min, v.max)
}

func (v StringNumberRangeValidator) MarkdownDescription(_ context.Context) string {
	if v.name != "" {
		return fmt.Sprintf("Must be a valid %s or %s range between %d and %d.", v.name, v.name, v.min, v.max)
	}

	return fmt.Sprintf("Must be a valid number or number range between %d and %d.", v.min, v.max)
}

func (v StringNumberRangeValidator) ValidateString(
	_ context.Context, req validator.StringRequest, resp *validator.StringResponse,
) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}
	summary := "Invalid Number or Number Range"
	if v.name != "" {
		summary = fmt.Sprintf("Invalid %s or %s Range", v.name, v.name)
	}

	value := req.ConfigValue.ValueString()
	if ok := regexp.MustCompile(`^\d+(-\d+)?$`).MatchString(value); !ok {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			summary,
			fmt.Sprintf(`expected value match regular expression \d+(-\d+)?, got %q`, value),
		)

		return
	}

	vSplit := strings.Split(value, "-")
	low, err := strconv.Atoi(vSplit[0])
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			summary,
			err.Error(),
		)

		return
	}

	high := low
	if len(vSplit) > 1 {
		high, err = strconv.Atoi(vSplit[1])
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				summary,
				err.Error(),
			)

			return
		}
	}

	if low > high {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			summary,
			fmt.Sprintf("low(%d) in %q bigger than high(%d)", low, value, high),
		)
	}
	if low < v.min {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			summary,
			fmt.Sprintf("low(%d) in %q is too small (min %d)", low, value, v.min),
		)
	}
	if high > v.max {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			summary,
			fmt.Sprintf("high(%d) in %q is too big (max %d)", high, value, v.max),
		)
	}
}
