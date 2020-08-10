package main

/**
 * @author Steven Dellamore
 * @email dellamoresteven@gmail.com
 */

/*
Notes:
Init Travelers
✓	- Variable length dna
✓	- Randomly selected from the Stores list
✓	- Make sure they all start at the same store
Playing the game
✓	- Need to use go routine to speed this up.
✓	- Count the number of unique stores the Travel has visited
✓	- Track the distance the traveler has traveled
	- Create a score algorithm to rate how well a travel did
✓		- Need to visit every unique store at least 1
		- Travel the least amount of distance
*/

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
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
	dna       []Store        // List of stores to go to in order
	numStores map[Store]bool // A set of stores visited. This set is unique
	distance  float64        // Distance traveled
	score     float64        // Overall score of the traveler
}

func newTraveler(dna []Store) Traveler {
	var t Traveler
	t.numStores = make(map[Store]bool)
	t.distance = 0
	t.score = 0
	t.dna = dna
	return t
}

const (
	numGames     = 100 // Number of games to be played
	population   = 100 // Number of Travelers
	mutationRate = .05 // Mutation rate of the dna of travelers
)

var travelers []Traveler // List of Travelers that will take on the chiptole world
var stores []Store       // The metadata of the stores read from the csv
var gameNumber int = 0

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

		s := Store{
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
	initPopulation()
	var wg sync.WaitGroup
	for i := 0; i < numGames; i++ {
		gameNumber++
		// fmt.Println("\n\nNEW GAME STARTED\n\n")
		for j := 0; j < population; j++ {
			wg.Add(1)
			go playGame(&travelers[j], &wg)
		}
		wg.Wait()
		breedNextGeneration()
	}

	printBestTraveler()
}

func initPopulation() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	for i := 0; i < population; i++ {
		var dna []Store
		dna = append(dna, stores[0])
		for j := 1; j < r1.Intn(len(stores)-1000)+1000; j++ {
			dna = append(dna, stores[r1.Intn(len(stores))])
		}
		t := newTraveler(dna)
		travelers = append(travelers, t)
	}
}

func playGame(t *Traveler, wg *sync.WaitGroup) {
	defer wg.Done()
	oldStore := t.dna[0]
	for i := 1; i < len(t.dna); i++ {
		newStore := t.dna[i]
		t.distance += math.Abs((math.Sqrt(
			math.Pow(newStore.log-oldStore.log, 2) +
				math.Pow(newStore.lat-oldStore.lat, 2))))
		t.numStores[newStore] = true
		oldStore = newStore
	}
	t.score = float64(len(t.numStores)) + (1 / t.distance)
}

func breedNextGeneration() {
	sort.SliceStable(travelers, func(i, j int) bool {
		return travelers[i].score > travelers[j].score
	})

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	travelers = travelers[:len(travelers)-3] // remove the worst 3 performing travelers
	travelers = append(travelers, breed(travelers[r1.Intn(10)+1], travelers[r1.Intn(10)+1]))
	travelers = append(travelers, breed(travelers[r1.Intn(10)+1], travelers[r1.Intn(10)+1]))
	travelers = append(travelers, breed(travelers[r1.Intn(10)+1], travelers[r1.Intn(10)+1]))

	maxScore := 0.0
	for _, t := range travelers {
		// fmt.Println(t.score)
		if maxScore < t.score {
			maxScore = t.score
		}
		t.score = 0
		t.numStores = make(map[Store]bool)
		t.distance = 0
	}
	fmt.Println("Max Score for game ", gameNumber, " : ", maxScore)
}

func breed(p1, p2 Traveler) (child Traveler) {
	var dna []Store

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	if mutationRate > float64((r1.Intn(100) / 100)) {
		dna = append(dna, stores[0])
		for j := 1; j < r1.Intn(len(stores)-1000)+1000; j++ {
			dna = append(dna, stores[r1.Intn(len(stores))])
		}
		child = newTraveler(dna)
	}

	dna = append(dna, p1.dna[0:len(p1.dna)/2]...)
	dna = append(dna, p2.dna[len(p2.dna)/2:len(p2.dna)]...)
	// fmt.Println("p1: ", p1.dna[1:2])
	// fmt.Println("p2: ", p2.dna[len(p2.dna)-1:])
	// fmt.Println("CHILD START: ", dna[1:2])
	// fmt.Println("CHILD END: ", dna[len(dna)-1:], "\n\n\n")

	child = newTraveler(dna)
	return
}

func printBestTraveler() {
	sort.SliceStable(travelers, func(i, j int) bool {
		return travelers[i].score > travelers[j].score
	})
	t := travelers[0]
	fmt.Println("\n\nBEST")
	fmt.Println("Score: ", t.score)
	fmt.Println("numStores: ", len(t.numStores))
	fmt.Println("distance: ", t.distance)
}
