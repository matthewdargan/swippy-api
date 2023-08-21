package ebay_test

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/matthewdargan/swippy-api/ebay"
)

func ExampleFindingClient_FindItemsByCategories() {
	params := map[string]string{
		"categoryId":            "9355",
		"itemFilter.name":       "MaxPrice",
		"itemFilter.value":      "500.0",
		"itemFilter.paramName":  "Currency",
		"itemFilter.paramValue": "EUR",
	}
	fc := ebay.NewFindingClient(&http.Client{Timeout: time.Second * 5}, "your_app_id")
	resp, err := fc.FindItemsByCategories(context.Background(), params)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp)
	}
	// Output:
	// ebay: failed to perform eBay Finding API request with status code 500
}

func ExampleFindingClient_FindItemsByKeywords() {
	params := map[string]string{
		"keywords":              "iphone",
		"itemFilter.name":       "MaxPrice",
		"itemFilter.value":      "500.0",
		"itemFilter.paramName":  "Currency",
		"itemFilter.paramValue": "EUR",
	}
	fc := ebay.NewFindingClient(&http.Client{Timeout: time.Second * 5}, "your_app_id")
	resp, err := fc.FindItemsByKeywords(context.Background(), params)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp)
	}
	// Output:
	// ebay: failed to perform eBay Finding API request with status code 500
}

func ExampleFindingClient_FindItemsAdvanced() {
	params := map[string]string{
		"categoryId":            "9355",
		"keywords":              "iphone",
		"itemFilter.name":       "MaxPrice",
		"itemFilter.value":      "500.0",
		"itemFilter.paramName":  "Currency",
		"itemFilter.paramValue": "EUR",
	}
	fc := ebay.NewFindingClient(&http.Client{Timeout: time.Second * 5}, "your_app_id")
	resp, err := fc.FindItemsAdvanced(context.Background(), params)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp)
	}
	// Output:
	// ebay: failed to perform eBay Finding API request with status code 500
}

func ExampleFindingClient_FindItemsByProduct() {
	params := map[string]string{
		"productId.@type":       "ISBN",
		"productId":             "9780131101630",
		"itemFilter.name":       "MaxPrice",
		"itemFilter.value":      "50.0",
		"itemFilter.paramName":  "Currency",
		"itemFilter.paramValue": "EUR",
	}
	fc := ebay.NewFindingClient(&http.Client{Timeout: time.Second * 5}, "your_app_id")
	resp, err := fc.FindItemsByProduct(context.Background(), params)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp)
	}
	// Output:
	// ebay: failed to perform eBay Finding API request with status code 500
}

func ExampleFindingClient_FindItemsInEBayStores() {
	params := map[string]string{
		"storeName":             "Supplytronics",
		"itemFilter.name":       "MaxPrice",
		"itemFilter.value":      "50.0",
		"itemFilter.paramName":  "Currency",
		"itemFilter.paramValue": "EUR",
	}
	fc := ebay.NewFindingClient(&http.Client{Timeout: time.Second * 5}, "your_app_id")
	resp, err := fc.FindItemsInEBayStores(context.Background(), params)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp)
	}
	// Output:
	// ebay: failed to perform eBay Finding API request with status code 500
}
