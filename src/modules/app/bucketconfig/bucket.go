package bucketconfig

import (
	"github.com/herb-go/herb-go/modules/cors"
)

type Bucket struct {
	Name              string
	Type              string
	Enabled           bool
	DateFormat        string
	LifetimeInSeconds int64
	Cors              cors.CORS
	Referrer          []string
	Sizelimit         int64
	BaseURL           string
}
