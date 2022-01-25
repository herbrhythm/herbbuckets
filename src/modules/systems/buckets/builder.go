package buckets

import (
	"fmt"
	"herbbuckets/modules/app/bucketconfig"
	"herbbuckets/modules/bucket"
	"herbbuckets/modules/bucket/localbucket"
)

func buildBucket(b *bucket.Bucket, loader func(v interface{}) error) error {
	switch b.Type {
	case "", localbucket.BucketType:
		return localbucket.Builder(b, loader)
	}
	return fmt.Errorf("unknown bucket type [%s]", b.Type)

}
func CreateBucket(config *bucketconfig.Config) (*bucket.Bucket, error) {
	b := bucket.New()
	err := buildBucket(b, config.Config)
	if err != nil {
		return nil, err
	}
	return b, nil
}
