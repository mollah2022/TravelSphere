package models

type WeatherResponse struct {
	Location WeatherLocation `json:"location"`
	Current  WeatherCurrent  `json:"current"`
}

type WeatherLocation struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

type WeatherCurrent struct {
	TempC      float64          `json:"temp_c"`
	TempF      float64          `json:"temp_f"`
	Condition  WeatherCondition `json:"condition"`
	Humidity   int              `json:"humidity"`
	WindKph    float64          `json:"wind_kph"`
	FeelsLikeC float64          `json:"feelslike_c"`
}

type WeatherCondition struct {
	Text string `json:"text"`
	Icon string `json:"icon"`
}

type WeatherDTO struct {
	Location   string  `json:"location"`
	Country    string  `json:"country"`
	TempC      float64 `json:"temp_c"`
	Condition  string  `json:"condition"`
	Icon       string  `json:"icon"`
	Humidity   int     `json:"humidity"`
	WindKph    float64 `json:"wind_kph"`
	FeelsLikeC float64 `json:"feels_like_c"`
	Available  bool    `json:"available"`
}
