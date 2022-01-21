package app

import (
	"sync/atomic"

	"github.com/herb-go/herb/middleware/csrf"
	"github.com/herb-go/herbconfig/source"
	"github.com/herb-go/util"
	"github.com/herb-go/util/config"
	"github.com/herb-go/util/config/tomlconfig"
)

//Csrf crsf  middleware component
var Csrf = &csrf.Config{}

var syncCsrf atomic.Value

//StoreCsrf atomically store csrf config
func (a *appSync) StoreCsrf(c *csrf.Config) {
	syncCsrf.Store(c)
}

//LoadCsrf atomically load csrf config
func (a *appSync) LoadCsrf() *csrf.Config {
	v := syncCsrf.Load()
	if v == nil {
		return nil
	}
	return v.(*csrf.Config)
}

func init() {
	config.RegisterLoader(util.ConfigFile("/csrf.toml"), func(configpath source.Source) {
		util.Must(tomlconfig.Load(configpath, Csrf))
		Sync.StoreCsrf(Csrf)
	})
}
