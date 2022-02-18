package s3bucket

import (
	"herbbuckets/modules/bucket"
	"io"
	"net/http"

	"github.com/herb-go/fetcher"
)

const BucketType = "s3"
const ForbiddenFilenamePart = "${filename}"

type Config struct {
	Public           bool
	Server           *fetcher.Server
	FieldAccessKeyId string
}

func (c *Config) ApplyTo(bu *bucket.Bucket, b *S3Bucket) error {
	var err error
	b.Public = c.Public
	b.Preset, err = c.Server.CreatePreset()
	if err != nil {
		return err
	}
	return nil
}

type S3Bucket struct {
	Public           bool
	Preset           *fetcher.Preset
	FieldAccessKeyId string
}

func (b *S3Bucket) newWebuploadInfo() *bucket.WebuploadInfo {
	info := bucket.NewWebuploadInfo()
	info.SuccessCodeMin = 200
	info.SuccessCodeMax = 299
	info.FileField = bucket.PostFieldFile
	return info
}
func (b *S3Bucket) GrantUploadInfo(bu *bucket.Bucket, id string, object string, opt *bucket.Options) (info *bucket.WebuploadInfo, err error) {
	return nil, nil
}
func (b *S3Bucket) newDownloadInfo() *bucket.DownloadInfo {
	info := bucket.NewDownloadInfo()
	return info
}
func (b *S3Bucket) GrantDownloadInfo(bu *bucket.Bucket, object string, opt *bucket.Options) (info *bucket.DownloadInfo, err error) {
	return nil, nil
}
func (b *S3Bucket) Permanent() bool {
	return b.Public
}
func (b *S3Bucket) Download(bu *bucket.Bucket, objectname string) (r io.ReadCloser, err error) {
	return nil, nil
}
func (b *S3Bucket) Upload(bu *bucket.Bucket, objectname string) (w io.WriteCloser, err error) {
	return nil, nil
}
func (b *S3Bucket) ServeHTTPDownload(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}
func (b *S3Bucket) GetFileinfo(bu *bucket.Bucket, objectname string) (info *bucket.Fileinfo, err error) {
	return nil, nil
}
func (b *S3Bucket) RemoveFile(bu *bucket.Bucket, objectname string) error {
	return nil
}
func (b *S3Bucket) Complete(bu *bucket.Bucket, id string, objectname string, opt *bucket.Options) (info *bucket.CompleteInfo, err error) {
	return nil, nil
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
