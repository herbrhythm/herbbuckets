package localbucket

import (
	"fmt"
	"herbbuckets/modules/app"
	"herbbuckets/modules/bucket"
	bucketsmiddlewares "herbbuckets/modules/systems/buckets/middlewares"
	"io"
	"net/http"
	"net/url"
	"os"
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

const ExtentionsSeparator = ","

type Config struct {
	Public   bool
	Hasher   string
	Location string
}

func (c *Config) ApplyTo(bu *bucket.Bucket, b *LocalBucket) error {
	b.Public = c.Public
	b.Hasher = c.Hasher
	b.Cors = bu.Cors
	b.Location = c.Location
	if b.Location == "" {
		b.Location = util.AppData(bucket.BucketsFolder, bu.Name)
	}
	if bu.BaseURL == nil {
		return fmt.Errorf("bucket [%s] BaseURL format error", bu.Name)
	}
	return nil
}

type LocalBucket struct {
	Public    bool
	Location  string
	Hasher    string
	Temporary bool
	Cors      *cors.CORS
}

func (b *LocalBucket) localpath(object string) string {
	return filepath.Join(b.Location, object)
}
func (b *LocalBucket) newDownloadInfo() *bucket.DownloadInfo {
	info := bucket.NewDownloadInfo()
	return info
}

func (b *LocalBucket) GrantDownloadInfo(bu *bucket.Bucket, object string, opt *bucket.Options) (info *bucket.DownloadInfo, err error) {
	info = b.newDownloadInfo()
	if b.Public {
		info.Permanent = true
		info.URL = bu.Join(bucket.PrefixDownload, bu.Name, object)
		return info, nil
	}
	p := urlencodesign.NewParams()
	p.Append(app.Sign.ObjectField, object)
	p.Append(app.Sign.AppidField, opt.Appid)
	expired := time.Now().Add(opt.Lifetime).Unix()
	ts := strconv.FormatInt(expired, 10)
	p.Append(app.Sign.TimestampField, ts)
	s, err := urlencodesign.Sign(hasher.Md5Hasher, secret.Secret(opt.Secret), app.Sign.SecretField, p, true)
	if err != nil {
		return nil, err
	}
	q := &url.Values{}
	q.Add(app.Sign.AppidField, opt.Appid)
	q.Add(app.Sign.TimestampField, ts)
	q.Add(app.Sign.SignField, s)
	info.URL = bu.Join(bucket.PrefixDownload, bu.Name, object) + "?" + q.Encode()
	info.Permanent = false
	info.ExpiredAt = expired
	return info, nil
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
func (b *LocalBucket) GrantUploadInfo(bu *bucket.Bucket, id string, object string, opt *bucket.Options) (info *bucket.WebuploadInfo, err error) {
	expired := time.Now().Add(opt.Lifetime).Unix()
	ts := strconv.FormatInt(expired, 10)
	sizelimit := strconv.FormatInt(opt.SizeLimit, 10)
	p := urlencodesign.NewParams()
	p.Append(app.Sign.AppidField, opt.Appid)
	p.Append(app.Sign.TimestampField, ts)
	p.Append(app.Sign.BucketField, bu.Name)
	p.Append(app.Sign.ObjectField, object)
	p.Append(app.Sign.SizeLimitField, sizelimit)
	s, err := urlencodesign.Sign(hasher.Md5Hasher, secret.Secret(opt.Secret), app.Sign.SecretField, p, true)
	if err != nil {
		return nil, err
	}
	downloadinfo, err := b.GrantDownloadInfo(bu, object, opt)
	if err != nil {
		return nil, err
	}
	q := &url.Values{}
	q.Add(app.Sign.AppidField, opt.Appid)
	q.Add(app.Sign.SignField, s)
	q.Add(app.Sign.TimestampField, ts)
	q.Add(app.Sign.BucketField, bu.Name)
	q.Add(app.Sign.SizeLimitField, sizelimit)
	info = b.newWebuploadInfo()
	info.UploadURL = bu.Join(bucket.PrefixUpload, bu.Name, object) + "?" + q.Encode()
	info.PreviewURL = downloadinfo.URL
	info.Bucket = bu.Name
	info.ID = id
	info.Object = object
	info.SizeLimit = opt.SizeLimit
	info.ExpiredAt = expired
	return info, nil
}
func (b *LocalBucket) Download(bu *bucket.Bucket, objectname string) (r io.ReadCloser, err error) {
	f, err := os.Open(b.localpath(objectname))
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
	info.Modtime = stat.ModTime().Unix()
	return info, nil
}

func (b *LocalBucket) RemoveFile(bu *bucket.Bucket, objectname string) error {
	err := os.Remove(b.localpath(objectname))
	if err != nil {
		if os.IsNotExist(err) {
			return bucket.ErrNotFound
		}
		return err
	}
	return nil
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
