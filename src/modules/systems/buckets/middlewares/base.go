package middlewares

import (
	"herbbuckets/modules/bucket"
	"herbbuckets/modules/systems/buckets"
	"net/http"
	"path"

	"github.com/herb-go/herb/middleware/router/httprouter"
)

var MiddlewareBase = func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ps := httprouter.GetParams(r)
	bucketname := path.Base(r.URL.Path)
	if bucketname == "" || bucketname == "/" || bucketname == "." {
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
	next(w, r)
}
