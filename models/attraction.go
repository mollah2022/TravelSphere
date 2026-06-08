package models

type AttractionResponse struct {
	Features []AttractionFeature `json:"features"`
}

type AttractionFeature struct {
	ID         string               `json:"id"`
	Name       string               `json:"name"`
	Geometry   AttractionGeometry   `json:"geometry"`
	Properties AttractionProperties `json:"properties"`
}

type AttractionGeometry struct {
	Coordinates []float64 `json:"coordinates"`
}

type AttractionProperties struct {
	XID   string          `json:"xid"`
	Name  string          `json:"name"`
	Dist  float64         `json:"dist"`
	Rate  int             `json:"rate"`
	OSM   string          `json:"osm"`
	Kinds string          `json:"kinds"`
	Point AttractionPoint `json:"point"`
}

type AttractionPoint struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type AttractionDTO struct {
	XID       string   `json:"xid"`
	Name      string   `json:"name"`
	Kinds     string   `json:"kinds"`
	KindsList []string `json:"kinds_list"`
	Distance  float64  `json:"distance"`
	Latitude  float64  `json:"latitude"`
	Longitude float64  `json:"longitude"`
}

type PopularAttraction struct {
	Name    string `json:"name"`
	Kinds   string `json:"kinds"`
	Country string `json:"country"`
}
