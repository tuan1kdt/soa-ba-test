package querybuilder

import (
	"strings"

	"gorm.io/gorm"
)

type Builder interface {
	Build(tx *gorm.DB) *gorm.DB
}

func notEmptyString(s string) bool {
	return strings.TrimSpace(s) != ""
}
