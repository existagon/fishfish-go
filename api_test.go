package fishfish_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/existentiality/fishfish-go"
)

var primaryKey = os.Getenv("FISHFISH_API_KEY")
var rawClient *fishfish.RawClient

func TestNewClient(t *testing.T) {
	var err error
	rawClient, err = fishfish.NewRaw(primaryKey, []fishfish.APIPermission{})

	mustPanic(err)

	t.Log("successfully created new client")
}

func TestErrors(t *testing.T) {
	_, err := fishfish.NewRaw(primaryKey, []fishfish.APIPermission{fishfish.APIPermissionDomains, fishfish.APIPermissionURLs})

	expected := "failed to create session token: unauthorized for specified permission(s)"
	if err.Error() != expected {
		panic(fmt.Errorf("incorrect error message for 403. expected %s got %s", expected, err))
	}

	_, err = fishfish.NewRaw("INVALID_KEY", []fishfish.APIPermission{fishfish.APIPermissionDomains, fishfish.APIPermissionURLs})

	expected = "failed to create session token: invalid FishFish API Token"
	if err.Error() != expected {
		panic(fmt.Errorf("incorrect error message for 401. expected %s got %s", expected, err))
	}

	_, err = rawClient.GetDomain("example.com", false)

	expected = "resource not found"
	if err.Error() != expected {
		panic(fmt.Errorf("incorrect error message for 404. expected %s got %s", expected, err))
	}
}

func mustPanic(err error) {
	if err != nil {
		panic(err)
	}
}
