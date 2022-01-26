package middlewares

import (
	"herbbuckets/modules/bucket"
	"net/http"

	"github.com/herb-go/herbmodules/protecter"
)

var MiddlewareSignDownload = func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	protecter.ProtectMiddleware(bucket.ProtecterDownload)(w, r, func(w http.ResponseWriter, r *http.Request) {
		MiddlewareAuthViewBucket(w, r, next)
	})
}
