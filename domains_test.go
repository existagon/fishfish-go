package fishfish_test

import (
	"testing"

	"github.com/existentiality/fishfish-go"
)

func TestGetDomains(t *testing.T) {
	domains, err := rawClient.GetDomains(fishfish.CategoryPhishing)

	mustPanic(err)

	t.Logf("got %d phishing domains", len(*domains))
}

func TestGetDomainsFull(t *testing.T) {
	domains, err := rawClient.GetDomainsFull()

	mustPanic(err)

	t.Logf("got %d domains with full data", len(*domains))
}

func TestGetDomain(t *testing.T) {
	domain, err := rawClient.GetDomain("fishfish.gg")

	mustPanic(err)

	t.Logf("got domain %s (category %s)", domain.Domain, domain.Category)
}

func TestAddDomain(t *testing.T) {
	if !rawClient.HasPermission(fishfish.APIPermissionDomains) {
		t.Skip("missing permission")
	}

	added, err := rawClient.AddDomain("fishfish.gg", fishfish.CreateDomainRequest{
		Category:    fishfish.CategorySafe,
		Description: "FishFish official site",
	})

	mustPanic(err)

	t.Logf("added domain %s with category %s", added.Domain, added.Category)
}

func TestUpdateDomain(t *testing.T) {
	if !rawClient.HasPermission(fishfish.APIPermissionDomains) {
		t.Skip("missing permission")
	}

	updated, err := rawClient.UpdateDomain("fishfish.gg", fishfish.UpdateDomainRequest{
		Category: fishfish.CategoryMalware,
	})

	mustPanic(err)

	t.Logf("updated domain %s", updated.Domain)
}

func TestDeleteDomain(t *testing.T) {
	if !rawClient.HasPermission(fishfish.APIPermissionDomains) {
		t.Skip("missing permission")
	}

	err := rawClient.DeleteDomain("fishfish.gg")

	mustPanic(err)

	t.Logf("successfully deleted domain")
}
