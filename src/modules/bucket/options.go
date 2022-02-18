package bucket

import (
	"herbbuckets/modules/app"
	"net/http"
	"net/url"
	"time"
)

type Options struct {
	Appid     string
	Secret    string
	SizeLimit int64
	Lifetime  time.Duration
}

func NewOptions() *Options {
	return &Options{}
}

type CompleteOptions struct {
	ID        string
	Bucket    string
	Object    string
	Timestamp string
	Sign      string
}

func (o *CompleteOptions) Encode() string {
	params := &url.Values{}
	params.Add(app.Sign.IDField, o.ID)
	params.Add(app.Sign.BucketField, o.Bucket)
	params.Add(app.Sign.ObjectField, o.Object)
	params.Add(app.Sign.TimestampField, o.Timestamp)
	params.Add(app.Sign.SignField, o.Sign)
	return params.Encode()
}
func (o *CompleteOptions) Decode(r *http.Request) {
	q := r.URL.Query()
	o.ID = q.Get(app.Sign.IDField)
	o.Bucket = q.Get(app.Sign.BucketField)
	o.Object = q.Get(app.Sign.ObjectField)
	o.Timestamp = q.Get(app.Sign.TimestampField)
	o.Sign = q.Get(app.Sign.SignField)
}
func NewCompleteOptions() *CompleteOptions {
	return &CompleteOptions{}
}
