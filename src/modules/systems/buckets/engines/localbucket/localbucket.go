package localbucket

import (
	"herbbuckets/modules/app"
	"herbbuckets/modules/bucket"
	bucketsmiddlewares "herbbuckets/modules/systems/buckets/middlewares"
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

const DefaultDataformat = "2006/01"

const ExtentionsSeparator = ","

type Config struct {
	Public   bool
	Hasher   string
	Location string
	Cors     cors.CORS
}

func (c *Config) ApplyTo(bu *bucket.Bucket, b *LocalBucket) error {
	b.Public = c.Public
	b.Hasher = c.Hasher
	b.Cors = &c.Cors
	b.Location = c.Location
	if b.Location == "" {
		b.Location = util.AppData(bucket.BucketsFolder, bu.Name)
	}
	return nil
}

type LocalBucket struct {
	Public   bool
	Location string
	Hasher   string
	Cors     *cors.CORS
}

func (b *LocalBucket) localpath(object string) string {
	return filepath.Join(b.Location, object)
}
func (b *LocalBucket) GrantDownloadURL(bu *bucket.Bucket, object string, opt *bucket.Options) (downloadurl string, err error) {
	if b.Public {
		return path.Join(bu.BaseURL, bucket.PrefixDownload, bu.Name, object), nil
	}
	p := urlencodesign.NewParams()
	p.Append(app.Sign.ObjectField, object)
	p.Append(app.Sign.AppidField, opt.Appid)
	ts := strconv.FormatInt(time.Now().Add(opt.Lifetime).Unix(), 10)
	p.Append(app.Sign.TimestampField, ts)
	s, err := urlencodesign.Sign(hasher.Md5Hasher, secret.Secret(opt.Secret), app.Sign.SecretField, p, true)
	if err != nil {
		return "", err
	}
	q := &url.Values{}
	q.Add(app.Sign.AppidField, opt.Appid)
	q.Add(app.Sign.TimestampField, ts)
	q.Add(app.Sign.SignField, s)
	return path.Join(bu.BaseURL, bucket.PrefixDownload, bu.Name, object) + "?" + q.Encode(), nil
}
func (b *LocalBucket) Permanent() bool {
	return b.Public
}
func (b *LocalBucket) newWebuploadInfo() *bucket.WebuploadInfo {
	info := bucket.NewWebuploadInfo()
	info.SuccessCodeMin = 200
	info.SuccessCodeMax = 299
	info.Permanent = b.Public
	info.FileField = bucket.PostFieldFile
	return info
}
func (b *LocalBucket) GrantUploadInfo(bu *bucket.Bucket, object string, opt *bucket.Options) (info *bucket.WebuploadInfo, err error) {
	ts := strconv.FormatInt(time.Now().Add(opt.Lifetime).Unix(), 10)
	sizelimit := strconv.FormatInt(opt.Sizelimit, 10)
	p := urlencodesign.NewParams()
	p.Append(app.Sign.AppidField, opt.Appid)
	p.Append(app.Sign.TimestampField, ts)
	p.Append(app.Sign.BucketField, bu.Name)
	p.Append(app.Sign.ObjectField, object)
	p.Append(app.Sign.SizelimitField, sizelimit)
	s, err := urlencodesign.Sign(hasher.Md5Hasher, secret.Secret(opt.Secret), app.Sign.SecretField, p, true)
	if err != nil {
		return nil, err
	}
	previewurl, err := b.GrantDownloadURL(bu, object, opt)
	if err != nil {
		return nil, err
	}
	q := &url.Values{}
	q.Add(app.Sign.AppidField, opt.Appid)
	q.Add(app.Sign.SignField, s)
	q.Add(app.Sign.TimestampField, ts)
	q.Add(app.Sign.BucketField, bu.Name)
	q.Add(app.Sign.SizelimitField, sizelimit)
	info = b.newWebuploadInfo()
	info.UploadURL = path.Join(bu.BaseURL+bucket.PrefixUpload, bu.Name, object) + "?" + q.Encode()
	info.PreviewURL = previewurl
	info.Bucket = bu.Name
	info.Objcet = object
	info.Sizelimit = opt.Sizelimit
	return info, nil
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
	lp := b.localpath(objectname)
	folder := filepath.Dir(lp)
	_, err = os.Stat(folder)
	if err == nil {
		return nil, bucket.ErrExists
	}
	if !os.IsNotExist(err) {
		return nil, err
	}
	os.MkdirAll(folder, util.DefaultFolderMode)
	return os.Create(lp)
}
func (b *LocalBucket) serveHTTPDownload(w http.ResponseWriter, r *http.Request) {
	objectname := httprouter.GetParams(r).Get(bucket.RouterParamObject)
	simplehttpserver.ServeFile(b.localpath(objectname)).ServeHTTP(w, r)
}
func (b *LocalBucket) ServeHTTPDownload(w http.ResponseWriter, r *http.Request) {
	if !b.Public {
		bucketsmiddlewares.MiddlewareSignDownload(w, r, b.serveHTTPDownload)
		return
	}
	b.serveHTTPDownload(w, r)
}
func (b *LocalBucket) GetFileinfo(bu *bucket.Bucket, objectname string) (info *bucket.Fileinfo, err error) {
	stat, err := os.Stat(b.localpath(objectname))
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

func (b *LocalBucket) RemoveFile(bu *bucket.Bucket, objectname string) error {
	return os.Remove(b.localpath(objectname))
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
func (b *LocalBucket) Start() error {
	return nil
}
func (b *LocalBucket) Stop() error {
	return nil
}
func New() *LocalBucket {
	return &LocalBucket{}
}
func Builder(b *bucket.Bucket, loader func(v interface{}) error) error {
	lb := New()
	config := &Config{}
	err := loader(config)
	if err != nil {
		return err
	}
	err = config.ApplyTo(b, lb)
	if err != nil {
		return err
	}
	b.Engine = lb
	return nil
}

func init() {
	bucket.Builders[""] = Builder
	bucket.Builders[BucketType] = Builder
}
