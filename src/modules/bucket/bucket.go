package bucket

import (
	"herbbuckets/modules/app"
	"herbbuckets/modules/app/bucketconfig"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"time"

	"github.com/herb-go/herb/middleware/cors"
)

var NameRegexp = regexp.MustCompile(`^[a-zA-Z0-9\-_\.]{1,64}$`)

type Fileinfo struct {
	Size    int64
	Modtime int64
}

func NewFileinfo() *Fileinfo {
	return &Fileinfo{}
}

type DownloadInfo struct {
	URL       string
	ExpiredAt int64
	Permanent bool
}

func NewDownloadInfo() *DownloadInfo {
	return &DownloadInfo{}
}

type WebuploadInfo struct {
	ID             string
	UploadURL      string
	PreviewURL     string
	Permanent      bool
	Bucket         string
	Object         string
	UploadType     string
	SizeLimit      int64
	ExpiredAt      int64
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
	SizeLimit int64
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
	BaseURL    *url.URL
	SizeLimit  int64
	Engine     Engine
}

func (b *Bucket) Verify() error {
	return nil
}
func (b *Bucket) Join(ele ...string) string {
	u := *b.BaseURL
	u.Path = path.Join(u.Path, path.Join(ele...))
	return u.String()
}
func (b *Bucket) InitWith(config *bucketconfig.Config) error {
	var err error
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
	b.SizeLimit = config.SizeLimit
	b.Cors = &config.Cors
	b.Referrer = config.Referrer
	burl := config.BaseURL
	if burl == "" {
		burl = app.HTTP.Config.BaseURL
	}
	if burl != "" {
		b.BaseURL, err = url.Parse(burl)
		if err != nil {
			return err
		}
	}
	b.Prefix = config.Prefix
	return b.Verify()
}
func New() *Bucket {
	return &Bucket{}
}

type Engine interface {
	GrantDownloadInfo(b *Bucket, objectname string, opt *Options) (info *DownloadInfo, err error)
	GrantUploadInfo(b *Bucket, id string, objectname string, opt *Options) (info *WebuploadInfo, err error)
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
