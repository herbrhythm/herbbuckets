package pathcleaner

import (
	"herbbuckets/modules/app"
	"herbbuckets/modules/bucket"
	"herbbuckets/modules/uniqueid"
	"path"
)

func CreateObjectID(bu *bucket.Bucket, filename string) string {
	var result = make([]string, 0, 4)
	if bu.Prefix != "" {
		result = append(result, bu.Prefix)
	}
	dateformat := bu.DateFormat
	if dateformat == "" {
		dateformat = app.System.DateFormat
	}
	if dateformat != "" {
		result = append(result, app.Time.FormatNow(dateformat))
	}
	result = append(result, uniqueid.MustGenerateID())
	result = append(result, filename)
	return Clean(path.Join(result...))
}
