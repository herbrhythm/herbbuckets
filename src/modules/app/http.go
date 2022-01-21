package app

import (
	"sync/atomic"

	"github.com/herb-go/herbconfig/source"
	"github.com/herb-go/util"
	"github.com/herb-go/util/config"
	"github.com/herb-go/util/config/httpconfig"
	"github.com/herb-go/util/config/tomlconfig"
)

//HTTP app http config
var HTTP *httpconfig.Config

var syncHTTP atomic.Value

//StoreHTTP atomically store http config
func (a *appSync) StoreHTTP(c *httpconfig.Config) {
	syncHTTP.Store(c)
}

//LoadHTTP atomically load http config
func (a *appSync) LoadHTTP() *httpconfig.Config {
	v := syncHTTP.Load()
	if v == nil {
		return nil
	}
	return v.(*httpconfig.Config)
}

func init() {
	config.RegisterLoaderAndWatch(util.ConfigFile("/http.toml"), func(configpath source.Source) {
		HTTP = httpconfig.New()
		util.Must(tomlconfig.Load(configpath, HTTP))
		Sync.StoreHTTP(HTTP)
		util.SetWarning("Forwarded", HTTP.Forwarded.Warnings()...)
	})
}
