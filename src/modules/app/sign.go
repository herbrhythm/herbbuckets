package app

import (
	"sync/atomic"

	"github.com/herb-go/herbconfig/source"
	"github.com/herb-go/util"
	"github.com/herb-go/util/config"
	"github.com/herb-go/util/config/tomlconfig"
)

//SignConfig sign config data struct.
//Fields must can be unmarshaled by Toml lib.
//You comment this struct if you use third party config struct.
type SignConfig struct {
	DefaultHasher  string
	AppidField     string
	TimestampField string
	SecretField    string
	SignField      string
	ObjectField    string
	BucketField    string
	SizeLimitField string
}

//Sign config instance of sign.
var Sign = &SignConfig{}

var syncSign atomic.Value

//StoreSign atomically store sign config
func (a *appSync) StoreSign(c *SignConfig) {
	syncSign.Store(c)
}

//LoadSign atomically load sign config
func (a *appSync) LoadSign() *SignConfig {
	v := syncSign.Load()
	if v == nil {
		return nil
	}
	return v.(*SignConfig)
}

func init() {
	//Register loader which will be execute when Config.LoadAll func be called.
	//You can put your init code after load.
	//You must panic if any error rasied when init.
	config.RegisterLoader(util.ConstantsFile("/sign.toml"), func(configpath source.Source) {
		util.Must(tomlconfig.Load(configpath, Sign))
		Sync.StoreSign(Sign)
	})
}
