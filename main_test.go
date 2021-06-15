package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

func TestHandleSearch(t *testing.T) {
	searcher := &Searcher{}
	err := searcher.Load("data.gz")
	if err != nil {
		log.Fatalf("unable to load search data due: %v", err)
	}

	//enter search word here
	s := "iphone pro "

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3001/search?q="+url.QueryEscape(s), nil)
	recorder := httptest.NewRecorder()

	handleSearch(searcher)(recorder, request)

	response := recorder.Result()
	body, _ := io.ReadAll(response.Body)
	bodyString := string(body)
	var arr []string
	_ = json.Unmarshal([]byte(bodyString), &arr)

	//TOTAL RESULT
	res := len(arr)
	// ENTER THE RESULTS YOU WANT
	pas := 13

	fmt.Printf("Searching : %s \n", s)
	fmt.Printf("Result : %s \n", strconv.Itoa(res))
	fmt.Printf("Pass : %s \n", strconv.Itoa(pas))

	// CHECK RESULT
	if res != pas {
		t.FailNow()
	}
}
