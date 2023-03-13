package tfdata

import (
	"fmt"

	"github.com/jeremmfr/terraform-provider-junos/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func ConvAtoi64Value(str string,
) (
	basetypes.Int64Value, error,
) {
	v, err := utils.ConvAtoi64(str)
	if err != nil {
		return types.Int64Null(), fmt.Errorf("converting value from '%s' to integer: %w", str, err)
	}

	return types.Int64Value(v), nil
}
