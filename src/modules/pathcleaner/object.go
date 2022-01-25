package pathcleaner

import (
	"herbbuckets/modules/app"
	"herbbuckets/modules/app/bucketconfig"
	"herbbuckets/modules/uniqueid"
	"path"
)

func CreateObjectID(bu *bucketconfig.Bucket, tag string, filename string) string {
	var result = make([]string, 0, 4)
	dateformat := bu.DateFormat
	if dateformat == "" {
		dateformat = app.System.DateFormat
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
