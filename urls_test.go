package fishfish_test

import (
	"testing"

	"github.com/existentiality/fishfish-go"
)

func TestGetURLs(t *testing.T) {
	urls, err := rawClient.GetURLs(fishfish.CategoryPhishing)

	mustPanic(err)

	t.Logf("got %d phishing urls", len(*urls))
}

func TestGetURLsFull(t *testing.T) {
	urls, err := rawClient.GetURLsFull()

	mustPanic(err)

	t.Logf("got %d urls with full data", len(*urls))
}

func TestGetURL(t *testing.T) {
	// There are currently no URLs in the databse, skip
	t.SkipNow()
	url, err := rawClient.GetURL("https://fishfish.gg/api.html")

	mustPanic(err)

	t.Logf("got url %s (category %s)", url.URL, url.Category)
}

func TestAddURL(t *testing.T) {
	if !rawClient.HasPermission(fishfish.APIPermissionURLs) {
		t.Skip("missing permission")
	}

	added, err := rawClient.AddURL("https://api.fishfish.gg/v1/docs", fishfish.CreateURLRequest{
		Category:    fishfish.CategorySafe,
		Description: "FishFish API v1 Docs",
	})

	mustPanic(err)

	t.Logf("added url %s with category %s", added.URL, added.Category)
}

func TestUpdateURL(t *testing.T) {
	if !rawClient.HasPermission(fishfish.APIPermissionURLs) {
		t.Skip("missing permission")
	}

	err := rawClient.UpdateURL("https://api.fishfish.gg/v1/docs", fishfish.UpdateURLRequest{
		Description: "Amazing FishFish API v1 Docs",
	})

	mustPanic(err)
}

func TestDeleteURL(t *testing.T) {
	if !rawClient.HasPermission(fishfish.APIPermissionURLs) {
		t.Skip("missing permission")
	}

	err := rawClient.DeleteURL("https://api.fishfish.gg/v1/docs")

	mustPanic(err)

	t.Logf("successfully deleted url")
}
