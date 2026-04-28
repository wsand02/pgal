package models

import "strings"

func Schemas() string {
	return strings.Join([]string{FolderSchema(), FileSchema()}, "\n")
}
