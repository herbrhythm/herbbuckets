package bucket

import (
	"herbbuckets/modules/app"
	"herbbuckets/modules/app/bucketconfig"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"time"

	"github.com/herb-go/herb/middleware/cors"
	"github.com/herb-go/herbsecurity/secret"
	"github.com/herb-go/herbsecurity/secret/hasher"
	"github.com/herb-go/herbsecurity/secret/hasher/urlencodesign"
)

var NameRegexp = regexp.MustCompile(`^[a-zA-Z0-9\-_\.]{1,64}$`)

type Bucket struct {
	Name     string
	Type     string
	Disabled bool
	Prefix   string

	DateFormat string
	Lifetime   time.Duration
	Offset     time.Duration
	Cors       *cors.CORS
	Referrer   []string
	BaseURL    *url.URL
	SizeLimit  int64
	Engine     Engine
}

func (b *Bucket) GrantCompleteOptions(id string, object string, opt *Options) (*CompleteOptions, error) {
	p := urlencodesign.NewParams()
	expired := time.Now().Add(opt.Lifetime).Unix()
	ts := strconv.FormatInt(expired, 10)
	p.Append(app.Sign.AppidField, opt.Appid)
	p.Append(app.Sign.BucketField, b.Name)
	p.Append(app.Sign.ObjectField, object)
	p.Append(app.Sign.IDField, id)
	p.Append(app.Sign.TimestampField, ts)
	s, err := urlencodesign.Sign(hasher.Md5Hasher, secret.Secret(opt.Secret), app.Sign.SecretField, p, true)
	if err != nil {
		return nil, err
	}
	o := NewCompleteOptions()
	o.ID = id
	o.Bucket = b.Name
	o.Object = object
	o.Timestamp = ts
	o.Sign = s
	return o, nil
}
func (b *Bucket) Verify() error {
	return nil
}
func (b *Bucket) Join(ele ...string) *url.URL {
	u := *b.BaseURL
	u.Path = path.Join(u.Path, path.Join(ele...))
	return &u
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
	b.Offset = time.Duration(config.OffsetInSeconds) * time.Second
	if b.Offset <= 0 {
		b.Offset = time.Duration(app.System.OffsetInSeconds) * time.Second
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
	Complete(b *Bucket, id string, objectname string, opt *Options) (info *CompleteInfo, err error)
	RemoveFile(b *Bucket, objectname string) error
	Permanent() bool
	ThirdpartyUpload() bool
	ThirdpartyDownload() bool
	BucketType() string
	ServeHTTPDownload(w http.ResponseWriter, r *http.Request)
	Download(b *Bucket, objectname string, w io.Writer) (err error)
	Upload(b *Bucket, objectname string, body io.Reader) (err error)
	GetFileinfo(b *Bucket, objectname string) (info *Fileinfo, err error)
	Start() error
	Stop() error
}
