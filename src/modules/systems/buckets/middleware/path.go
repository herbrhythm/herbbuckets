package middleware

import (
	"herbbuckets/modules/bucket"
	"herbbuckets/modules/pathcleaner"
	"herbbuckets/modules/systems/buckets"
	"net/http"
	"strings"

	"github.com/herb-go/herb/middleware/router/httprouter"
)

var MiddlewarePath = func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ps := httprouter.GetParams(r)
	p := strings.TrimPrefix(r.URL.Path, "/")
	list := strings.SplitN(p, "/", 2)
	if len(list) != 2 {
		http.NotFound(w, r)
		return
	}
	bucketname := pathcleaner.Clean(list[0])
	objectname := pathcleaner.Clean(list[1])
	if bucketname == "" || objectname == "" {
		http.NotFound(w, r)
		return
	}
	b := buckets.GetBucket(bucketname)
	if b == nil {
		http.NotFound(w, r)
		return
	}
	bucket.SaveBucketToRequest(&r, b)
	ps.Set(bucket.RouterParamBucket, bucketname)
	ps.Set(bucket.RouterParamObject, objectname)
	next(w, r)
}
