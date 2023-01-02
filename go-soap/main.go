package main

import (
	"encoding/xml"
	"log"
	"net/http"
	"time"

	"github.com/tiaguinho/gosoap"
)

// GetIPLocationResponse will hold the Soap response
type GetCountryCurrencyResponse struct {
	CountryCurrencyResponse string `xml:"CountryCurrencyResponse"`
}

// GetIPLocationResult will
type GetCountryCurrencyResult struct {
	XMLName  xml.Name
	SISOCode string `xml:"sISOCode"`
	SName    string `xml:"sName"`
}

type CurrencyResponse struct {
	XMLName                    xml.Name
	GetCountryCurrencyResponse struct {
		XMLName               xml.Name
		CountryCurrencyResult struct {
			XMLName  xml.Name
			SISOCode string `xml:"sISOCode"`
			SName    string `xml:"sName"`
		} `xml:"CountryCurrencyResult"`
	} `xml:"GetCountryCurrencyResponse"`
}

var r GetCountryCurrencyResponse

func main() {
	httpClient := &http.Client{
		Timeout: 1500 * time.Millisecond,
	}
	// set custom envelope
	gosoap.SetCustomEnvelope("soapenv", map[string]string{
		"xmlns:soapenv": "http://schemas.xmlsoap.org/soap/envelope/",
		"xmlns:tem":     "http://tempuri.org/",
	})

	soap, err := gosoap.SoapClient("http://webservices.oorsprong.org/websamples.countryinfo/CountryInfoService.wso?WSDL", httpClient)
	if err != nil {
		log.Fatalf("SoapClient error: %s", err)
	}

	// Use gosoap.ArrayParams to support fixed position params
	params := gosoap.Params{
		"sCountryISOCode": "INDIA",
	}

	res, err := soap.Call("CountryCurrency", params)
	if err != nil {
		log.Fatalf("Call error: %s", err)
	}
	response := new(CurrencyResponse)
	//err = res.Unmarshal(&r)
	//if err != nil {
	//	return
	//}

	// GetIpLocationResult will be a string. We need to parse it to XML
	err = xml.Unmarshal(res.Body, &response)
	if err != nil {
		log.Fatalf("xml.Unmarshal error: %s", err)
	}

	log.Println("Country Name: ", response.GetCountryCurrencyResponse.CountryCurrencyResult.SName)
	log.Println("Code: ", response.GetCountryCurrencyResponse.CountryCurrencyResult.SISOCode)
}
