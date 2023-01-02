package main

import (
	"countryinfoservice/gen"
	"fmt"
	"github.com/hooklift/gowsdl/soap"
	"log"
	"net/http"
	"time"
)

//	func main() {
//		client := soap.NewClient("http://127.0.0.1:8000")
//		service := gen.NewCountryInfoServiceSoapType(client)
//		reply, err := service.CapitalCity(&gen.CapitalCity{
//			SCountryISOCode: "INDIA",
//		})
//		if err != nil {
//			log.Fatalf("could't get trade prices: %v", err)
//		}
//		log.Println(reply)
//	}
var done = make(chan struct{})

func client() {
	client := soap.NewClient("http://127.0.0.1:8000")
	service := gen.NewCountryInfoServiceSoapType(client)
	reply, err := service.CapitalCity(&gen.CapitalCity{
		SCountryISOCode: "INDIA",
	})
	if err != nil {
		log.Fatalf("could't get trade prices: %v", err)
	}
	fmt.Println(reply)
	done <- struct{}{}
}

// use fixtures/test.wsdl
func main() {
	http.HandleFunc("/", gen.Endpoint)
	go func() {
		time.Sleep(time.Second * 1)
		client()
	}()
	go func() {
		http.ListenAndServe(":8000", nil)
	}()
	<-done
}
