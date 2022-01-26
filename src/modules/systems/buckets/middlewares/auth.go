package middlewares

import (
	"net/http"

	"github.com/herb-go/herbmodules/httpinfomanager"
	au "github.com/herb-go/herbmodules/protecter/payloadauthorizer"
	"github.com/herb-go/util"
)

const PolicyRoot = "root"
const PolicyViewAll = "viewall"
const PolicyViewBucket = "view:bucket={{0}}"
const PolicyViewUpload = "upload"
const PolicyViewUploadAll = "upload:bucket={{0}}"

var MiddlewareAuthViewBucket func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
var MiddlewareAuthUploadBucket func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)

func init() {
	util.RegisterInitiator(ModuleName, "initmiddleware", func() {
		MiddlewareAuthViewBucket = au.New().WithAny(au.MustParse(PolicyRoot), au.MustParse(PolicyViewAll), au.MustParseWith(PolicyViewBucket, httpinfomanager.MustGetStringField("bucket").LoadStringInfo)).ServeMiddleware
		MiddlewareAuthUploadBucket = au.New().WithAny(au.MustParse(PolicyRoot), au.MustParse(PolicyViewUploadAll), au.MustParseWith(PolicyViewUpload, httpinfomanager.MustGetStringField("bucket").LoadStringInfo)).ServeMiddleware
	})
}
