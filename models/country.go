package models

type CountryResponse struct {
	Name       CountryName         `json:"name"`
	TLD        []string            `json:"tld"`
	CCA2       string              `json:"cca2"`
	CCA3       string              `json:"cca3"`
	Capital    []string            `json:"capital"`
	Region     string              `json:"region"`
	Subregion  string              `json:"subregion"`
	Population int64               `json:"population"`
	Flags      CountryFlag         `json:"flags"`
	Currencies map[string]Currency `json:"currencies"`
	Languages  map[string]string   `json:"languages"`
	LatLng     []float64           `json:"latlng"`
	Timezones  []string            `json:"timezones"`
}

type CountryName struct {
	Common   string `json:"common"`
	Official string `json:"official"`
}

type CountryFlag struct {
	PNG string `json:"png"`
	SVG string `json:"svg"`
	Alt string `json:"alt"`
}

type Currency struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type CountryDTO struct {
	Slug         string  `json:"slug"`
	Name         string  `json:"name"`
	OfficialName string  `json:"official_name"`
	Capital      string  `json:"capital"`
	Region       string  `json:"region"`
	Subregion    string  `json:"subregion"`
	Population   int64   `json:"population"`
	FlagURL      string  `json:"flag_url"`
	FlagAlt      string  `json:"flag_alt"`
	Currencies   string  `json:"currencies"`
	Languages    string  `json:"languages"`
	CCA2         string  `json:"cca2"`
	CCA3         string  `json:"cca3"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
}
