package middleware

import (
	"herbbuckets/modules/bucket"
	"herbbuckets/modules/pathcleaner"
	"herbbuckets/modules/systems/buckets"
	"net/http"

	"github.com/herb-go/herb/middleware/router/httprouter"
)

var MiddlewareQuery = func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ps := httprouter.GetParams(r)
	q := r.URL.Query()
	bucketname := pathcleaner.Clean(q.Get(bucket.QueryFieldBucket))
	objectname := pathcleaner.Clean(q.Get(bucket.QueryFieldObject))
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
