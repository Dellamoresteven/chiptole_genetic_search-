package main

/*
Notes:
Init Travelers
✓	- Varible length dna
✓	- Randomly selected from the Stores list
Playing the game
✓	- Need to use go routies to speed this up.
	- 
*/

import (
	"fmt"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"math/rand"
)

// Store holds all the metadata for each store in the csv
type Store struct {
	name    string
	city    string
	address string
	lat     float64
	log     float64
}

type Traveler struct {
	dna 	  []Store // List of stores to go to in order
	numStores int     // Number of stores the traveler has been too
	distance  float64 // Distance traveled 
	score	  float64 // Overall score of the traveler
}

const (
	numGames = 500     // Number of games to be played
	population = 100   // Number of Travelers
	mutationRate = .05 // Mutation rate of the dna of travelers
)

var travlers []Traveler // List of Travelers that will take on the chiptole world
var stores   []Store    // The metadata of the stores read from the csv

func main() {
	csvfile, err := os.Open("data/datasets_804019_1378604_chipotle_stores.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	r := csv.NewReader(csvfile)

	for {
		record, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		s := Store {
			name:    record[0],
			city:    record[1],
			address: record[2],
		}
		s.lat, _ = strconv.ParseFloat(record[3], 64)
		s.log, _ = strconv.ParseFloat(record[4], 64)

		stores = append(stores, s)
	}
	startGame()
}

func startGame() {
	initPopulation();
	var wg sync.WaitGroup
	fmt.Println(len(stores))
	for i := 0; i < numGames; i++ {
		for j := 0; j < population; j++ {
			wg.Add(1)
			go playGame(travlers[j], &wg)
		}
		wg.Wait()
		// breedNextGeneration()
	}
}

func initPopulation() {
	for i := 0; i < population; i++ {
		var dna []Store
		for j := 0; j < rand.Intn(len(stores) - 1000) + 1000; j++ {
			dna = append(dna, stores[rand.Intn(len(stores))])
		}
		t := Traveler {
			dna: dna,
			numStores: 0,
			distance: 0,
			score: 0,
		}
		travlers = append(travlers, t)
	}
}

func playGame(t Traveler, wg *sync.WaitGroup) {

	defer wg.Done()

}
