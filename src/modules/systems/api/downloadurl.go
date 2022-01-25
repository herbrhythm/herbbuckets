package api

import (
	"herbbuckets/modules/bucket"
	"herbbuckets/modules/pathcleaner"
	"herbbuckets/modules/systems/buckets"
	"net/http"
	"strconv"
	"time"

	"github.com/herb-go/herb/middleware/action"
	"github.com/herb-go/herb/middleware/router"
	"github.com/herb-go/herb/ui/render"
	"github.com/herb-go/herbmodules/protecter"
	"github.com/herb-go/herbsecurity/authority"
)

type DownloadURLResult struct {
	URL string
}

var DownloadURLAction = action.New(func(w http.ResponseWriter, r *http.Request) {
	var err error
	var expired int64
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.NotFound(w, r)
		return
	}
	filename = pathcleaner.Clean(filename)
	bucketname := router.GetParams(r).Get("bucket")
	if bucketname == "" {
		http.NotFound(w, r)
		return
	}
	b := buckets.Buckets[bucketname]
	if b == nil {
		http.NotFound(w, r)
		return
	}
	ts := r.URL.Query().Get("timestamp")
	if ts == "" {
		expired = time.Now().Unix()
	} else {
		expired, err = strconv.ParseInt(ts, 10, 64)
		if err != nil {
			http.NotFound(w, r)
			return
		}
	}
	auth := protecter.LoadAuth(r)

	opt := &bucket.Options{}
	opt.Appid = auth.Authority().String()
	opt.Secret = auth.Payloads().LoadString(authority.PayloadSignSecret)
	opt.ExpiredAt = expired
	downloadurl, err := b.GetDownloadURL(bucketname, filename, opt)
	if err != nil {
		if err == bucket.ErrNotFound {
			http.NotFound(w, r)
			return
		}
		panic(err)
	}
	render.MustJSON(w, &DownloadURLResult{URL: downloadurl}, 200)
})
