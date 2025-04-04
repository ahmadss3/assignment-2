// File: assignment-2/services/external_services.go
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"assignment-2/constants"
	"assignment-2/firebase"
	"assignment-2/structs"
)

// Function variables for test stubbing
var (
	FetchCountryInfo   func(countryOrISO string) (*structs.CountryInfo, error) = realFetchCountryInfo
	FetchMeteoData     func(lat, lon float64) (*structs.MeteoData, error)      = realFetchMeteoData
	FetchCurrencyRates func(base string) (structs.CurrencyRates, error)        = realFetchCurrencyRates
)

// realFetchCountryInfo checks Firestore cache first, then calls callRestCountries if not found or parse fails

func realFetchCountryInfo(countryOrISO string) (*structs.CountryInfo, error) {
	ctx := context.Background()
	cacheKey := "country:" + strings.ToUpper(countryOrISO)

	// 1) Attempt to read from Firestore-based cache
	cached, err := firebase.GetCacheEntry(ctx, cacheKey)
	if err == nil && cached != nil {
		var cInfo structs.CountryInfo
		if unmarshalErr := json.Unmarshal(cached.Data, &cInfo); unmarshalErr == nil {
			// Cache HIT: if not unmarshal properly, return cached data
			return &cInfo, nil
		}
		// If unmarshal fails, we fall through and do a real call
	}

	// 2) Not in cache or failed to unmarshal means do real call
	cInfo, err := callRestCountries(countryOrISO)
	if err != nil {
		return nil, err
	}

	// 3) Save to cache so future lookups can avoid a real call
	rawBytes, marshalErr := json.Marshal(cInfo)
	if marshalErr == nil {
		saveErr := firebase.SaveCacheEntry(ctx, structs.CacheEntry{
			Key:         cacheKey,
			Data:        rawBytes,
			LastFetched: time.Now(),
			TTLHours:    24,
		})
		if saveErr != nil {

		}
	}

	// Return the newly fetched data
	return cInfo, nil
}

// callRestCountries does a real HTTP request to REST Countries
func callRestCountries(countryOrISO string) (*structs.CountryInfo, error) {
	url := fmt.Sprintf("%s%s?fields=name,capital,population,area,latlng,currencies",
		constants.REST_COUNTRIES_NAME,
		countryOrISO,
	)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call REST Countries: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("REST Countries returned %d => %s", resp.StatusCode, string(body))
	}

	// parse "currencies" as a map to find the first key
	type restCountry struct {
		Name struct {
			Common string `json:"common"`
		} `json:"name"`
		Capital    []string               `json:"capital"`
		Population int64                  `json:"population"`
		Area       float64                `json:"area"`
		Latlng     []float64              `json:"latlng"`
		Currencies map[string]interface{} `json:"currencies"`
	}

	var parsed []restCountry
	if decodeErr := json.NewDecoder(resp.Body).Decode(&parsed); decodeErr != nil {
		return nil, fmt.Errorf("failed to decode restcountries JSON: %v", decodeErr)
	}
	if len(parsed) == 0 {
		return nil, fmt.Errorf("no country data found for %s", countryOrISO)
	}

	first := parsed[0]
	cInfo := &structs.CountryInfo{
		Name:         first.Name.Common,
		Capital:      "",
		Population:   first.Population,
		Area:         first.Area,
		BaseCurrency: "",
		Coordinates:  structs.Coordinates{},
	}

	// Capital
	if len(first.Capital) > 0 {
		cInfo.Capital = first.Capital[0]
	}
	// lat/long
	if len(first.Latlng) == 2 {
		cInfo.Coordinates.Lat = first.Latlng[0]
		cInfo.Coordinates.Lon = first.Latlng[1]
	}
	// currency
	if len(first.Currencies) > 0 {
		for key := range first.Currencies {
			cInfo.BaseCurrency = key
			break
		}
	}

	return cInfo, nil
}

// realFetchMeteoData fetches average temperature and precipitation from open-meteo
func realFetchMeteoData(lat, lon float64) (*structs.MeteoData, error) {
	url := fmt.Sprintf("%s?latitude=%.4f&longitude=%.4f&hourly=temperature_2m,precipitation",
		constants.OPEN_METEO_API,
		lat,
		lon,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call open-meteo: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("open-meteo returned %d => %s", resp.StatusCode, string(body))
	}

	var parsed struct {
		Hourly struct {
			Temperature2m []float64 `json:"temperature_2m"`
			Precipitation []float64 `json:"precipitation"`
		} `json:"hourly"`
	}
	if dErr := json.NewDecoder(resp.Body).Decode(&parsed); dErr != nil {
		return nil, fmt.Errorf("decode error from open-meteo: %v", dErr)
	}

	temps := parsed.Hourly.Temperature2m
	precs := parsed.Hourly.Precipitation

	var sumT, sumP float64
	for _, v := range temps {
		sumT += v
	}
	avgT := 0.0
	if len(temps) > 0 {
		avgT = sumT / float64(len(temps))
	}
	for _, p := range precs {
		sumP += p
	}
	avgP := 0.0
	if len(precs) > 0 {
		avgP = sumP / float64(len(precs))
	}

	return &structs.MeteoData{
		AverageTemp:          avgT,
		AveragePrecipitation: avgP,
	}, nil
}

// realFetchCurrencyRates calls the currency API to retrieve exchange rates
func realFetchCurrencyRates(base string) (structs.CurrencyRates, error) {
	url := fmt.Sprintf("%s%s", constants.CURRENCY_API, strings.ToUpper(base))
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call currency API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("currency API returned %d => %s", resp.StatusCode, string(body))
	}

	var parsed struct {
		Result string             `json:"result"`
		Rates  map[string]float64 `json:"rates"`
	}
	if decodeErr := json.NewDecoder(resp.Body).Decode(&parsed); decodeErr != nil {
		return nil, fmt.Errorf("decode error currency API: %v", decodeErr)
	}
	if parsed.Result != "success" {
		return nil, fmt.Errorf("currency API: result=%s (not success)", parsed.Result)
	}
	return parsed.Rates, nil
}
