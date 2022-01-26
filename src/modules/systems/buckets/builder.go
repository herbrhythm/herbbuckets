package buckets

import (
	"fmt"
	"herbbuckets/modules/app/bucketconfig"
	"herbbuckets/modules/bucket"
)

func buildBucket(b *bucket.Bucket, loader func(v interface{}) error) error {
	builder := bucket.Builders[b.Type]
	if builder == nil {
		return fmt.Errorf("unknown bucket type [%s]", b.Type)
	}
	return builder(b, loader)
}
func CreateBucket(config *bucketconfig.Config) (*bucket.Bucket, error) {
	b := bucket.New()
	err := b.InitWith(config)
	if err != nil {
		return nil, err
	}
	err = buildBucket(b, config.Config)
	if err != nil {
		return nil, err
	}
	return b, nil
}
