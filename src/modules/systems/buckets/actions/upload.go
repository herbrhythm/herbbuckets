package actions

import (
	"herbbuckets/modules/app"
	"herbbuckets/modules/bucket"
	"herbbuckets/modules/pathcleaner"
	"herbbuckets/modules/uniqueid"
	"io"
	"net/http"
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

type SaveResult struct {
	ID     string
	Bucket string
	Object string
}

var ActionSave = action.New(func(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	formerr := &validator.Validator{}
	bu := bucket.GetBucketFromRequest(r)
	q := r.URL.Query()
	id := uniqueid.MustGenerateID()
	filename := q.Get(bucket.QueryFieldFilename)
	if filename == "" {
		formerr.AddPlainError(bucket.QueryFieldFilename, "Filename required")
		render.MustJSON(w, formerr.Errors(), 422)
		return
	}
	if path.Base(filename) != filename {
		formerr.AddPlainError(bucket.QueryFieldFilename, "Filename format error")
		render.MustJSON(w, formerr.Errors(), 422)
		return
	}
	objectname := pathcleaner.CreateObjectID(bu, id, filename)
	err := bu.Engine.Upload(bu, objectname, r.Body)
	if err != nil {
		if err == bucket.ErrExists {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		panic(err)
	}
	result := SaveResult{
		ID:     id,
		Bucket: bu.Name,
		Object: objectname,
	}
	render.MustJSON(w, result, 200)
})
var ActionUpload = action.New(func(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	formerr := &validator.Validator{}
	bu := bucket.GetBucketFromRequest(r)
	if bu == nil || bu.Engine.ThirdpartyUpload() {
		http.NotFound(w, r)
		return
	}
	q := r.URL.Query()
	objectname := httprouter.GetParams(r).Get(bucket.RouterParamObject)
	sizelimitq := q.Get(bucket.QueryFieldSizeLimit)
	if sizelimitq == "" {
		sizelimitq = "0"
	}
	sizelimit, err := strconv.ParseInt(sizelimitq, 10, 64)
	if err != nil {
		formerr.AddPlainError(bucket.QueryFieldSizeLimit, "Sizelimit format error")
		render.MustJSON(w, formerr.Errors(), 422)
		return
	}
	if bu.SizeLimit > 0 {
		if sizelimit > bu.SizeLimit {
			sizelimit = bu.SizeLimit
		}
	}
	var reader io.Reader = r.Body
	if sizelimit > 0 {
		reader = io.LimitReader(reader, sizelimit)
	}
	err = bu.Engine.Upload(bu, objectname, reader)
	if err != nil {
		if err == bucket.ErrExists {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		panic(err)
	}
	render.MustJSON(w, "success", 200)
})
var ActionGrantUploadInfo = action.New(func(w http.ResponseWriter, r *http.Request) {
	formerr := &validator.Validator{}
	bu := bucket.GetBucketFromRequest(r)
	q := r.URL.Query()
	filename := q.Get(bucket.QueryFieldFilename)
	if filename == "" {
		formerr.AddPlainError(bucket.QueryFieldFilename, "Filename required")
		render.MustJSON(w, formerr.Errors(), 422)
		return
	}
	if path.Base(filename) != filename {
		formerr.AddPlainError(bucket.QueryFieldFilename, "Filename format error")
		render.MustJSON(w, formerr.Errors(), 422)
		return
	}
	ttl := q.Get(bucket.QueryFieldTTL)
	if ttl == "" {
		ttl = "0"
	}
	ttli, err := strconv.ParseInt(ttl, 10, 64)
	if err != nil {
		formerr.AddPlainError(bucket.QueryFieldTTL, "TTL format error")
		render.MustJSON(w, formerr.Errors(), 422)
		return
	}
	lifetime := time.Duration(ttli) * time.Second
	if lifetime <= 0 {
		lifetime = bu.Lifetime
	}
	sizelimitq := q.Get(bucket.QueryFieldSizeLimit)
	if sizelimitq == "" {
		sizelimitq = "0"
	}
	sizelimit, err := strconv.ParseInt(sizelimitq, 10, 64)
	if err != nil {
		formerr.AddPlainError(bucket.QueryFieldSizeLimit, "Sizelimit format error")
		render.MustJSON(w, formerr.Errors(), 422)
		return
	}
	if sizelimit <= 0 {
		sizelimit = bu.SizeLimit
	}
	if (bu.SizeLimit > 0) && (sizelimit > bu.SizeLimit) {
		sizelimit = bu.SizeLimit
	}
	sizeq := q.Get(bucket.QueryFieldSize)
	size, err := strconv.ParseInt(sizeq, 10, 64)
	if err != nil {
		formerr.AddPlainError(bucket.QueryFieldSize, "size format error")
		render.MustJSON(w, formerr.Errors(), 422)
		return
	}
	if (sizelimit > 0) && size > sizelimit {
		formerr.AddPlainError(bucket.QueryFieldSize, "size too large")
		render.MustJSON(w, formerr.Errors(), 422)
		return
	}
	sizelimit = size
	id := uniqueid.MustGenerateID()
	object := pathcleaner.CreateObjectID(bu, id, filename)
	opt := bucket.NewOptions()
	auth := protecter.LoadAuth(r)
	opt.Appid = auth.Authority().String()
	opt.Lifetime = lifetime
	opt.SizeLimit = sizelimit
	opt.Secret = auth.Payloads().LoadString(authority.PayloadSignSecret)
	info, err := bu.Engine.GrantUploadInfo(bu, id, object, opt)
	if err != nil {
		panic(err)
	}
	render.MustJSON(w, info, 200)
})

var ActionComplete = action.New(func(w http.ResponseWriter, r *http.Request) {
	objectname := httprouter.GetParams(r).Get(bucket.RouterParamObject)
	bu := bucket.GetBucketFromRequest(r)
	auth := protecter.LoadAuth(r)
	opt := bucket.NewOptions()
	opt.Appid = auth.Authority().String()
	opt.Secret = auth.Payloads().LoadString(authority.PayloadSignSecret)
	opt.Lifetime = bu.Offset
	q := r.URL.Query()
	id := q.Get(app.Sign.IDField)
	info, err := bu.Engine.Complete(bu, id, objectname, opt)
	if err != nil {
		if err == bucket.ErrNotFound {
			http.NotFound(w, r)
			return
		}
		panic(err)
	}
	render.MustJSON(w, info, 200)
})
