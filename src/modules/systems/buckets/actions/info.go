package actions

import (
	"herbbuckets/modules/bucket"
	"net/http"

	"github.com/herb-go/herb/middleware/action"
	"github.com/herb-go/herb/middleware/router/httprouter"
	"github.com/herb-go/herb/ui/render"
	"github.com/herb-go/herbmodules/protecter"
	"github.com/herb-go/herbsecurity/authority"
)

var ActionInfo = action.New(func(w http.ResponseWriter, r *http.Request) {
	bu := bucket.GetBucketFromRequest(r)
	objectname := httprouter.GetParams(r).Get(bucket.RouterParamObject)
	info, err := bu.Engine.GetFileinfo(bu, objectname)
	if err != nil {
		if err == bucket.ErrNotFound {
			http.NotFound(w, r)
			return
		}
		panic(err)
	}
	render.MustJSON(w, info, 200)
})

type CurrentResult struct {
	AppID string
	Owner string
	Roles string
}

var ActionCurrent = action.New(func(w http.ResponseWriter, r *http.Request) {
	auth := protecter.LoadAuth(r)
	result := CurrentResult{
		AppID: auth.Authority().String(),
		Owner: auth.Principal().String(),
		Roles: auth.Payloads().LoadString(authority.PayloadRoles),
	}
	render.MustJSON(w, result, 200)
})
