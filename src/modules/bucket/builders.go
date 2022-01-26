package bucket

var Builders = map[string]Builder{}

type Builder func(b *Bucket, loader func(v interface{}) error) error
