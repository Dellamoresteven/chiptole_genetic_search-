package main

import (
	//"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

// Store holds all the metadata for each store in the csv
type Store struct {
	name    string
	city    string
	address string
	lat     float64
	log     float64
}

const (
	numGames = 500
)

func main() {
	csvfile, err := os.Open("data/datasets_804019_1378604_chipotle_stores.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}
	var stores []Store

	r := csv.NewReader(csvfile)

	for {
		record, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		s := Store{
			name:    record[0],
			city:    record[1],
			address: record[2],
		}
		s.lat, _ = strconv.ParseFloat(record[3], 64)
		s.log, _ = strconv.ParseFloat(record[4], 64)

		stores = append(stores, s)

		startGame()
	}
}

func startGame() {
	for i := 0; i < numGames; i++ {

	}
}
