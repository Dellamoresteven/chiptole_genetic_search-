package main

/**
 * @author Steven Dellamore
 * @email dellamoresteven@gmail.com
 */

/*
Notes:
Init Travelers
✓	- Varible length dna
✓	- Randomly selected from the Stores list
Playing the game
✓	- Need to use go routies to speed this up.
✓	- Count the number of unique stores the Travel has visited
✓	- Track the distance the traveler has traveled
	- Create a score algorthm to rate how well a travel did
		- Need to visit every unique store at least 1
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
	numGames     = 2   // Number of games to be played
	population   = 100 // Number of Travelers
	mutationRate = .05 // Mutation rate of the dna of travelers
)

var travlers []Traveler // List of Travelers that will take on the chiptole world
var stores []Store      // The metadata of the stores read from the csv

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
		fmt.Println("\n\nNEW GAME STARTED\n\n")
		for j := 0; j < population; j++ {
			wg.Add(1)
			go playGame(&travlers[j], &wg)
		}
		wg.Wait()
		breedNextGeneration()
	}
}

func initPopulation() {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	for i := 0; i < population; i++ {
		var dna []Store
		for j := 0; j < r1.Intn(len(stores)-1000)+1000; j++ {
			dna = append(dna, stores[r1.Intn(len(stores))])
		}
		t := newTraveler(dna)
		travlers = append(travlers, t)
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
	t.score = float64(len(t.numStores))
	t.distance = 0
}

func breedNextGeneration() {
	sort.SliceStable(travlers, func(i, j int) bool {
		return travlers[i].score > travlers[j].score
	})

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	travlers = travlers[:len(travlers)-3] // remove the worst 3 performing travlers
	travlers = append(travlers, breed(travlers[r1.Intn(10)+1], travlers[r1.Intn(10)+1]))
	travlers = append(travlers, breed(travlers[r1.Intn(10)+1], travlers[r1.Intn(10)+1]))
	travlers = append(travlers, breed(travlers[r1.Intn(10)+1], travlers[r1.Intn(10)+1]))

	for _, t := range travlers {
		fmt.Println(t.score)
		t.score = 0
	}
}

func breed(p1, p2 Traveler) (child Traveler) {
	var dna []Store

	min := len(p1.dna)
	if min > len(p2.dna) {
		min = len(p2.dna)
	}
	crossover := min / 2

	child = newTraveler(dna)
	return
}
