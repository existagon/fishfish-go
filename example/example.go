package main

import (
	"fmt"
	"os"

	"github.com/existentiality/fishfish-go"
)

func main() {
	apiToken := os.Getenv("FISHFISH_API_KEY")

	ffClient, err := fishfish.NewAutoSync(apiToken, []fishfish.APIPermission{})

	if err != nil {
		panic(err)
	}

	// Start automatically syncing domains every 5 minutes
	ffClient.StartAutoSync()

	// Get domain
	domain, err := ffClient.GetDomain("fishfish.gg")

	if err != nil {
		// Domain does not exist
		fmt.Println(err)
	} else {
		fmt.Println("Domain: ", domain.Domain)
		fmt.Println("Description: ", domain.Description)
		fmt.Println("Category: ", domain.Category)
	}
}
