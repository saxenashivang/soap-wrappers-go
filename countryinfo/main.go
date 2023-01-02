package main

import (
	"countryinfo/gen"
	"crypto/tls"
	"github.com/hooklift/gowsdl/soap"
	"log"
	"net/http"
	"time"
)

var done = make(chan struct{})

func CallClient() {
	client := soap.NewClient(
		"http://127.0.0.1:8000",
		soap.WithTimeout(time.Second*5),
		soap.WithBasicAuth("usr", "psw"),
		soap.WithTLS(&tls.Config{InsecureSkipVerify: true}),
	)
	service := gen.NewCountryInfoServiceSoapType(client)
	reply, err := service.CountryCurrency(&gen.CountryCurrency{SCountryISOCode: "INDIA"})
	if err != nil {
		log.Fatalf("could't get trade prices: %v", err)
	}
	log.Println(reply)
	done <- struct{}{}
}
func main() {
	//http.HandleFunc("/", gen.Endpoint)
	go func() {
		log.Fatal(http.ListenAndServe(":8000", nil))
	}()
	go func() {
		time.Sleep(time.Second * 1)
		CallClient()
	}()
	<-done
}
