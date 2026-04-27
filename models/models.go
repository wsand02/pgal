package models

import (
	"strings"
)

func Schemas() string {
	return strings.Join([]string{folderSchema(), fileSchema()}, "\n")
}
