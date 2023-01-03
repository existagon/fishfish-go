package fishfish_test

import (
	"github.com/existentiality/fishfish-go"
	"os"
	"testing"
)

var primaryKey = os.Getenv("FISHFISH_API_KEY")
var client *fishfish.Client

func TestNewClient(t *testing.T) {
	var err error
	client, err = fishfish.New(primaryKey, []fishfish.APIPermission{})

	mustPanic(err)

	t.Log("successfully created new client")
}

func mustPanic(err error) {
	if err != nil {
		panic(err)
	}
}
