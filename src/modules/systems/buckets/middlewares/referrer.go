package middlewares

import (
	"herbbuckets/modules/bucket"
	"net/http"
)

var MiddlewareReferrer = func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	b := bucket.GetBucketFromRequest(r)
	if b == nil {
		http.NotFound(w, r)
		return
	}
}
