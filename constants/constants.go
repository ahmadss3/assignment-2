// File: assignment-2/constants/constants.go
// This file defines various constants used throughout the service.

package constants

// VERSION specifies the current version of the API.
const VERSION = "v1"

// BASE_PATH forms the base path for the API, for example "/dashboard/v1/".
const BASE_PATH = "/dashboard/" + VERSION + "/"

// Paths for specific resources:
const REGISTRATIONS_PATH = BASE_PATH + "registrations/"
const DASHBOARDS_PATH = BASE_PATH + "dashboards/"
const NOTIFICATIONS_PATH = BASE_PATH + "notifications/"
const STATUS_PATH = BASE_PATH + "status/"

// DefaultPort defines the default port for the service
const DefaultPort = "8080"

// Production external API endpoints
const REST_COUNTRIES_ALPHA = "http://129.241.150.113:8080/v3.1/alpha/"
const REST_COUNTRIES_NAME = "http://129.241.150.113:8080/v3.1/name/"
const CURRENCY_API = "http://129.241.150.113:9090/currency/"
const OPEN_METEO_API = "https://api.open-meteo.com/v1/forecast"

// Firebase collection names (if used in a real Firestore scenario)
const REGISTRATIONS_COLLECTION = "registrations"
const NOTIFICATIONS_COLLECTION = "notifications"
const CACHE_COLLECTION = "cache"

// Local mock_data file paths
const MOCKDATA_RESTCOUNTRIES_NORWAY = "mock_data/restcountries_norway.json"
const MOCKDATA_WEATHER_NORWAY = "mock_data/weather_norway.json"
const MOCKDATA_CURRENCY_NOK = "mock_data/currency_response_nok.json"

// ServiceVersion can be used in logs or status endpoints
const ServiceVersion = "v1.0.0"
