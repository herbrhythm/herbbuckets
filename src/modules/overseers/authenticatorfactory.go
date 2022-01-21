package overseers

import (
	overseer "github.com/herb-go/herbmodule-drivers/protecter-drivers/overseers/authenticatorfactoryoverseer"
	"github.com/herb-go/herbmodules/protecter/authenticator"
	worker "github.com/herb-go/worker"
)

//AuthenticatorFactoryWorker authenticator factory worker.
var AuthenticatorFactoryWorker authenticator.AuthenticatorFactory

//AuthenticatorFactoryOverseer authenticator factory overseer
var AuthenticatorFactoryOverseer = worker.NewOrverseer("authenticatorfactory", &AuthenticatorFactoryWorker)

func init() {
	AuthenticatorFactoryOverseer.WithInitFunc(func(t *worker.OverseerTranning) error {
		return overseer.New().ApplyTo(AuthenticatorFactoryOverseer)
	})
	worker.Appoint(AuthenticatorFactoryOverseer)
}
