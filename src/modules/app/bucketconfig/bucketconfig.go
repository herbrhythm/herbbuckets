package bucketconfig

type Bucket struct {
	Name    string
	Type    string
	Enabled bool
	Config  func(interface{}) error
}
