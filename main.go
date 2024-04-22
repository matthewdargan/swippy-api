// Copyright 2024 Matthew P. Dargan.
// SPDX-License-Identifier: Apache-2.0

// Swippy retrieves from the eBay Finding API and stores results in a
// PostgreSQL database.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/matthewdargan/ebay"
)

var (
	method = flag.String("m", "", "eBay client method to call")
	params = flag.String("p", "", "query parameters")
	appID  = os.Getenv("EBAY_APP_ID")
	dbURL  = os.Getenv("DB_URL")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: swippy -m method -p params\n")
	os.Exit(2)
}

func main() {
	log.SetPrefix("swippy: ")
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()
	if *method == "" || *params == "" {
		usage()
	}
	queryParams, err := parseParams(*params)
	if err != nil {
		log.Fatal(err)
	}
	c := ebay.NewFindingClient(&http.Client{Timeout: time.Second * 10}, appID)
	var resps []ebay.FindItemsResponse
	switch *method {
	case "advanced":
		var r *ebay.FindItemsAdvancedResponse
		r, err = c.FindItemsAdvanced(context.Background(), queryParams)
		if err != nil {
			log.Fatal(err)
		}
		resps = r.ItemsResponse
	case "category":
		var r *ebay.FindItemsByCategoryResponse
		r, err = c.FindItemsByCategory(context.Background(), queryParams)
		if err != nil {
			log.Fatal(err)
		}
		resps = r.ItemsResponse
	case "keywords":
		var r *ebay.FindItemsByKeywordsResponse
		r, err = c.FindItemsByKeywords(context.Background(), queryParams)
		if err != nil {
			log.Fatal(err)
		}
		resps = r.ItemsResponse
	case "product":
		var r *ebay.FindItemsByProductResponse
		r, err = c.FindItemsByProduct(context.Background(), queryParams)
		if err != nil {
			log.Fatal(err)
		}
		resps = r.ItemsResponse
	case "ebay-stores":
		var r *ebay.FindItemsInEBayStoresResponse
		r, err = c.FindItemsInEBayStores(context.Background(), queryParams)
		if err != nil {
			log.Fatal(err)
		}
		resps = r.ItemsResponse
	default:
		usage()
	}
	if len(resps) == 0 {
		os.Exit(0)
	}
	if len(resps[0].ErrorMessage) > 0 {
		log.Fatal(resps[0].ErrorMessage)
	}
	log.Print(resps)
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer conn.Close(context.Background())
	insertItems(conn, resps)
}

func parseParams(ps string) (map[string]string, error) {
	params := make(map[string]string)
	for _, p := range strings.Split(ps, "&") {
		parts := strings.Split(p, "=")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid parameter %q", p)
		}
		params[parts[0]] = parts[1]
	}
	return params, nil
}

type eBayItem struct {
	timestamp                                  time.Time
	version                                    string
	conditionDisplayName                       string
	conditionID                                int
	country                                    string
	galleryURL                                 *string
	globalID                                   string
	isMultiVariationListing                    bool
	itemID                                     int64
	listingInfoBestOfferEnabled                bool
	listingInfoBuyItNowAvailable               bool
	listingInfoEndTime                         time.Time
	listingInfoListingType                     string
	listingInfoStartTime                       time.Time
	listingInfoWatchCount                      *int
	location                                   *string
	postalCode                                 *string
	primaryCategoryID                          int64
	primaryCategoryName                        string
	productIDType                              *string
	productIDValue                             *int64
	sellingStatusConvertedCurrentPriceCurrency *string
	sellingStatusConvertedCurrentPriceValue    *float64
	sellingStatusCurrentPriceCurrency          *string
	sellingStatusCurrentPriceValue             *float64
	sellingStatusSellingState                  *string
	sellingStatusTimeLeft                      *string
	shippingServiceCostCurrency                *string
	shippingServiceCostValue                   *float64
	shippingType                               *string
	shipToLocations                            *string
	subtitle                                   *string
	title                                      string
	topRatedListing                            bool
	viewItemURL                                *string
}

