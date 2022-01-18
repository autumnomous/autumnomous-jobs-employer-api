package zipcode_test

import (
	"autumnomous-jobs-employer-api/shared/services/zipcode"
	"log"
	"os"
	"testing"

	// "57channels.io/geocode/zipcode"

	"github.com/stretchr/testify/assert"
)

var apiKey = os.Getenv("api_key")
var testZip = "45505"
var chars = "Pittsbu"
var latitude = 33.4499
var longitude = -112.0712
var zipcode1 = "85258"
var zipcode2 = "85004"

func TestGetZipCode(t *testing.T) {

	assert := assert.New(t)
	var gateway = zipcode.NewZipCodeGateway(apiKey)
	zip, zipErr := gateway.GetZipCode(testZip)

	assert.Nil(zipErr)
	assert.NotNil(zip)

	assert.Equal(zip.City, "Springfield", "cities should be the same")
	assert.Equal(zip.State, "OH", "states should be the same")
	assert.Equal(zip.ZipCode, testZip, "IP Addresses should be equal")
}

func TestGetZipCode2(t *testing.T) {

	assert := assert.New(t)
	var gateway = zipcode.NewZipCodeGateway(apiKey)

	testZip = "15218"
	zip, zipErr := gateway.GetZipCode(testZip)

	assert.Nil(zipErr)
	assert.NotNil(zip)

	assert.Equal(zip.City, "Pittsburgh", "cities should be the same")
	assert.Equal(zip.State, "PA", "states should be the same")
	assert.Equal(zip.ZipCode, testZip, "zip codes should be equal")
}

func TestGetZipCodeAutoComplete(t *testing.T) {

	assert := assert.New(t)
	var gateway = zipcode.NewZipCodeGateway(apiKey)
	locations, zipErr := gateway.GetAutoComplete(chars)

	assert.Nil(zipErr)
	assert.True(len(locations) > 0)
}

func TestGetZipCodeJSAutoComplete(t *testing.T) {

	assert := assert.New(t)
	var gateway = zipcode.NewZipCodeGateway(apiKey)
	locations, zipErr := gateway.GetJSAutoComplete(chars)

	assert.Nil(zipErr)
	assert.True(len(locations) > 0)
}

func TestGetLocationByLatLong(t *testing.T) {
	assert := assert.New(t)

	var gateway = zipcode.NewZipCodeGateway(apiKey)
	location, locationErr := gateway.GetLocationByLatLong(longitude, latitude)

	assert.Nil(locationErr)
	assert.NotNil(location)
	assert.Equal(location.ZipCode, "85004", "Zip Codes should be the same.")
	assert.Equal(location.City, "Phoenix", "Cities should be the same.")
	assert.Equal(location.State, "AZ", "States should be the same.")
	assert.Equal(location.Country, "US", "Countries should be the same.")
	// assert.Equal(location.AreaCode, "602", "area codes should be the same")
}

func TestGetDistanceBetweenZipcodes(t *testing.T) {
	assert := assert.New(t)

	var gateway = zipcode.NewZipCodeGateway(apiKey)
	result, distanceErr := gateway.GetDistanceBetweenZipCodes(zipcode1, zipcode2)

	assert.Nil(distanceErr)
	assert.NotNil(result)
	assert.True(result.ZipCode1 == zipcode1)
	assert.True(result.ZipCode2 == zipcode2)
	assert.True(result.DistanceInKilometers > 0)
	assert.True(result.DistanceInMiles > 0)
}

func TestGetZipCodesInRadius(t *testing.T) {
	assert := assert.New(t)

	var gateway = zipcode.NewZipCodeGateway(apiKey)
	zips, zipsErr := gateway.GetZipCodesInRadius("43403", 15.5)

	log.Println(zips)
	log.Print(zipsErr)

	//assert.Nil(zipsErr)

	var zipCount = len(zips)
	assert.True(zipCount > 0)
}
