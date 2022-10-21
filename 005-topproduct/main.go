package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"text/tabwriter"
)

type Product struct {
	Name   string `json:"product"`
	Price  int    `json:"price"`
	Rating int    `json:"rating"`
}

type Metrics struct {
	topPrice  []*Product
	topRating []*Product
	maxPrice  int
	maxRating int
}

func (m *Metrics) addTop(product *Product) {
	if product.Price >= m.maxPrice {
		if product.Price > m.maxPrice {
			m.topPrice = m.topPrice[:0]
			m.maxPrice = product.Price
		}
		m.topPrice = append(m.topPrice, product)
	}
	if product.Rating >= m.maxRating {
		if product.Rating > m.maxRating {
			m.topRating = m.topRating[:0]
			m.maxRating = product.Rating
		}
		m.topRating = append(m.topRating, product)
	}
}

func main() {
	fileName := flag.String("f", "db.json", "path to the file")
	flag.Parse()

	var metrics Metrics

	if filepath.Ext(*fileName) == ".json" {
		readJSON(*fileName, &metrics)
	}

	if filepath.Ext(*fileName) == ".csv" {
		readCSV(*fileName, &metrics)
	}

	printTop("Top Price:", metrics.topPrice)
	fmt.Print("\n") //nolint: forbidigo
	printTop("Top Rating:", metrics.topRating)
}

func readJSON(fileName string, metrics *Metrics) {
	inputFile, err := os.Open(fileName)
	if err != nil {
		log.Printf("TopProduct: %v\n", err)
		return
	}
	defer inputFile.Close()

	decoderJSON := json.NewDecoder(inputFile)
	for {
		var product Product
		if err := decoderJSON.Decode(&product); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			log.Printf("TopProduct: %v\n", err)
			return
		}

		metrics.addTop(&product)
	}
}

func readCSV(fileName string, metrics *Metrics) {
	inputFile, err := os.Open(fileName)
	if err != nil {
		log.Printf("TopProduct: %v\n", err)
		return
	}
	defer inputFile.Close()

	reader := csv.NewReader(inputFile)
	reader.FieldsPerRecord = 3
	for {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		// Skip header
		if record[0] == "Product" {
			continue
		}
		price, err := strconv.Atoi(record[1])
		if err != nil {
			log.Printf("TopProduct: price: convert: %s", record[1])
			return
		}
		rating, err := strconv.Atoi(record[2])
		if err != nil {
			log.Printf("TopProduct: rating: convert: %s", record[2])
			return
		}

		metrics.addTop(&Product{record[0], price, rating})
	}
	if err != nil && !errors.Is(err, io.EOF) {
		log.Printf("TopProduct: %v", err)
	}
}

func printTop(header string, slTop []*Product) {
	fmt.Println(header) //nolint: forbidigo
	const formatHeader = "%s\t%s\t%s\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, formatHeader, "Name", "Price", "Rating")
	fmt.Fprintf(tw, formatHeader, "----", "-----", "------")

	for _, product := range slTop {
		fmt.Fprintf(tw, "%s\t%d\t%d\n", product.Name, product.Price, product.Rating)
	}

	tw.Flush()
}
