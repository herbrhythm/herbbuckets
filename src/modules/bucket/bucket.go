package bucket

import "io"

type Fileinfo struct {
}

type Verifier struct {
	CodeMin  int
	CodeMax  int
	Required string
}
type Options struct {
	Appid     string
	Secret    string
	ExpiredAt int64
}
type Bucket interface {
	GrantDownloadURL(bucketname string, objectname string, opt *Options) (url string, err error)
	GrantUploadURL(bucketname string, objectname string, opt *Options) (uploadurl string, err error)
	RemoveFile(bucketname string, objectname string) error
	Permanent() bool
	ThirdpartyUpload() bool
	GetVerifier() *Verifier
	BucketType() string
	Download(bucketname string, objectname string, w io.Writer) (err error)
	Upload(bucketname string, objectname string, r io.Reader) (err error)
	GetFileinfo(bucketname string, objectname string) (info *Fileinfo, err error)
}

type Factory func(bucketname string, loader func(v interface{}) error) (Bucket, error)