func insertItems(conn *pgx.Conn, rs []ebay.FindItemsResponse) {
	var eBayItems []eBayItem
	for _, r := range rs {
		items, err := responseToItems(r)
		if err != nil {
			log.Printf("failed to convert eBay API response to items: %v", err)
			continue
		}
		eBayItems = append(eBayItems, items...)
	}
	_, err := conn.CopyFrom(
		context.Background(), pgx.Identifier{"item"},
		[]string{
			"timestamp", "version", "condition_display_name", "condition_id",
			"country", "gallery_url", "global_id",
			"is_multi_variation_listing", "item_id",
			"listing_info_best_offer_enabled",
			"listing_info_buy_it_now_available", "listing_info_end_time",
			"listing_info_listing_type",
			"listing_info_start_time", "listing_info_watch_count", "location",
			"postal_code", "primary_category_id", "primary_category_name",
			"product_id_type", "product_id_value",
			"selling_status_converted_current_price_currency",
			"selling_status_converted_current_price_value",
			"selling_status_current_price_currency",
			"selling_status_current_price_value",
			"selling_status_selling_state", "selling_status_time_left",
			"shipping_service_cost_currency", "shipping_service_cost_value",
			"shipping_type", "ship_to_locations", "subtitle", "title",
			"top_rated_listing", "view_item_url",
		},
		pgx.CopyFromSlice(len(eBayItems), func(i int) ([]any, error) {
			return []any{
				eBayItems[i].timestamp, eBayItems[i].version,
				eBayItems[i].conditionDisplayName, eBayItems[i].conditionID,
				eBayItems[i].country, eBayItems[i].galleryURL,
				eBayItems[i].globalID, eBayItems[i].isMultiVariationListing,
				eBayItems[i].itemID,
				eBayItems[i].listingInfoBestOfferEnabled,
				eBayItems[i].listingInfoBuyItNowAvailable,
				eBayItems[i].listingInfoEndTime,
				eBayItems[i].listingInfoListingType,
				eBayItems[i].listingInfoStartTime,
				eBayItems[i].listingInfoWatchCount, eBayItems[i].location,
				eBayItems[i].postalCode, eBayItems[i].primaryCategoryID,
				eBayItems[i].primaryCategoryName, eBayItems[i].productIDType,
				eBayItems[i].productIDValue,
				eBayItems[i].sellingStatusConvertedCurrentPriceCurrency,
				eBayItems[i].sellingStatusConvertedCurrentPriceValue,
				eBayItems[i].sellingStatusCurrentPriceCurrency,
				eBayItems[i].sellingStatusCurrentPriceValue,
				eBayItems[i].sellingStatusSellingState,
				eBayItems[i].sellingStatusTimeLeft,
				eBayItems[i].shippingServiceCostCurrency,
				eBayItems[i].shippingServiceCostValue,
				eBayItems[i].shippingType, eBayItems[i].shipToLocations,
				eBayItems[i].subtitle, eBayItems[i].title,
				eBayItems[i].topRatedListing, eBayItems[i].viewItemURL,
			}, nil
		}),
	)
	if err != nil {
		log.Printf("failed to insert data: %v", err)
	}
}

func responseToItems(resp ebay.FindItemsResponse) ([]eBayItem, error) {
	items := make([]eBayItem, len(resp.SearchResult[0].Item))
	for i := range items {
		it, err := item(resp.SearchResult[0].Item[i])
		if err != nil {
			return nil, err
		}
		it.timestamp = resp.Timestamp[0]
		it.version = resp.Version[0]
		items[i] = *it
	}
	return items, nil
}

