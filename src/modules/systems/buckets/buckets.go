package buckets

import "herbbuckets/modules/bucket"

var Buckets = map[string]*bucket.Bucket{}

func GetBucket(name string) *bucket.Bucket {
	return Buckets[name]
}
