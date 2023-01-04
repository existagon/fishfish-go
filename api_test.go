package fishfish_test

import (
	"github.com/existentiality/fishfish-go"
	"os"
	"testing"
)

var primaryKey = os.Getenv("FISHFISH_API_KEY")
var rawClient *fishfish.RawClient

func TestNewClient(t *testing.T) {
	var err error
	rawClient, err = fishfish.NewRaw(primaryKey, []fishfish.APIPermission{})

	mustPanic(err)

	t.Log("successfully created new client")
}

func mustPanic(err error) {
	if err != nil {
		panic(err)
	}
}
