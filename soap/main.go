package main

import (
	"fmt"
	"log"
)

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
