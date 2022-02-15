package actions

import (
	"herbbuckets/modules/bucket"
	"net/http"

	"github.com/herb-go/herb/middleware/action"
	"github.com/herb-go/herb/middleware/router/httprouter"
	"github.com/herb-go/herb/ui/render"
)

var ActionRemove = action.New(func(w http.ResponseWriter, r *http.Request) {
	bu := bucket.GetBucketFromRequest(r)
	objectname := httprouter.GetParams(r).Get(bucket.RouterParamObject)
	err := bu.Engine.RemoveFile(bu, objectname)
	if err != nil {
		if err == bucket.ErrNotFound {
			http.NotFound(w, r)
			return
		}
		panic(err)
	}
	render.MustJSON(w, "success", 200)
})
