package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

// Create a cache with a default expiration time of 5 minutes, and which
// purges expired items every 10 minutes
var kc = cache.New(5*time.Minute, 10*time.Minute)

func main() {
	// initialize searcher
	searcher := &Searcher{}
	err := searcher.Load("data.gz")
	if err != nil {
		log.Fatalf("unable to load search data due: %v", err)
	}

	// define port, we need to set it as env for Heroku deployment
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	// define http handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleMain)
	mux.HandleFunc("/search", handleSearch(searcher))

	errorHandler := &ErrorHandler{
		Handler: mux,
	}

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: errorHandler,
	}

	// start server
	fmt.Printf("Server is listening on %s... \n", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("unable to start server due: %v \n", err)
	}
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	// panic("handleMain")
	http.ServeFile(w, r, "./static/index.html")
}

func handleSearch(s *Searcher) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			//panic("handleSearch")

			// fetch query string from query params
			q := r.URL.Query().Get("q")

			//Get Cache
			foo, found := kc.Get(q)
			if found {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, foo)
				return
			}

			if len(q) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("missing search query in query params"))
				return
			}
			// search relevant records
			records, err := s.Search(q)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			// output success response
			encjson, error := json.Marshal(records)
			if error != nil {
				panic(error.Error())
			}

			json := string(encjson)
			//Set Cache
			kc.Set(q, json, cache.DefaultExpiration)

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, json)
		},
	)
}

type Searcher struct {
	records []Record
}

func (s *Searcher) Load(filepath string) error {
	// open file
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("unable to open source file due: %v", err)
	}
	defer file.Close()
	// read as gzip
	reader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("unable to initialize gzip reader due: %v", err)
	}
	// read the reader using scanner to contstruct records
	var records []Record
	cs := bufio.NewScanner(reader)
	for cs.Scan() {
		var r Record
		err = json.Unmarshal(cs.Bytes(), &r)
		if err != nil {
			continue
		}
		records = append(records, r)
	}
	s.records = records

	return nil
}

func (s *Searcher) Search(query string) ([]Record, error) {
	var result []Record
	word := strings.Fields(strings.ToLower(query))
	for _, record := range s.records {
		if contains(word, strings.ToLower(record.Title)) || contains(word, strings.ToLower(record.Content)) {
			result = append(result, record)
		}
	}
	return result, nil
}

type Record struct {
	ID        int64    `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	ThumbURL  string   `json:"thumb_url"`
	Tags      []string `json:"tags"`
	UpdatedAt int64    `json:"updated_at"`
	ImageURLs []string `json:"image_urls"`
}

func contains(s []string, str string) bool {
	k := strings.Fields(str)
	wordTrue := false
	for _, v := range s {
		wordTrue = false
		for _, i := range k {
			if i == v {
				wordTrue = true
				break
			}
		}
		if !wordTrue {
			return false
		}
	}

	return true
}

type ErrorHandler struct {
	Handler http.Handler
}

func (errorHandler *ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %s", err)
			fmt.Printf("Error: %s", err)
		}
	}()
	errorHandler.Handler.ServeHTTP(w, r)
}
