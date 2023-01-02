package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

var getTemplate = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <CountryCurrency xmlns="http://www.oorsprong.org/websamples.countryinfo">
      <sCountryISOCode>{{.ISOCode}}</sCountryISOCode>
    </CountryCurrency>
  </soap:Body>
</soap:Envelope>`

type Request struct {
	ISOCode string
}

func populateRequest() *Request {
	req := Request{}
	req.ISOCode = "INDIA"
	return &req
}

type Response struct {
	XMLName  xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	SoapBody *SOAPBodyResponse
}
type SOAPBodyResponse struct {
	XMLName xml.Name `xml:"Body"`
	Resp    *ResponseBody
}

type ResponseBody struct {
	XMLName  xml.Name `xml:"Body"`
	SISOCode string   `xml:"sISOCode"`
	SName    string   `xml:"sName"`
	Status   string   `xml:"Status"`
}

func generateSOAPRequest(req *Request) (*http.Request, error) {
	// Using the var getTemplate to construct request
	temp, err := template.New("InputRequest").Parse(getTemplate)
	if err != nil {
		fmt.Printf("Error while marshling object. %s \n", err.Error())
		return nil, err
	}

	doc := &bytes.Buffer{}
	// Replacing the doc from temp with actual req values
	err = temp.Execute(doc, req)
	if err != nil {
		fmt.Printf("temp.Execute error. %s \n", err.Error())
		return nil, err
	}

	buffer := &bytes.Buffer{}
	encoder := xml.NewEncoder(buffer)
	err = encoder.Encode(doc.String())
	if err != nil {
		fmt.Printf("encoder.Encode error. %s \n", err.Error())
		return nil, err
	}

	r, err := http.NewRequest(http.MethodPost, "http://webservices.oorsprong.org/websamples.countryinfo/CountryInfoService.wso", bytes.NewBuffer([]byte(doc.String())))
	if err != nil {
		fmt.Printf("Error making a request. %s \n", err.Error())
		return nil, err
	}

	return r, nil
}

func soapCall(req *http.Request) (*Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r := &Response{}
	err = xml.Unmarshal(body, &r)

	if err != nil {
		return nil, err
	}

	if r.SoapBody.Resp.Status != "200" {
		return nil, err
	}

	return r, nil
}
func main() {
	request := populateRequest()
	client, err := generateSOAPRequest(request)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	call, err := soapCall(client)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println(call)
}
