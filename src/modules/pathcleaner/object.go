package pathcleaner

import (
	"herbbuckets/modules/app"
	"herbbuckets/modules/uniqueid"
	"path"
)

func CreateObjectID(dateformat string, tag string, filename string) string {
	var result = make([]string, 0, 4)
	if dateformat == "" {
		dateformat = app.System.DataFormat
	}
	if dateformat != "" {
		result = append(result, app.Time.FormatNow(dateformat))
	}
	result = append(result, uniqueid.MustGenerateID())
	if tag != "" {
		result = append(result, tag)
	}
	result = append(result, filename)
	return Clean(path.Join(result...))
}
