// file: assignment-2/handlers/dashboards_handler.go
package handlers

import (
	"log"
	"net/http"
	"strings"
	"time"

	"assignment-2/constants"
	"assignment-2/firebase"
	"assignment-2/services"
	"assignment-2/structs"
	"assignment-2/tools"
)

// DashboardsRouter handles GET /dashboard/v1/dashboards/{id}
func DashboardsRouter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		tools.WriteJsonErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed on dashboards")
		return
	}
	if r.URL.Path == constants.DASHBOARDS_PATH {
		// We do not allow listing all dashboards
		tools.WriteJsonErrorResponse(w, http.StatusMethodNotAllowed, "Cannot list dashboards")
		return
	}
	// There's something after /dashboards/
	id := strings.TrimPrefix(r.URL.Path, constants.DASHBOARDS_PATH)
	handleGetDashboardByID(w, r, id)
}

// handleGetDashboardByID fetches the corresponding registration and then retrieves real data from external APIs.
func handleGetDashboardByID(w http.ResponseWriter, r *http.Request, id string) {
	reg, err := firebase.GetRegistrationByID(r.Context(), id)
	if err != nil {
		log.Printf("Error retrieving registration for dashboard: %v\n", err)
		tools.WriteJsonErrorResponse(w, http.StatusNotFound, "Registration not found")
		return
	}

	// Build the Dashboard response
	var dash structs.Dashboard
	dash.Country = reg.Country
	dash.ISOCode = reg.ISOCode

	var df structs.DashboardFeatures

	// Decide if we need to fetch country info
	needCountry := reg.Features.Capital || reg.Features.Population || reg.Features.Area || reg.Features.Coordinates

	var cInfo *structs.CountryInfo
	if needCountry {
		key := reg.Country
		if key == "" {
			key = reg.ISOCode
		}
		cInfo, err = services.FetchCountryInfo(key)
		if err != nil {
			log.Printf("Warning: could not fetch country info for '%s': %v\n", key, err)
		}
	}
	if cInfo != nil {
		if reg.Features.Capital {
			df.Capital = cInfo.Capital
		}
		if reg.Features.Coordinates {
			df.Coordinates = &structs.Coordinates{
				Lat: cInfo.Coordinates.Lat,
				Lon: cInfo.Coordinates.Lon,
			}
		}
		if reg.Features.Population {
			df.Population = cInfo.Population
		}
		if reg.Features.Area {
			df.Area = cInfo.Area
		}
	}

	// If temperature/precipitation... call open-meteo
	if (reg.Features.Temperature || reg.Features.Precipitation) && cInfo != nil {
		mData, errM := services.FetchMeteoData(cInfo.Coordinates.Lat, cInfo.Coordinates.Lon)
		if errM == nil && mData != nil {
			if reg.Features.Temperature {
				df.Temperature = mData.AverageTemp
			}
			if reg.Features.Precipitation {
				df.Precipitation = mData.AveragePrecipitation
			}
		} else {
			log.Printf("Warning: fetch meteo data lat=%.2f lon=%.2f: %v\n",
				cInfo.Coordinates.Lat, cInfo.Coordinates.Lon, errM)
		}
	}

	// If targetCurrencies... call currency API if cInfo.BaseCurrency is not empty
	if len(reg.Features.TargetCurrencies) > 0 && cInfo != nil && cInfo.BaseCurrency != "" {
		rates, errC := services.FetchCurrencyRates(cInfo.BaseCurrency)
		if errC == nil && rates != nil {
			tcMap := make(map[string]float64)
			for _, cur := range reg.Features.TargetCurrencies {
				if val, ok := rates[cur]; ok {
					tcMap[cur] = val
				}
			}
			df.TargetCurrencies = tcMap
		} else {
			log.Printf("Warning: fetch currency rates for base=%s: %v\n", cInfo.BaseCurrency, errC)
		}
	}

	dash.Features = df
	dash.LastRetrieval = time.Now()

	// Return the assembled JSON
	tools.WriteJsonResponse(w, http.StatusOK, dash)

	// Trigger 'INVOKE' event
	countryKey := reg.Country
	if countryKey == "" {
		countryKey = reg.ISOCode
	}
	TriggerWebhookEventVar("INVOKE", countryKey)
}
