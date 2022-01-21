package localbucket

import (
	"herbbuckets/modules/app"
	"herbbuckets/modules/bucket"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/herb-go/herb/middleware/cors"
	"github.com/herb-go/herbsecurity/secret"
	"github.com/herb-go/herbsecurity/secret/hasher"
	"github.com/herb-go/herbsecurity/secret/hasher/urlencodesign"
	"github.com/herb-go/util"
)

const BucketType = "local"

const BucketFolder = "buckets"

const DefaultDataformat = "2006/01"

const UploadRouter = "/upload/"
const FileRouter = "/file/"

const ExtentionsSeparator = ","

var Verifier = &bucket.Verifier{
	CodeMin:  200,
	CodeMax:  200,
	Required: "success",
}

type Config struct {
	BaseURL           string
	Public            bool
	DataFormat        string
	Hasher            string
	Location          string
	LifetimeInSeconds int64
	Cors              cors.CORS
}

func (c *Config) ApplyTo(bucketname string, b LocalBucket) error {
	b.Public = c.Public
	b.Hasher = c.Hasher
	b.Lifetime = time.Duration(time.Second * time.Duration(c.LifetimeInSeconds))
	b.Cors = &c.Cors
	b.Location = c.Location
	if b.Location == "" {
		b.Location = util.AppData(BucketFolder, bucketname)
	}
	b.DateFormat = c.DataFormat
	burl := c.BaseURL
	if burl == "" {
		burl = app.HTTP.Config.BaseURL
	}
	u, err := url.Parse(burl)
	if err != nil {
		return err
	}
	b.BasePath = u.Scheme + "//" + u.Host + u.Path
	b.BasePath = u.Path
	return nil
}

type LocalBucket struct {
	BaseURL    string
	BasePath   string
	Public     bool
	Location   string
	DateFormat string
	Hasher     string
	Lifetime   time.Duration
	Cors       *cors.CORS
}

func (b *LocalBucket) GrantDownloadURL(bucketname string, object string, opt *bucket.Options) (downloadurl string, err error) {
	urlpath := path.Join(b.BasePath, bucketname, object)
	p := urlencodesign.NewParams()
	p.Append(app.Sign.PathField, urlpath)
	p.Append(app.Sign.AppidField, opt.Appid)
	ts := strconv.FormatInt(opt.ExpiredAt, 10)
	p.Append(app.Sign.TimestampField, ts)
	s, err := urlencodesign.Sign(hasher.Md5Hasher, secret.Secret(opt.Secret), app.Sign.SecretField, p, true)
	if err != nil {
		return "", err
	}
	q := &url.Values{}
	q.Add(app.Sign.AppidField, opt.Appid)
	q.Add(app.Sign.TimestampField, ts)
	q.Add(app.Sign.SignField, s)
	return b.BaseURL + urlpath + "?" + q.Encode(), nil

}
func (b *LocalBucket) Permanent() bool {
	return b.Public
}

func (b *LocalBucket) GrantUploadURL(bucketname string, object string, opt *bucket.Options) (uploadurl string, err error) {
	ts := strconv.FormatInt(opt.ExpiredAt, 10)
	p := urlencodesign.NewParams()
	p.Append(app.Sign.AppidField, opt.Appid)
	p.Append(app.Sign.TimestampField, ts)
	p.Append(app.Sign.BucketField, bucketname)
	p.Append(app.Sign.ObjectField, object)
	s, err := urlencodesign.Sign(hasher.Md5Hasher, secret.Secret(opt.Secret), app.Sign.SecretField, p, true)
	if err != nil {
		return "", err
	}
	q := &url.Values{}
	q.Add(app.Sign.AppidField, opt.Appid)
	q.Add(app.Sign.SignField, s)
	q.Add(app.Sign.TimestampField, ts)
	q.Add(app.Sign.BucketField, bucketname)
	q.Add(app.Sign.ObjectField, object)
	return b.BasePath + UploadRouter + "?" + q.Encode(), nil
}
func (b *LocalBucket) Download(bucketname string, objectname string, w io.Writer) (err error) {
	f, err := os.Open(filepath.Join(b.Location, bucketname, objectname))
	if err != nil {
		if os.IsNotExist(err) {
			return bucket.ErrNotFound
		}
	}
	defer f.Close()
	_, err = io.Copy(w, f)
	return err
}
func (b *LocalBucket) Upload(bucketname string, objectname string, r io.Reader) (err error) {
	f, err := os.Open(filepath.Join(b.Location, bucketname, objectname))
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, r)
	return err
}

func (b *LocalBucket) GetFileinfo(bucketname string, path string) (info *bucket.Fileinfo, err error) {
	return nil, nil
}
func (b *LocalBucket) GetVerifier() *bucket.Verifier {
	return Verifier
}
func (b *LocalBucket) RemoveFile(bucketname string, path string) error {
	return nil
}
func (b *LocalBucket) ThirdpartyUpload() bool {
	return false
}
func (b *LocalBucket) BucketType() string {
	return BucketType
}
func New() *LocalBucket {
	return &LocalBucket{}
}
func Factory(bucketname string, loader func(v interface{}) error) (bucket.Bucket, error) {
	b := New()
	return b, nil
}
