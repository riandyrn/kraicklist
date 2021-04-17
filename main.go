package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/ahmetb/go-linq/v3"
	"github.com/sahilm/fuzzy"
)

func main() {
	// initialize searcher
	searcher := &Searcher{}
	err := searcher.Load("data.gz")
	if err != nil {
		log.Fatalf("unable to load search data due: %v", err)
	}

	// define http handlers
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/search", handleSearch(searcher))

	// define port, we need to set it as env for Heroku deployment
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	// start server
	fmt.Printf("Server is listening on %s...", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatalf("unable to start server due: %v", err)
	}
}

func handleSearch(s *Searcher) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// fetch query string from query params
			q := r.URL.Query().Get("q")
			if len(q) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("missing search query in query params"))
				return
			}

			// simple cache
			if q != s.prevQuery {
				// search relevant records
				records, err := s.FuzzySearch(q)
				s.prevRecords = records
				s.prevQuery = q

				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(err.Error()))
					return
				}
			}

			// output success response
			buf := new(bytes.Buffer)
			encoder := json.NewEncoder(buf)
			res := map[string]interface{}{}

			pageQuery := r.URL.Query().Get("page")
			page, err := strconv.Atoi(pageQuery)
			if err != nil || page <= 0 {
				page = 1
			}

			sizeQuery := r.URL.Query().Get("size")
			size, err := strconv.Atoi(sizeQuery)
			if err != nil {
				size = 10
			}

			res["data"] = linq.From(s.prevRecords).Skip((page - 1) * size).Take((size)).Results()
			res["page"] = page
			res["size"] = size
			res["totalPage"] = math.Ceil(float64(s.prevRecords.Len()) / float64(size))
			res["total"] = s.prevRecords.Len()

			encoder.Encode(res)
			w.Header().Set("Content-Type", "application/json")
			w.Write(buf.Bytes())
		},
	)
}

type Searcher struct {
	records     Records
	prevRecords Records
	prevQuery   string
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
	var records Records
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

func (s *Searcher) FuzzySearch(query string) (Records, error) {
	var indexScore []interface{}
	var result Records

	fuzzyResults := fuzzy.FindFrom(query, s.records)

	for _, r := range fuzzyResults {
		indexScore = append(indexScore, KV{r.Index, r.Score})
	}

	sort.Slice(indexScore, func(i, j int) bool {
		return indexScore[i].(KV).Value > indexScore[j].(KV).Value
	})

	// limit the score relative to the first entry
	indexScore = linq.From(indexScore).Where(func(r interface{}) bool {
		return r.(KV).Value >= (indexScore[0].(KV).Value - 25)
	}).Results()

	for _, k := range indexScore {
		result = append(result, s.records[k.(KV).Key])
	}

	return result, nil
}

type KV struct {
	Key   int
	Value int
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

type Records []Record

func (r Records) String(i int) string {
	return r[i].Title + " " + r[i].Content
}

func (r Records) Len() int {
	return len(r)
}
