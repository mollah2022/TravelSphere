package utils

import (
	"TravelSphere/models"
	"strings"
)

// GetMockCountriesData returns fallback countries data when API fails
func GetMockCountriesData() []models.CountryResponse {
	return []models.CountryResponse{
		{
			Name:       models.CountryName{Common: "France", Official: "French Republic"},
			CCA2:       "FR",
			CCA3:       "FRA",
			Capital:    []string{"Paris"},
			Region:     "Europe",
			Subregion:  "Western Europe",
			Population: 67970000,
			Flags: models.CountryFlag{
				PNG: "https://flagcdn.com/w320/fr.png",
				SVG: "https://flagcdn.com/fr.svg",
				Alt: "The flag of France",
			},
			Currencies: map[string]models.Currency{
				"EUR": {Name: "Euro", Symbol: "€"},
			},
			Languages: map[string]string{"fra": "French"},
			LatLng:    []float64{46.0, 2.0},
		},
		{
			Name:       models.CountryName{Common: "Spain", Official: "Kingdom of Spain"},
			CCA2:       "ES",
			CCA3:       "ESP",
			Capital:    []string{"Madrid"},
			Region:     "Europe",
			Subregion:  "Southern Europe",
			Population: 47615000,
			Flags: models.CountryFlag{
				PNG: "https://flagcdn.com/w320/es.png",
				SVG: "https://flagcdn.com/es.svg",
				Alt: "The flag of Spain",
			},
			Currencies: map[string]models.Currency{
				"EUR": {Name: "Euro", Symbol: "€"},
			},
			Languages: map[string]string{"spa": "Spanish"},
			LatLng:    []float64{40.0, -3.0},
		},
		{
			Name:       models.CountryName{Common: "Italy", Official: "Italian Republic"},
			CCA2:       "IT",
			CCA3:       "ITA",
			Capital:    []string{"Rome"},
			Region:     "Europe",
			Subregion:  "Southern Europe",
			Population: 58940000,
			Flags: models.CountryFlag{
				PNG: "https://flagcdn.com/w320/it.png",
				SVG: "https://flagcdn.com/it.svg",
				Alt: "The flag of Italy",
			},
			Currencies: map[string]models.Currency{
				"EUR": {Name: "Euro", Symbol: "€"},
			},
			Languages: map[string]string{"ita": "Italian"},
			LatLng:    []float64{41.871940, 12.567380},
		},
		{
			Name:       models.CountryName{Common: "Japan", Official: "Japan"},
			CCA2:       "JP",
			CCA3:       "JPN",
			Capital:    []string{"Tokyo"},
			Region:     "Asia",
			Subregion:  "Eastern Asia",
			Population: 125124000,
			Flags: models.CountryFlag{
				PNG: "https://flagcdn.com/w320/jp.png",
				SVG: "https://flagcdn.com/jp.svg",
				Alt: "The flag of Japan",
			},
			Currencies: map[string]models.Currency{
				"JPY": {Name: "Japanese yen", Symbol: "¥"},
			},
			Languages: map[string]string{"jpn": "Japanese"},
			LatLng:    []float64{36.2048, 138.2529},
		},
		{
			Name:       models.CountryName{Common: "Brazil", Official: "Federative Republic of Brazil"},
			CCA2:       "BR",
			CCA3:       "BRA",
			Capital:    []string{"Brasília"},
			Region:     "Americas",
			Subregion:  "South America",
			Population: 215313000,
			Flags: models.CountryFlag{
				PNG: "https://flagcdn.com/w320/br.png",
				SVG: "https://flagcdn.com/br.svg",
				Alt: "The flag of Brazil",
			},
			Currencies: map[string]models.Currency{
				"BRL": {Name: "Brazilian real", Symbol: "R$"},
			},
			Languages: map[string]string{"por": "Portuguese"},
			LatLng:    []float64{-14.235, -51.9253},
		},
		{
			Name:       models.CountryName{Common: "United States", Official: "United States of America"},
			CCA2:       "US",
			CCA3:       "USA",
			Capital:    []string{"Washington, D.C."},
			Region:     "Americas",
			Subregion:  "North America",
			Population: 339996563,
			Flags: models.CountryFlag{
				PNG: "https://flagcdn.com/w320/us.png",
				SVG: "https://flagcdn.com/us.svg",
				Alt: "The flag of the United States of America",
			},
			Currencies: map[string]models.Currency{
				"USD": {Name: "United States dollar", Symbol: "$"},
			},
			Languages: map[string]string{"eng": "English"},
			LatLng:    []float64{37.0902, -95.7129},
		},
		{
			Name:       models.CountryName{Common: "Australia", Official: "Commonwealth of Australia"},
			CCA2:       "AU",
			CCA3:       "AUS",
			Capital:    []string{"Canberra"},
			Region:     "Oceania",
			Subregion:  "Australia and New Zealand",
			Population: 26068792,
			Flags: models.CountryFlag{
				PNG: "https://flagcdn.com/w320/au.png",
				SVG: "https://flagcdn.com/au.svg",
				Alt: "The flag of Australia",
			},
			Currencies: map[string]models.Currency{
				"AUD": {Name: "Australian dollar", Symbol: "$"},
			},
			Languages: map[string]string{"eng": "English"},
			LatLng:    []float64{-27.0, 133.0},
		},
		{
			Name:       models.CountryName{Common: "Egypt", Official: "Arab Republic of Egypt"},
			CCA2:       "EG",
			CCA3:       "EGY",
			Capital:    []string{"Cairo"},
			Region:     "Africa",
			Subregion:  "Northern Africa",
			Population: 110672000,
			Flags: models.CountryFlag{
				PNG: "https://flagcdn.com/w320/eg.png",
				SVG: "https://flagcdn.com/eg.svg",
				Alt: "The flag of Egypt",
			},
			Currencies: map[string]models.Currency{
				"EGP": {Name: "Egyptian pound", Symbol: "£"},
			},
			Languages: map[string]string{"ara": "Arabic"},
			LatLng:    []float64{26.8206, 30.8025},
		},
	}
}

// searchMockCountriesByName searches for countries by name in mock data
func searchMockCountriesByName(name string) []models.CountryResponse {
	name = strings.ToLower(strings.TrimSpace(name))
	mockData := GetMockCountriesData()
	result := make([]models.CountryResponse, 0)

	for _, country := range mockData {
		if strings.Contains(strings.ToLower(country.Name.Common), name) ||
			strings.Contains(strings.ToLower(country.Name.Official), name) {
			result = append(result, country)
		}
	}

	return result
}