func item(it ebay.SearchItem) (*eBayItem, error) {
	conditionID, err := strconv.Atoi(it.Condition[0].ConditionID[0])
	if err != nil {
		return nil, fmt.Errorf("cannot convert conditionID to int: %w", err)
	}
	isMultiVariationListing, err := strconv.ParseBool(it.IsMultiVariationListing[0])
	if err != nil {
		return nil, fmt.Errorf("cannot convert isMultiVariationListing to bool: %w", err)
	}
	itemID, err := strconv.ParseInt(it.ItemID[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("cannot convert itemID to int64: %w", err)
	}
	bestOfferEnabled, err := strconv.ParseBool(it.ListingInfo[0].BestOfferEnabled[0])
	if err != nil {
		return nil, fmt.Errorf("cannot convert bestOfferEnabled to bool: %w", err)
	}
	buyItNowAvailable, err := strconv.ParseBool(it.ListingInfo[0].BuyItNowAvailable[0])
	if err != nil {
		return nil, fmt.Errorf("cannot convert buyItNowAvailable to bool: %w", err)
	}
	var watchCount *int
	if len(it.ListingInfo[0].WatchCount) > 0 {
		var v int
		v, err = strconv.Atoi(it.ListingInfo[0].WatchCount[0])
		if err != nil {
			return nil, fmt.Errorf("cannot convert watchCount to int: %w", err)
		}
		watchCount = &v
	}
	primaryCategoryID, err := strconv.ParseInt(it.PrimaryCategory[0].CategoryID[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("cannot convert primaryCategoryID to int64: %w", err)
	}
	var productIDType *string
	var productIDValue *int64
	if len(it.ProductID) > 0 {
		productIDType = &it.ProductID[0].Type
		var v int64
		v, err = strconv.ParseInt(it.ProductID[0].Value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot convert productID value to int64: %w", err)
		}
		productIDValue = &v
	}
	var sellingStatusSellingState, sellingStatusTimeLeft *string
	if len(it.SellingStatus[0].SellingState) > 0 {
		sellingStatusSellingState = &it.SellingStatus[0].SellingState[0]
		sellingStatusTimeLeft = &it.SellingStatus[0].TimeLeft[0]
	}
	var sellingStatusPriceCurrency, sellingStatusConvertedPriceCurrency *string
	var sellingStatusPriceValue, sellingStatusConvertedPriceValue *float64
	if len(it.SellingStatus[0].CurrentPrice) > 0 {
		sellingStatusPriceCurrency = &it.SellingStatus[0].CurrentPrice[0].CurrencyID
		var v float64
		v, err = strconv.ParseFloat(it.SellingStatus[0].CurrentPrice[0].Value, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot convert selling status current price value to float64: %w", err)
		}
		sellingStatusPriceValue = &v
		sellingStatusConvertedPriceCurrency = &it.SellingStatus[0].ConvertedCurrentPrice[0].CurrencyID
		v, err = strconv.ParseFloat(it.SellingStatus[0].ConvertedCurrentPrice[0].Value, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot convert selling status converted current price value to float64: %w", err)
		}
		sellingStatusConvertedPriceValue = &v
	}
	var shippingServiceCurrency, shippingType, shipToLocations *string
	var shippingServiceValue *float64
	if len(it.ShippingInfo[0].ShippingServiceCost) > 0 {
		shippingServiceCurrency = &it.ShippingInfo[0].ShippingServiceCost[0].CurrencyID
		var v float64
		v, err = strconv.ParseFloat(it.ShippingInfo[0].ShippingServiceCost[0].Value, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot convert shipping service cost value to float64: %w", err)
		}
		shippingServiceValue = &v
		shippingType = &it.ShippingInfo[0].ShippingType[0]
		shipToLocations = &it.ShippingInfo[0].ShipToLocations[0]
	}
	topRatedListing, err := strconv.ParseBool(it.TopRatedListing[0])
	if err != nil {
		return nil, fmt.Errorf("cannot convert topRatedListing to bool: %w", err)
	}
	return &eBayItem{
		conditionDisplayName:         it.Condition[0].ConditionDisplayName[0],
		conditionID:                  conditionID,
		country:                      it.Country[0],
		galleryURL:                   firstElem(it.GalleryURL),
		globalID:                     it.GlobalID[0],
		isMultiVariationListing:      isMultiVariationListing,
		itemID:                       itemID,
		listingInfoBestOfferEnabled:  bestOfferEnabled,
		listingInfoBuyItNowAvailable: buyItNowAvailable,
		listingInfoEndTime:           it.ListingInfo[0].EndTime[0],
		listingInfoListingType:       it.ListingInfo[0].ListingType[0],
		listingInfoStartTime:         it.ListingInfo[0].StartTime[0],
		listingInfoWatchCount:        watchCount,
		location:                     firstElem(it.Location),
		postalCode:                   firstElem(it.PostalCode),
		primaryCategoryID:            primaryCategoryID,
		primaryCategoryName:          it.PrimaryCategory[0].CategoryName[0],
		productIDType:                productIDType,
		productIDValue:               productIDValue,
		sellingStatusConvertedCurrentPriceCurrency: sellingStatusConvertedPriceCurrency,
		sellingStatusConvertedCurrentPriceValue:    sellingStatusConvertedPriceValue,
		sellingStatusCurrentPriceCurrency:          sellingStatusPriceCurrency,
		sellingStatusCurrentPriceValue:             sellingStatusPriceValue,
		sellingStatusSellingState:                  sellingStatusSellingState,
		sellingStatusTimeLeft:                      sellingStatusTimeLeft,
		shippingServiceCostCurrency:                shippingServiceCurrency,
		shippingServiceCostValue:                   shippingServiceValue,
		shippingType:                               shippingType,
		shipToLocations:                            shipToLocations,
		subtitle:                                   firstElem(it.Subtitle),
		title:                                      it.Title[0],
		topRatedListing:                            topRatedListing,
		viewItemURL:                                firstElem(it.ViewItemURL),
	}, nil
}

func firstElem(ss []string) *string {
	if len(ss) > 0 {
		return &ss[0]
	}
	return nil
}
