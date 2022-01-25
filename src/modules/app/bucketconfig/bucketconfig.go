package bucketconfig

type Config struct {
	Bucket
	Config func(interface{}) error
}
