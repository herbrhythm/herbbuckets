package bucket

import (
	"context"
	"herbbuckets/modules/protecters"
	"net/http"
)

var ProtecterDownload = protecters.Get("appkeysecret")

const BucketsFolder = "buckets"

const PrefixDownload = "/file"

const PrefixUpload = "/upload"

const QueryFieldBucket = "bucket"
const QueryFieldObject = "object"

const RouterParamBucket = "bucket"
const RouterParamObject = "object"

type ContextKey string

const ContextKeyBucket = ContextKey("bucket")

func GetBucketFromRequest(r *http.Request) *Bucket {
	v := r.Context().Value(ContextKeyBucket)
	if v == nil {
		return nil
	}
	b, ok := v.(*Bucket)
	if !ok {
		return nil
	}
	return b
}

func SaveBucketToRequest(req **http.Request, b *Bucket) {
	ctx := context.WithValue((*req).Context(), ContextKeyBucket, b)
	*req = (*req).WithContext(ctx)
}
