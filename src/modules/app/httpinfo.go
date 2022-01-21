package app

import (
	"sync/atomic"

	"github.com/herb-go/herbconfig/source"
	"github.com/herb-go/herbmodules/httpinfomanager"
	"github.com/herb-go/util"
	"github.com/herb-go/util/config"
	"github.com/herb-go/util/config/tomlconfig"
)

//Httpinfo config instance of httpinfo.
var Httpinfo = &httpinfomanager.Config{}

var syncHttpinfo atomic.Value

//StoreHttpinfo atomically store httpinfo config
func (a *appSync) StoreHttpinfo(c *httpinfomanager.Config) {
	syncHttpinfo.Store(c)
}

//LoadHttpinfo atomically load httpinfo config
func (a *appSync) LoadHttpinfo() *httpinfomanager.Config {
	v := syncHttpinfo.Load()
	if v == nil {
		return nil
	}
	return v.(*httpinfomanager.Config)
}

func init() {
	//Register loader which will be execute when Config.LoadAll func be called.
	//You can put your init code after load.
	//You must panic if any error rasied when init.
	config.RegisterLoader(util.ConfigFile("/httpinfo.toml"), func(configpath source.Source) {
		util.Must(tomlconfig.Load(configpath, Httpinfo))
		Sync.StoreHttpinfo(Httpinfo)
	})
}
