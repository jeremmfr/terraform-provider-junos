package tfdata

import (
	"reflect"
	"strings"

	"github.com/jeremmfr/go-utils/basicalter"
)

const fieldTag = "tfdata"

func tagsOfStructField(f reflect.StructField) []string {
	tags := strings.Split(f.Tag.Get(fieldTag), ",")

	basicalter.ReplaceInSliceWith(tags, strings.TrimSpace)

	return tags
}
