package fishfish_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/existentiality/fishfish-go"
)

var autoClient *fishfish.AutoSyncClient

func TestAutoSyncNewClient(t *testing.T) {
	var err error
	autoClient, err = fishfish.NewAutoSync(primaryKey, []fishfish.APIPermission{})

	mustPanic(err)

	t.Log("successfully created new autosync")
}

// This is just to start it
func TestAutoSyncForceSync(t *testing.T) {
	err := autoClient.ForceSync()

	mustPanic(err)
}

func TestAutoSyncStart(t *testing.T) {
	autoClient.StartAutoSync(time.Minute * 5)
}

func TestAutoSyncGetDomains(t *testing.T) {
	domains := autoClient.GetDomains()

	if len(domains) == 0 {
		panic("no domains")
	}
}

func TestAutoSyncGetDomain(t *testing.T) {
	domain, err := autoClient.GetDomain("fishfish.gg")

	mustPanic(err)

	expected := fishfish.Domain{
		Domain:      "fishfish.gg",
		Description: "Submitted via FishFish Discord",
		Category:    fishfish.CategorySafe,
		Added:       1667617118,
		Checked:     1667617118,
	}

	if *domain != expected {
		panic(fmt.Errorf("expected domain %v, got %v", expected, domain))
	}
}

func TestAutoSyncStop(t *testing.T) {
	autoClient.StopAutoSync()
}
