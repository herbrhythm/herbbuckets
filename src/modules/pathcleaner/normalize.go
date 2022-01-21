package pathcleaner

import "net/url"

func Normalize(raw string) string {
	return url.PathEscape(raw)
}
