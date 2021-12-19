package zipcode

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	// "57channels.io/geocode/geoerror"
)

type GeoCodeError struct {
	Message string `json:"error"`
}

type ZipCodeGateway struct {
	apiKey string
}

func NewZipCodeGateway(apiKey string) *ZipCodeGateway {
	return &ZipCodeGateway{apiKey}
}

func (gateway *ZipCodeGateway) GetZipCode(zip string) (*ZipCodeResponse, error) {

	var zipcode ZipCodeResponse
	var zipRequest = ZipCodeRequest{ZipCode: zip}

	var uri = "https://api.zipcodeservices.io/v1/zipcode"

	jsonRequest, marshalErr := json.Marshal(zipRequest)

	if marshalErr != nil {
		return &zipcode, marshalErr
	}

	request, requestErr := http.NewRequest("POST", uri, bytes.NewBuffer(jsonRequest))

	if requestErr != nil {
		return &zipcode, requestErr
	}

	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(gateway.apiKey, "")
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, responseErr := client.Do(request)

	if responseErr != nil {
		return &zipcode, responseErr
	}

	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		var errorString GeoCodeError
		json.Unmarshal(body, &errorString)
		log.Println(errorString.Message)
		return &zipcode, errors.New(errorString.Message)
	}

	json.Unmarshal(body, &zipcode)

	return &zipcode, nil
}

func (gateway *ZipCodeGateway) GetAutoComplete(chars string) ([]CityAutoCompleteResponse, error) {

	var autocompleteResponse []CityAutoCompleteResponse
	var autoCompleteRequest = AutoCompleteRequest{Chars: chars}

	var uri = "https://api.zipcodeservices.io/v1/autocomplete"

	jsonRequest, marshalErr := json.Marshal(autoCompleteRequest)

	if marshalErr != nil {
		return autocompleteResponse, marshalErr
	}

	request, requestErr := http.NewRequest("POST", uri, bytes.NewBuffer(jsonRequest))

	if requestErr != nil {
		return autocompleteResponse, requestErr
	}

	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(gateway.apiKey, "")
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, responseErr := client.Do(request)

	if responseErr != nil {
		return autocompleteResponse, responseErr
	}

	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		var errorString GeoCodeError
		json.Unmarshal(body, &errorString)
		log.Println(errorString.Message)
		return autocompleteResponse, errors.New(errorString.Message)
	}

	json.Unmarshal(body, &autocompleteResponse)

	return autocompleteResponse, nil
}

func (gateway *ZipCodeGateway) GetJSAutoComplete(chars string) ([]CityLatLongAutoCompleteResponse, error) {

	var autocompleteResponse []CityLatLongAutoCompleteResponse
	var autoCompleteRequest = AutoCompleteRequest{Chars: chars}

	var uri = "https://api.zipcodeservices.io/v1/jsautocomplete"

	jsonRequest, marshalErr := json.Marshal(autoCompleteRequest)

	if marshalErr != nil {
		return autocompleteResponse, marshalErr
	}

	request, requestErr := http.NewRequest("POST", uri, bytes.NewBuffer(jsonRequest))

	if requestErr != nil {
		return autocompleteResponse, requestErr
	}

	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(gateway.apiKey, "")
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, responseErr := client.Do(request)

	if responseErr != nil {
		return autocompleteResponse, responseErr
	}

	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		var errorString GeoCodeError
		json.Unmarshal(body, &errorString)
		log.Println(errorString.Message)
		return autocompleteResponse, errors.New(errorString.Message)
	}

	json.Unmarshal(body, &autocompleteResponse)

	return autocompleteResponse, nil
}

func (gateway *ZipCodeGateway) GetLocationByLatLong(longitude float64, latitude float64) (*ZipCodeResponse, error) {
	var zipcode ZipCodeResponse
	var locationRequest = LocationByLatLongRequest{Latitude: latitude, Longitude: longitude}

	var uri = "https://api.zipcodeservices.io/v1/location"

	jsonRequest, marshalErr := json.Marshal(locationRequest)

	if marshalErr != nil {
		return &zipcode, marshalErr
	}

	request, requestErr := http.NewRequest("POST", uri, bytes.NewBuffer(jsonRequest))

	if requestErr != nil {
		return &zipcode, requestErr
	}

	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(gateway.apiKey, "")
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, responseErr := client.Do(request)

	if responseErr != nil {
		return &zipcode, responseErr
	}

	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		var errorString GeoCodeError
		json.Unmarshal(body, &errorString)
		log.Println(errorString.Message)
		return &zipcode, errors.New(errorString.Message)
	}

	json.Unmarshal(body, &zipcode)

	return &zipcode, nil
}

func (gateway *ZipCodeGateway) GetDistanceBetweenZipCodes(zipcode1 string, zipcode2 string) (*ZipCodesDistanceResponse, error) {
	var zipDistance ZipCodesDistanceResponse
	var zipDistanceRequest = ZipCodesDistanceRequest{ZipCode1: zipcode1, ZipCode2: zipcode2}

	var uri = "https://api.zipcodeservices.io/v1/zipcodedistance"

	jsonRequest, marshalErr := json.Marshal(zipDistanceRequest)

	if marshalErr != nil {
		return &zipDistance, marshalErr
	}

	request, requestErr := http.NewRequest("POST", uri, bytes.NewBuffer(jsonRequest))

	if requestErr != nil {
		return &zipDistance, requestErr
	}

	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(gateway.apiKey, "")
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, responseErr := client.Do(request)

	if responseErr != nil {
		return &zipDistance, responseErr
	}

	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		var errorString GeoCodeError
		json.Unmarshal(body, &errorString)
		log.Println(errorString.Message)
		return &zipDistance, errors.New(errorString.Message)
	}

	json.Unmarshal(body, &zipDistance)

	return &zipDistance, nil
}

func (gateway *ZipCodeGateway) GetZipCodesInRadius(zipcode string, radius float64) ([]ZipCodeResponseWithDistance, error) {

	var zipcodes []ZipCodeResponseWithDistance
	var radiusRequest = ZipCodesInRadiusRequest{ZipCode: zipcode, Radius: radius}

	var uri = "https://api.zipcodeservices.io/v1/zipsinradius"

	jsonRequest, marshalErr := json.Marshal(radiusRequest)

	if marshalErr != nil {
		return zipcodes, marshalErr
	}

	request, requestErr := http.NewRequest("POST", uri, bytes.NewBuffer(jsonRequest))

	if requestErr != nil {
		return zipcodes, requestErr
	}

	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(gateway.apiKey, "")
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, responseErr := client.Do(request)

	if responseErr != nil {
		return zipcodes, responseErr
	}

	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		var errorString GeoCodeError
		json.Unmarshal(body, &errorString)
		log.Println(errorString.Message)
		return zipcodes, errors.New(errorString.Message)
	}

	json.Unmarshal(body, &zipcodes)

	return zipcodes, nil
}
