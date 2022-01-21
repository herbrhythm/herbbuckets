package pathcleaner

import (
	"path"
	"strings"
)

func Clean(filename string) string {
	return strings.TrimPrefix(path.Clean("/"+filename), "/")
}
