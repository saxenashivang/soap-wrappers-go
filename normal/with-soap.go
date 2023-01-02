package main

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type UserList struct {
	XMLName xml.Name
	Body    struct {
		XMLName                 xml.Name
		CountryCurrencyResponse struct {
			XMLName                 xml.Name
			CountryCurrencyResponse struct {
				SISOCode string `xml:"sISOCode"`
				SName    string `xml:"sName"`
			} `xml:"CountryCurrencyResult"`
		} `xml:"CountryCurrencyResponse"`
	}
}

func main() {
	// wsdl service url
	url := fmt.Sprintf("%s%s%s",
		"http://webservices.oorsprong.org",
		"/websamples",
		".countryinfo/CountryInfoService.wso",
	)

	// payload
	payload := []byte(strings.TrimSpace(`
		<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <CountryCurrency xmlns="http://www.oorsprong.org/websamples.countryinfo">
      <sCountryISOCode>US</sCountryISOCode>
    </CountryCurrency>
  </soap:Body>
</soap:Envelope>`,
	))

	httpMethod := "POST"

	// soap action
	//soapAction := "urn:listUsers"

	// authorization credentials
	username := "admin"
	password := "admin"

	log.Println("-> Preparing the request")

	// prepare the request
	req, err := http.NewRequest(httpMethod, url, bytes.NewReader(payload))
	if err != nil {
		log.Fatal("Error on creating request object. ", err.Error())
		return
	}

	// set the content type header, as well as the oter required headers
	req.Header.Set("Content-type", "text/xml")
	//req.Header.Set("SOAPAction", soapAction)
	req.SetBasicAuth(username, password)

	// prepare the client request
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	log.Println("-> Dispatching the request")

	// dispatch the request
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Error on dispatching request. ", err.Error())
		return
	}

	log.Println("-> Retrieving and parsing the response")

	// read and parse the response body
	result := new(UserList)
	err = xml.NewDecoder(res.Body).Decode(result)
	if err != nil {
		log.Fatal("Error on unmarshaling xml. ", err.Error())
		return
	}

	log.Println("-> Everything is good, printing users data")

	// print the users data
	iso := result.Body.CountryCurrencyResponse.CountryCurrencyResponse.SName
	name := result.Body.CountryCurrencyResponse.CountryCurrencyResponse.SISOCode
	fmt.Println(iso + " , " + name)
}
