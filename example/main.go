package main

import (
	"fmt"
	"os"
	"time"

	recyclebelgium "github.com/meyskens/go-recycle-belgium"
)

func main() {
	secret := os.Getenv("XSECRET")

	if secret == "" {
		fmt.Println("XSECRET variable is not set")
		os.Exit(1)
	}

	api := recyclebelgium.NewAPI(secret)

	// get the ZIP Code ID for 1060 Sint-Gillis
	zipResp, err := api.GetZipCodes("1060")
	if err != nil {
		fmt.Printf("Error getting ZIP codes: %v\n", err)
		os.Exit(1)
	}
	if len(zipResp.Items) < 1 {
		fmt.Println("Got no ZIP Code matches")
		os.Exit(1)
	}

	zipID := zipResp.Items[0].ID // assuming an exact match was given

	cityResp, err := api.GetStreets(zipID, "Coenraetsstraat")
	if err != nil {
		fmt.Printf("Error getting streets: %v\n", err)
		os.Exit(1)
	}
	if len(cityResp.Items) < 1 {
		fmt.Println("Got no street matches")
		os.Exit(1)
	}
	streetID := cityResp.Items[0].ID

	collections, err := api.GetCollections(zipID, streetID, "72", time.Now(), time.Now().Add(7*24*time.Hour), 100)
	if len(cityResp.Items) < 1 {
		fmt.Printf("Error getting collections: %v\n", err)
		os.Exit(1)
	}

	for _, collection := range collections.Items {
		fmt.Printf("%s will be collected %s\n", collection.Fraction.Name.EN, collection.Timestamp.Format("2006-01-02"))
	}

}
