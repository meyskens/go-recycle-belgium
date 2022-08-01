# Go API for Recycle! (Belgium)

The [Recycle! App](https://recycleapp.be) is an app/website that allows you to get dates for waste collection in Belgium.
There is however no official/open API for this. This project aims to provide a Go interface to the API used by the website, it is different than the app but will be easier to troubleshoot.

## What is `x-secret`?

This is a random string that is needed to get an auth token for some reasons... Just aim your webbrowser inspection tool at https://recycleapp.be/home and look for `https://api.fostplus.be/recycle-public/app/v1/access-token`, its not as secret as it's name might suggest.

## Example use

There is a minimal functional example in `example/main.go`

```go
package main

import (
	"fmt"
	"os"
	"time"

	recyclebelgium "github.com/meyskens/go-recycle-belgium"
)

func main() {
	api := recyclebelgium.NewAPI("---INSERT x-secret HERE---")

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
```
