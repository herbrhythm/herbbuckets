package bucket

import (
	"herbbuckets/modules/app"
	"herbbuckets/modules/app/bucketconfig"
	"io"
	"time"

	"github.com/herb-go/herb-go/modules/cors"
)

type Fileinfo struct {
	Size int64
}

func NewFileinfo() *Fileinfo {
	return &Fileinfo{}
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

type Bucket struct {
	Name       string
	Type       string
	Enabled    bool
	DateFormat string
	Lifetime   time.Duration
	Cors       *cors.CORS
	Referrer   []string
	BaseURL    string
	Engine     Engine
}

func (b *Bucket) Verify() error {
	return nil
}
func (b *Bucket) InitWith(config *bucketconfig.Config) error {
	b.Name = config.Name
	b.Type = config.Type
	b.Enabled = config.Enabled
	b.DateFormat = config.DateFormat
	if b.DateFormat == "" {
		b.DateFormat = app.System.DateFormat
	}
	b.Lifetime = time.Duration(config.LifetimeInSeconds) * time.Second
	b.Cors = &config.Cors
	b.Referrer = config.Referrer
	b.BaseURL = config.BaseURL
	if b.BaseURL == "" {
		b.BaseURL = app.HTTP.Config.BaseURL
	}
	return b.Verify()
}
func New() *Bucket {
	return &Bucket{}
}

type Engine interface {
	GrantDownloadURL(b *Bucket, objectname string, opt *Options) (url string, err error)
	GrantUploadURL(b *Bucket, objectname string, opt *Options) (uploadurl string, err error)
	RemoveFile(b *Bucket, objectname string) error
	Permanent() bool
	ThirdpartyUpload() bool
	ThirdpartyDownload() bool
	GetVerifier() *Verifier
	BucketType() string
	Download(b *Bucket, objectname string, w io.Writer) (err error)
	Upload(b *Bucket, objectname string, r io.Reader) (err error)
	GetFileinfo(b *Bucket, objectname string) (info *Fileinfo, err error)
}

type Builder func(b *bucketconfig.Bucket, loader func(v interface{}) error) error
