package actions

import (
	"herbbuckets/modules/bucket"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/herb-go/herb/middleware/action"
	"github.com/herb-go/herb/middleware/router/httprouter"
	"github.com/herb-go/herb/ui/render"
	"github.com/herb-go/herb/ui/validator"
	"github.com/herb-go/herbmodules/protecter"
	"github.com/herb-go/herbsecurity/authority"
)

var ActionDownload = action.New(func(w http.ResponseWriter, r *http.Request) {
	bu := bucket.GetBucketFromRequest(r)
	if bu == nil || bu.Engine.ThirdpartyDownload() {
		http.NotFound(w, r)
		return
	}
	if len(bu.Referrer) != 0 {
		var matched bool
		referrer := r.Header.Get("referer")
		u, err := url.Parse(referrer)
		if err != nil {
			panic(err)
		}
		for _, v := range bu.Referrer {
			ok, err := path.Match(v, u.Host)
			if err != nil {
				panic(err)
			}
			if ok {
				matched = true
				break
			}
		}
		if !matched {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
	}
	bu.Engine.ServeHTTPDownload(w, r)
})

var ActionGrantDownloadURL = action.New(func(w http.ResponseWriter, r *http.Request) {
	formerr := &validator.Validator{}
	bu := bucket.GetBucketFromRequest(r)
	objectname := httprouter.GetParams(r).Get(bucket.RouterParamObject)
	q := r.URL.Query()
	ttl := q.Get("ttl")
	if ttl == "" {
		ttl = "0"
	}
	ttli, err := strconv.ParseInt(ttl, 10, 64)
	if err != nil {
		formerr.AddPlainError("ttl", "TTL format error")
		render.MustJSON(w, formerr.Errors(), 422)
		return
	}
	lifetime := time.Duration(ttli) * time.Second
	if lifetime < 0 {
		lifetime = bu.Lifetime
	}
	opt := bucket.NewOptions()
	auth := protecter.LoadAuth(r)
	opt.Appid = auth.Authority().String()
	opt.Lifetime = lifetime
	opt.Secret = auth.Payloads().LoadString(authority.PayloadSignSecret)
	u, err := bu.Engine.GrantDownloadURL(bu, objectname, opt)
	if err != nil {
		if err == bucket.ErrNotFound {
			http.NotFound(w, r)
			return
		}
		panic(err)
	}
	w.Write([]byte(u))
})
