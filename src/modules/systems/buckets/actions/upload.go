package actions

import (
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
	formerr := &validator.Validator{}
	bu := bucket.GetBucketFromRequest(r)
	q := r.URL.Query()
	file, _, err := r.FormFile(bucket.PostFieldFile)
	if err != nil {
		if err == http.ErrMissingFile || err == http.ErrNotMultipart {
			formerr.AddPlainError(bucket.PostFieldFile, "File required")
			render.MustJSON(w, formerr.Errors(), 422)
			return
		}
		panic(err)
	}
	defer file.Close()
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
	writer, err := bu.Engine.Upload(bu, objectname)
	if err != nil {
		if err == bucket.ErrExists {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		panic(err)
	}
	defer writer.Close()
	_, err = io.Copy(writer, file)
	if err != nil {
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
	formerr := &validator.Validator{}
	bu := bucket.GetBucketFromRequest(r)
	if bu == nil || bu.Engine.ThirdpartyUpload() {
		http.NotFound(w, r)
		return
	}
	q := r.URL.Query()
	objectname := httprouter.GetParams(r).Get(bucket.RouterParamObject)
	file, h, err := r.FormFile(bucket.PostFieldFile)
	if err != nil {
		if err == http.ErrMissingFile || err == http.ErrNotMultipart {
			formerr.AddPlainError(bucket.PostFieldFile, "File required")
			render.MustJSON(w, formerr.Errors(), 422)
			return
		}
		panic(err)
	}
	defer file.Close()
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
	if sizelimit > 0 {
		if h.Size > sizelimit {
			formerr.AddPlainError(bucket.PostFieldFile, "File too large")
			render.MustJSON(w, formerr.Errors(), 422)
			return
		}
	}
	writer, err := bu.Engine.Upload(bu, objectname)
	if err != nil {
		if err == bucket.ErrExists {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		panic(err)
	}
	defer writer.Close()
	_, err = io.Copy(writer, file)
	if err != nil {
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
