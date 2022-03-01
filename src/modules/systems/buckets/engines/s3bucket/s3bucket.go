package s3bucket

import (
	"context"
	"herbbuckets/modules/bucket"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/herb-go/providers/s3compatible"
)

const BucketType = "s3"

type Config struct {
	Public bool
	Bucket string
	s3compatible.S3Config
}

func (b *S3Bucket) ConvertError(err error) error {
	if s3compatible.IsHTTPError(err, 404) {
		return bucket.ErrNotFound
	}
	return err
}
func (c *Config) ApplyTo(bu *bucket.Bucket, b *S3Bucket) error {
	api := s3compatible.New()
	err := c.S3Config.ApplyTo(api)
	if err != nil {
		return err
	}
	b.API = api
	b.Public = c.Public
	b.Bucket = c.Bucket
	return nil
}

type S3Bucket struct {
	Public bool
	Bucket string
	API    *s3compatible.API
}

func (b *S3Bucket) newWebuploadInfo() *bucket.WebuploadInfo {
	info := bucket.NewWebuploadInfo()
	info.SuccessCodeMin = 200
	info.SuccessCodeMax = 299
	info.FileField = ""
	info.UploadType = "put"
	return info
}
func (b *S3Bucket) GrantUploadInfo(bu *bucket.Bucket, id string, object string, opt *bucket.Options) (info *bucket.WebuploadInfo, err error) {
	info = b.newWebuploadInfo()
	o := s3compatible.NewUploadOptions()
	expired := time.Now().Add(opt.Lifetime).Unix()
	url, err := b.API.PresignPutObject(context.TODO(), b.Bucket, object, opt.Lifetime, o)
	if err != nil {
		return nil, err
	}
	info.UploadURL = url
	co, err := bu.GrantCompleteOptions(id, object, opt)
	if err != nil {
		return nil, err
	}
	info.Complete = co.Encode()
	info.Bucket = bu.Name
	info.ID = id
	info.Object = object
	info.SizeLimit = opt.SizeLimit
	info.ExpiredAt = expired
	return info, nil
}
func (b *S3Bucket) newDownloadInfo() *bucket.DownloadInfo {
	info := bucket.NewDownloadInfo()
	return info
}
func (b *S3Bucket) GrantDownloadInfo(bu *bucket.Bucket, object string, opt *bucket.Options) (info *bucket.DownloadInfo, err error) {
	info = b.newDownloadInfo()
	surl, err := b.API.PresignGetObject(context.TODO(), b.Bucket, object, opt.Lifetime)
	if err != nil {
		return nil, err
	}
	info.Permanent = b.Public
	if b.Public {
		u, err := url.Parse(surl)
		if err != nil {
			return nil, err
		}
		u.RawQuery = ""
		info.URL = u.String()
	} else {
		info.ExpiredAt = time.Now().Add(opt.Lifetime).Unix()
		info.URL = surl
	}
	return info, nil
}
func (b *S3Bucket) Permanent() bool {
	return b.Public
}
func (b *S3Bucket) Download(bu *bucket.Bucket, objectname string, w io.Writer) (err error) {
	_, err = b.API.Load(context.TODO(), b.Bucket, objectname, w)
	return b.ConvertError(err)
}
func (b *S3Bucket) Upload(bu *bucket.Bucket, objectname string, body io.Reader) (err error) {
	err = b.API.Save(context.TODO(), b.Bucket, objectname, body)
	return b.ConvertError(err)
}
func (b *S3Bucket) ServeHTTPDownload(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}
func (b *S3Bucket) GetFileinfo(bu *bucket.Bucket, objectname string) (info *bucket.Fileinfo, err error) {
	result, err := b.API.Info(context.TODO(), b.Bucket, objectname)
	if err != nil {
		return nil, b.ConvertError(err)
	}
	info = bucket.NewFileinfo()
	info.Modtime = result.LastModified.Unix()
	info.Size = result.Size
	return info, nil
}
func (b *S3Bucket) RemoveFile(bu *bucket.Bucket, objectname string) error {
	return b.ConvertError(b.API.Remove(context.TODO(), b.Bucket, objectname))
}
func (b *S3Bucket) newCompleteInfo() *bucket.CompleteInfo {
	info := bucket.NewCompleteInfo()
	return info
}
func (b *S3Bucket) Complete(bu *bucket.Bucket, id string, objectname string, opt *bucket.Options) (info *bucket.CompleteInfo, err error) {
	fi, err := b.GetFileinfo(bu, objectname)
	if err != nil {
		return nil, err
	}
	di, err := b.GrantDownloadInfo(bu, objectname, opt)
	if err != nil {
		return nil, err
	}
	info = b.newCompleteInfo()
	info.ID = id
	info.Bucket = bu.Name
	info.Object = objectname
	info.Size = fi.Size
	info.Preview = di
	return info, nil
}
func (b *S3Bucket) ThirdpartyUpload() bool {
	return true
}
func (b *S3Bucket) ThirdpartyDownload() bool {
	return true
}
func (b *S3Bucket) BucketType() string {
	return BucketType
}
func (b *S3Bucket) Start() error {
	return nil
}
func (b *S3Bucket) Stop() error {
	return nil
}
func New() *S3Bucket {
	return &S3Bucket{}
}
func Builder(b *bucket.Bucket, loader func(v interface{}) error) error {
	s3b := New()
	config := &Config{}
	err := loader(config)
	if err != nil {
		return err
	}
	err = config.ApplyTo(b, s3b)
	if err != nil {
		return err
	}
	b.Engine = s3b
	return nil
}

func init() {
	bucket.Builders[BucketType] = Builder
}
