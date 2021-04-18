package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// load data to sqlite
	loadData("data.gz", "data.db")

	// define http handlers
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/search", handleSearch())

	// define port, we need to set it as env for Heroku deployment
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	// start server
	fmt.Printf("Server is listening on %s...", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatalf("unable to start server due: %v", err)
	}
}

func loadData(filepath string, dbpath string) error {
	if _, err := os.Stat(dbpath); os.IsNotExist(err) {
		fmt.Println("Loading data...")

		db, err := sql.Open("sqlite3", dbpath)
		if err != nil {
			return fmt.Errorf("%q", err)
		}
		defer db.Close()

		sqlStatement := `
			CREATE TABLE record (
				id INTEGER NOT NULL PRIMARY KEY,
				title TEXT,
				content TEXT,
				thumb_url TEXT,
				tags TEXT,
				updated_at INTEGER,
				image_urls TEXT
			);
		`
		_, err = db.Exec(sqlStatement)
		if err != nil {
			return fmt.Errorf("%q: %s", err, sqlStatement)
		}

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

		cs := bufio.NewScanner(reader)
		for cs.Scan() {
			var r Record
			err = json.Unmarshal(cs.Bytes(), &r)
			if err != nil {
				continue
			}

			script := fmt.Sprintf(`INSERT INTO record(id, title, content, thumb_url, tags,updated_at,image_urls) values(
				%d,
				"%s",
				"%s",
				"%s",
				"%s",
				%d,
				"%s"
			);
			`,
				r.ID,
				r.Title,
				r.Content,
				r.ThumbURL,
				strings.Join(r.Tags[:], "\t"), // tsv formatted
				r.UpdatedAt,
				strings.Join(r.ImageURLs[:], "\t"), // tsv formatted
			)

			_, err = db.Exec(script)
			if err != nil {
				log.Fatal(err)
				continue
			}
		}
	}

	fmt.Println("Database is ready...")
	return nil
}

func handleSearch() http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// fetch query string from query params
			q := r.URL.Query().Get("q")
			if len(q) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("missing search query in query params"))
				return
			}

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

			cursorQuery := r.URL.Query().Get("cursor")
			cursor, err := strconv.Atoi(cursorQuery)
			if err != nil {
				cursor = -1
			}

			// output success response
			buf := new(bytes.Buffer)
			encoder := json.NewEncoder(buf)
			res := map[string]interface{}{}

			res["q"] = q
			res["page"] = page
			res["size"] = size
			res["cursor"] = cursor

			encoder.Encode(res)
			w.Header().Set("Content-Type", "application/json")
			w.Write(buf.Bytes())
		},
	)
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
