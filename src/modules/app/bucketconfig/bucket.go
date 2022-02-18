package bucketconfig

import (
	"github.com/herb-go/herb/middleware/cors"
)

type Bucket struct {
	Name              string
	Type              string
	Disabled          bool
	DateFormat        string
	LifetimeInSeconds int64
	OffsetInSeconds   int64
	Cors              cors.CORS
	Referrer          []string
	SizeLimit         int64
	BaseURL           string
	Prefix            string
}
