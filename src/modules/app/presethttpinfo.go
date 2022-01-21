package app

import (
	"sync/atomic"

	"github.com/herb-go/herbmodules/httpinfomanager"

	"github.com/herb-go/herbconfig/source"
	"github.com/herb-go/util"
	"github.com/herb-go/util/config"
	"github.com/herb-go/util/config/tomlconfig"
)

//Presethttpinfo config instance of presethttpinfo.
var Presethttpinfo = &httpinfomanager.Config{}

var syncPresethttpinfo atomic.Value

//StorePresethttpinfo atomically store presethttpinfo config
func (a *appSync) StorePresethttpinfo(c *httpinfomanager.Config) {
	syncPresethttpinfo.Store(c)
}

//LoadPresethttpinfo atomically load presethttpinfo config
func (a *appSync) LoadPresethttpinfo() *httpinfomanager.Config {
	v := syncPresethttpinfo.Load()
	if v == nil {
		return nil
	}
	return v.(*httpinfomanager.Config)
}

func init() {
	//Register loader which will be execute when Config.LoadAll func be called.
	//You can put your init code after load.
	//You must panic if any error rasied when init.
	config.RegisterLoader(util.ConstantsFile("/presethttpinfo.toml"), func(configpath source.Source) {
		util.Must(tomlconfig.Load(configpath, Presethttpinfo))
		Sync.StorePresethttpinfo(Presethttpinfo)
	})
}
