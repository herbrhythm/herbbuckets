package middlewares

import (
	"herbbuckets/modules/bucket"
	"net/http"
)

var MiddlewareWebuploadCORS = func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	bu := bucket.GetBucketFromRequest(r)
	switch r.Method {
	case "OPTIONS":
		bu.Cors.Preflight(w, r)
		return
	case "POST":
		bu.Cors.ServeMiddleware(w, r, next)
		return
	}
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
