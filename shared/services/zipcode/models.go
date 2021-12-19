package zipcode

type AutoCompleteRequest struct {
	Chars string `json:"chars"`
}

type CityAutoCompleteResponse struct {
	City      string  `json:"city"`
	State     string  `json:"state"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type ZipCodeRequest struct {
	ZipCode string `json:"zip_code"`
}

type ZipCodeResponse struct {
	ZipCode   string  `json:"zip_code"`
	City      string  `json:"city"`
	State     string  `json:"state"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	County    string  `json:"county"`
}

type LocationByLatLongRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type ZipCodeResponseWithDistance struct {
	ZipCode                  string  `json:"zip_code"`
	City                     string  `json:"city"`
	State                    string  `json:"state"`
	Country                  string  `json:"country"`
	Latitude                 float64 `json:"latitude"`
	Longitude                float64 `json:"longitude"`
	County                   string  `json:"county"`
	DistanceAwayInMiles      float64 `json:"distance_away_in_miles"`
	DistanceAwayInKilometers float64 `json:"distance_away_in_kilometers"`
}

type ZipCodesDistanceRequest struct {
	ZipCode1 string `json:"zip_code_1"`
	ZipCode2 string `json:"zip_code_2"`
}

type ZipCodesDistanceResponse struct {
	ZipCode1             string  `json:"zip_code1"`
	ZipCode2             string  `json:"zip_code2"`
	DistanceInMiles      float64 `json:"distance_miles"`
	DistanceInKilometers float64 `json:"distance_km"`
}

type ZipCodesInRadiusRequest struct {
	ZipCode string  `json:"zip_code"`
	Radius  float64 `json:"radius"`
}

type CityLatLongAutoCompleteResponse struct {
	Location string    `json:"label"`
	Point    CityPoint `json:"value"`
}

type CityPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
