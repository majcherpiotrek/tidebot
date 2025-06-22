package worldtides

import "time"

type WorldTidesResponse struct {
	Status        int       `json:"status"`
	Error         string    `json:"error,omitempty"`
	RequestDatum  string    `json:"requestDatum,omitempty"`
	ResponseDatum string    `json:"responseDatum,omitempty"`
	RequestLat    float64   `json:"requestLat,omitempty"`
	RequestLon    float64   `json:"requestLon,omitempty"`
	ResponseLat   float64   `json:"responseLat,omitempty"`
	ResponseLon   float64   `json:"responseLon,omitempty"`
	CallCount     int       `json:"callCount,omitempty"`
	Atlas         string    `json:"atlas,omitempty"`
	Station       string    `json:"station,omitempty"`
	Copyright     string    `json:"copyright,omitempty"`
	Timezone      string    `json:"timezone,omitempty"`
	Heights       []Height  `json:"heights,omitempty"`
	Extremes      []Extreme `json:"extremes,omitempty"`
	Datums        []Datum   `json:"datums,omitempty"`
	Stations      []Station `json:"stations,omitempty"`
}

type Height struct {
	Dt     int64   `json:"dt"`
	Date   string  `json:"date"`
	Height float64 `json:"height"`
}

type Extreme struct {
	Dt     int64   `json:"dt"`
	Date   string  `json:"date"`
	Height float64 `json:"height"`
	Type   string  `json:"type"`
}

type Datum struct {
	Name   string  `json:"name"`
	Height float64 `json:"height"`
}

type Station struct {
	ID   string  `json:"id"`
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

// Helper methods for Height
func (h *Height) Time() time.Time {
	return time.Unix(h.Dt, 0)
}

// Helper methods for Extreme
func (e *Extreme) Time() time.Time {
	return time.Unix(e.Dt, 0)
}

func (e *Extreme) IsHighTide() bool {
	return e.Type == "High"
}

func (e *Extreme) IsLowTide() bool {
	return e.Type == "Low"
}
