package services

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"log"
	"os"

	"github.com/knightazura/contracts"
	"github.com/knightazura/data/model"
)

// A service to load and manage model that need
// to be formatted into common type, GeneralDocument
// Note: This service only for testing purpose
type Seeder struct {}

func InitSeeder() *Seeder {
	return &Seeder{}
}

// Load model data from file and format it to GeneralDocument
func (s *Seeder) LoadData(path string, model contracts.Indexer) (out model.GeneralDocuments){
	// Open file
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("unable to open source file due: %v", err)
		return
	}
	defer file.Close()

	// read as gzip
	reader, err := gzip.NewReader(file)
	if err != nil {
		log.Fatalf("unable to initialize gzip reader due: %v", err)
		return
	}

	// read the reader using scanner to contstruct records
	cs := bufio.NewScanner(reader)
	for cs.Scan() {
		err = json.Unmarshal(cs.Bytes(), &model)
		if err != nil {
			continue
		}

		// Convert any models into general document
		item, _ := model.ConvertToGeneralDocs()

		out = append(out, item)
	}

	return
}