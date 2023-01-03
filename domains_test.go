package fishfish_test

import (
	"testing"
	"time"

	"github.com/existentiality/fishfish-go"
)

func TestGetDomains(t *testing.T) {
	domains, err := client.GetDomains(fishfish.CategoryPhishing, false)

	mustPanic(err)

	t.Logf("got %d phishing domains", len(*domains))
}

func TestGetDomainsFull(t *testing.T) {
	domains, err := client.GetDomainsFull()

	mustPanic(err)

	t.Logf("got %d domains with full data", len(*domains))
}

func TestGetDomain(t *testing.T) {
	domain, err := client.GetDomain("fishfish.gg", true)

	mustPanic(err)

	t.Logf("got domain %s (category %s)", domain.Domain, domain.Category)
}

func TestAddDomain(t *testing.T) {
	if !client.HasPermission(fishfish.APIPermissionDomains) {
		t.Skip("missing permission")
	}

	added, err := client.AddDomain("fishfish.gg", fishfish.CategorySafe, true)

	mustPanic(err)

	t.Logf("added domain %s with category %s", added.Domain, added.Category)
}

func TestUpdateDomain(t *testing.T) {
	if !client.HasPermission(fishfish.APIPermissionDomains) {
		t.Skip("missing permission")
	}

	updated, err := client.UpdateDomain("fishfish.gg", fishfish.CategorySafe)

	mustPanic(err)

	t.Logf("updated domain %s", updated.Domain)
}

func TestUpdateDomainMetadata(t *testing.T) {
	if !client.HasPermission(fishfish.APIPermissionDomains) {
		t.Skip("missing permission")
	}

	updated, err := client.UpdateDomainMetadata("fishfish.gg", fishfish.DomainMetadata{
		Target: "fishfish",
		Active: time.Now(),
	})

	mustPanic(err)

	t.Logf("updated metadata for fishfish.gg (last active %s, target %s)", updated.Active, updated.Target)
}

func TestDeleteDomain(t *testing.T) {
	if !client.HasPermission(fishfish.APIPermissionDomains) {
		t.Skip("missing permission")
	}

	err := client.DeleteDomain("fishfish.gg")

	mustPanic(err)

	t.Logf("successfully deleted domain")
}
