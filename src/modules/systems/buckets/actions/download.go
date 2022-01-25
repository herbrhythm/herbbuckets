package actions

import (
	"herbbuckets/modules/bucket"
	"net/http"

	"github.com/herb-go/herb/middleware/action"
)

var ActionDownload = action.New(func(w http.ResponseWriter, r *http.Request) {
	bu := bucket.GetBucketFromRequest(r)
	if bu == nil || !bu.Engine.ThirdpartyDownload() {
		http.NotFound(w, r)
		return
	}
	bu.Engine.ServeHTTPDownload(w, r)
})
