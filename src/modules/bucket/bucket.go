package bucket

import (
	"herbbuckets/modules/app"
	"herbbuckets/modules/app/bucketconfig"
	"io"
	"net/http"
	"time"

	"github.com/herb-go/herb/middleware/cors"
)

type Fileinfo struct {
	Size int64
}

func NewFileinfo() *Fileinfo {
	return &Fileinfo{}
}

type WebuploadInfo struct {
	UploadURL      string
	PreviewURL     string
	Permanent      bool
	Bucket         string
	Objcet         string
	UploadType     string
	Sizelimit      int64
	PostBody       map[string]string
	FileField      string
	SuccessCodeMin int
	SuccessCodeMax int
}

func NewWebuploadInfo() *WebuploadInfo {
	return &WebuploadInfo{
		PostBody: map[string]string{},
	}
}

type Options struct {
	Appid     string
	Secret    string
	Sizelimit int64
	Lifetime  time.Duration
}

func NewOptions() *Options {
	return &Options{}
}

type Bucket struct {
	Name       string
	Type       string
	Disabled   bool
	Prefix     string
	DateFormat string
	Lifetime   time.Duration
	Cors       *cors.CORS
	Referrer   []string
	BaseURL    string
	Sizelimit  int64
	Engine     Engine
}

func (b *Bucket) Verify() error {
	return nil
}
func (b *Bucket) InitWith(config *bucketconfig.Config) error {
	b.Name = config.Name
	b.Type = config.Type
	b.Disabled = config.Disabled
	b.DateFormat = config.DateFormat
	if b.DateFormat == "" {
		b.DateFormat = app.System.DateFormat
	}
	b.Lifetime = time.Duration(config.LifetimeInSeconds) * time.Second
	if b.Lifetime <= 0 {
		b.Lifetime = time.Duration(app.System.LifetimeInSeconds) * time.Second
	}
	b.Sizelimit = config.Sizelimit
	b.Cors = &config.Cors
	b.Referrer = config.Referrer
	b.BaseURL = config.BaseURL
	if b.BaseURL == "" {
		b.BaseURL = app.HTTP.Config.BaseURL
	}
	b.Prefix = config.Prefix
	return b.Verify()
}
func New() *Bucket {
	return &Bucket{}
}

type Engine interface {
	GrantDownloadURL(b *Bucket, objectname string, opt *Options) (url string, err error)
	GrantUploadInfo(b *Bucket, objectname string, opt *Options) (info *WebuploadInfo, err error)
	RemoveFile(b *Bucket, objectname string) error
	Permanent() bool
	ThirdpartyUpload() bool
	ThirdpartyDownload() bool
	BucketType() string
	ServeHTTPDownload(w http.ResponseWriter, r *http.Request)
	Download(b *Bucket, objectname string) (r io.ReadCloser, err error)
	Upload(b *Bucket, objectname string) (w io.WriteCloser, err error)
	GetFileinfo(b *Bucket, objectname string) (info *Fileinfo, err error)
	Start() error
	Stop() error
}
