// File: assignment-2/structs/features.go
package structs

// Features defines which fields should be included in the dashboard.
type Features struct {
	Temperature      bool     `json:"temperature"`      // Temperature should be fetched.
	Precipitation    bool     `json:"precipitation"`    // Precipitation should be fetched.
	Capital          bool     `json:"capital"`          // The name of the capital city should be included.
	Coordinates      bool     `json:"coordinates"`      // The latitude and longitude of the country should be included.
	Population       bool     `json:"population"`       // The total population should be included.
	Area             bool     `json:"area"`             // the total area should be included.
	TargetCurrencies []string `json:"targetCurrencies"` // Is a list of currency codes for which the user
	// wants to see exchange rates relative to the country's base currency.

}

// NewFeatures creates a Features struct with a guaranteed empty slice for TargetCurrencies
// instead of nil.
func NewFeatures() Features {
	return Features{
		Temperature:      false,
		Precipitation:    false,
		Capital:          false,
		Coordinates:      false,
		Population:       false,
		Area:             false,
		TargetCurrencies: []string{},
	}
}
