package localbucket

import (
	"herbbuckets/modules/app"
	"herbbuckets/modules/bucket"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/herb-go/herb/file/simplehttpserver"
	"github.com/herb-go/herb/middleware/cors"
	"github.com/herb-go/herb/middleware/router/httprouter"
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
	Public            bool
	LifetimeInSeconds int64
	Hasher            string
	Location          string
	Cors              cors.CORS
}

func (c *Config) ApplyTo(bu *bucket.Bucket, b *LocalBucket) error {
	b.Public = c.Public
	b.Hasher = c.Hasher
	b.Lifetime = time.Duration(time.Second * time.Duration(c.LifetimeInSeconds))
	b.Cors = &c.Cors
	b.Location = c.Location
	if b.Location == "" {
		b.Location = util.AppData(BucketFolder, bu.Name)
	}
	return nil
}

type LocalBucket struct {
	Public   bool
	Location string
	Hasher   string
	Lifetime time.Duration
	Cors     *cors.CORS
}

func (b *LocalBucket) localpath(bucketname string, object string) string {
	return filepath.Join(b.Location, bucketname, object)
}
func (b *LocalBucket) GrantDownloadURL(bu *bucket.Bucket, object string, opt *bucket.Options) (downloadurl string, err error) {
	p := urlencodesign.NewParams()
	p.Append(app.Sign.ObjectField, object)
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
	return path.Join(bu.BaseURL, bu.Name, object) + "?" + q.Encode(), nil
}
func (b *LocalBucket) Permanent() bool {
	return b.Public
}

func (b *LocalBucket) GrantUploadURL(bu *bucket.Bucket, object string, opt *bucket.Options) (uploadurl string, err error) {
	ts := strconv.FormatInt(opt.ExpiredAt, 10)
	p := urlencodesign.NewParams()
	p.Append(app.Sign.AppidField, opt.Appid)
	p.Append(app.Sign.TimestampField, ts)
	p.Append(app.Sign.BucketField, bu.Name)
	p.Append(app.Sign.ObjectField, object)
	s, err := urlencodesign.Sign(hasher.Md5Hasher, secret.Secret(opt.Secret), app.Sign.SecretField, p, true)
	if err != nil {
		return "", err
	}
	q := &url.Values{}
	q.Add(app.Sign.AppidField, opt.Appid)
	q.Add(app.Sign.SignField, s)
	q.Add(app.Sign.TimestampField, ts)
	q.Add(app.Sign.BucketField, bu.Name)
	q.Add(app.Sign.ObjectField, object)
	return bu.BaseURL + UploadRouter + "?" + q.Encode(), nil
}
func (b *LocalBucket) Download(bu *bucket.Bucket, objectname string) (r io.ReadCloser, err error) {
	f, err := os.Open(filepath.Join(b.Location, bu.Name, objectname))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, bucket.ErrNotFound
		}
	}
	return f, nil
}
func (b *LocalBucket) Upload(bu *bucket.Bucket, objectname string) (w io.WriteCloser, err error) {
	lp := b.localpath(bu.Name, objectname)
	folder := filepath.Dir(lp)
	_, err = os.Stat(folder)
	if err == nil {
		return nil, bucket.ErrExists
	}
	if !os.IsNotExist(err) {
		return nil, err
	}
	return os.Open(lp)
}
func (b *LocalBucket) ServeHTTPDownload(w http.ResponseWriter, r *http.Request) {
	bu := bucket.GetBucketFromRequest(r)
	objectname := httprouter.GetParams(r).Get(bucket.RouterParamObject)
	simplehttpserver.ServeFile(filepath.Join(b.Location, bu.Name, objectname)).ServeHTTP(w, r)
}
func (b *LocalBucket) GetFileinfo(bu *bucket.Bucket, objectname string) (info *bucket.Fileinfo, err error) {
	stat, err := os.Stat(b.localpath(bu.Name, objectname))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, bucket.ErrNotFound
		}
		return nil, err
	}
	info = bucket.NewFileinfo()
	info.Size = stat.Size()
	return info, nil
}
func (b *LocalBucket) GetVerifier() *bucket.Verifier {
	return Verifier
}
func (b *LocalBucket) RemoveFile(bu *bucket.Bucket, objectname string) error {
	return os.Remove(b.localpath(bu.Name, objectname))
}
func (b *LocalBucket) ThirdpartyUpload() bool {
	return false
}
func (b *LocalBucket) ThirdpartyDownload() bool {
	return false
}
func (b *LocalBucket) BucketType() string {
	return BucketType
}
func New() *LocalBucket {
	return &LocalBucket{}
}
func Builder(b *bucket.Bucket, loader func(v interface{}) error) error {
	lb := New()
	b.Engine = lb
	return nil
}
