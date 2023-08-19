package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

func proxyHandler(wr http.ResponseWriter, rq *http.Request) {
	// Target URL needs to come from somewhere else
	targetUrl := "https://google.com"

	newRequest, err := http.NewRequest(rq.Method, targetUrl, rq.Body)
	if err != nil {
		http.Error(wr, "error creating proxy request", http.StatusInternalServerError)
		return
	}
	// Copy the headers to the new request
	for header, headerValues := range rq.Header {
		for _, value := range headerValues {
			newRequest.Header.Add(header, value)
		}
	}

	// Creating new client and sending request
	client := http.Client{}
	resp, err := client.Do(newRequest)
	if err != nil {
		http.Error(wr, "Error sending proxy request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy remote response headers to the proxy response
	for name, values := range resp.Header {
		for _, value := range values {
			wr.Header().Add(name, value)
		}
	}
	// Set proper status code
	wr.WriteHeader(resp.StatusCode)
	// Copy remote response body to the proxy response
	io.Copy(wr, resp.Body)
}

func StartProxy(port uint16) {
	http.HandleFunc("/", proxyHandler)

	err := startServer(port)
	if err != nil {
		log.Fatal(err)
	}
}

func startServer(port uint16) error {
	strPort := ":" + strconv.FormatUint(uint64(port), 10)
	err := http.ListenAndServe(strPort, nil)
	if err == nil {
		fmt.Printf("Starting server on port %s\n", strPort)
	}
	return err
}
