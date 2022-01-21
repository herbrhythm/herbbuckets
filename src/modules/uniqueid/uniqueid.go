package uniqueid

import (
	"herbbuckets/modules/app"

	"github.com/herb-go/uniqueid"
	"github.com/herb-go/util"
)

//ModuleName module name
const ModuleName = "100uniqueid"

//Generator unique id generator
var Generator *uniqueid.Generator

//GenerateID generate unique id.
//Return  generated id and any error if rasied.
func GenerateID() (string, error) {
	return Generator.GenerateID()
}

//MustGenerateID generate unique id.
//Return  generated id.
//Panic if any error raised
func MustGenerateID() string {
	return Generator.MustGenerateID()
}

func init() {
	util.RegisterModule(ModuleName, func() {
		Generator = uniqueid.NewGenerator()
		util.Must(app.UniqueID.ApplyTo(Generator))
		uniqueid.DefaultGenerator = Generator
	})
}
