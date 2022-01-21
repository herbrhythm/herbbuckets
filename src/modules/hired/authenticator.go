package hired

import (
	"github.com/herb-go/herbmodule-drivers/protecter-drivers/authenticator/tomlappkey"
)

var TOMLAppKey = tomlappkey.AuthenticatorFactory

var TOMLAppKeySigner = tomlappkey.SignerAuthenticatorFactory
