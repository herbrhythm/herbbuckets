package middleware

import (
	"herbbuckets/modules/bucket"
	"herbbuckets/modules/systems/buckets"
	"net/http"

	"github.com/herb-go/herb/middleware/router/httprouter"
)

var MiddlewareReferrer = func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	bucketname := httprouter.GetParams(r).Get(bucket.RouterParamBucket)
	if bucketname == "" {
		http.NotFound(w, r)
		return
	}
	b := buckets.Buckets[bucketname]
	if b == nil {
		http.NotFound(w, r)
		return
	}
}
