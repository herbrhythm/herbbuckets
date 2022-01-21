package buckets

import (
	"fmt"
	"herbbuckets/modules/app/bucketconfig"
	"herbbuckets/modules/bucket"
	"herbbuckets/modules/bucket/localbucket"
)

func CreateBucket(config *bucketconfig.Bucket) (bucket.Bucket, error) {
	switch config.Type {
	case "", localbucket.BucketType:
		return localbucket.Factory(config.Name, config.Config)
	}
	return nil, fmt.Errorf("unknown bucket type [%s]", config.Type)
}
