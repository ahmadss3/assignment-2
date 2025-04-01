// File: assignment-2/structs/external_structs.go
package structs

import "time"

// CacheEntryEx holds cached data, typically used to store JSON from external APIs.
type CacheEntryEx struct {
	Key         string    `firestore:"key"`
	Data        []byte    `firestore:"data"`
	LastFetched time.Time `firestore:"lastFetched"`
	TTLHours    int       `firestore:"ttlHours"`
}

// CountryInfo represents data returned from REST Countries.
type CountryInfo struct {
	Name         string
	Capital      string
	Population   int64
	Area         float64
	BaseCurrency string
	Coordinates  CoordinatesEx
}

// CoordinatesEx represents latitude and longitude.
type CoordinatesEx struct {
	Lat float64
	Lon float64
}

// MeteoData holds average temperature and precipitation from Open-Meteo.
type MeteoData struct {
	AverageTemp          float64
	AveragePrecipitation float64
}

// CurrencyRates is a map from currency code to exchange rate.
type CurrencyRates map[string]float64
