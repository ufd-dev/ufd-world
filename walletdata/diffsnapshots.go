package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"maps"
	"math"
	"os"
	"strconv"
)

const (
	colAccount = iota
	colTokenAccount
	colQuantity
	colPercentage
)

type balanceChange struct {
	account string
	start   float64
	end     float64
	diff    float64
}

func main() {
	f1, err := os.Open("top1000snapshots/20250321.csv")
	if err != nil {
		log.Fatal("Error while reading first file", err)
	}
	defer f1.Close()

	reader1 := csv.NewReader(f1)
	// slice of slices of string
	records1, err := reader1.ReadAll()
	if err != nil {
		fmt.Println("Error reading first file's records")
	}

	f2, err := os.Open("top1000snapshots/20250503.csv")
	if err != nil {
		log.Fatal("Error while reading second file", err)
	}
	defer f1.Close()

	reader2 := csv.NewReader(f2)
	// slice of slices of string
	records2, err := reader2.ReadAll()
	if err != nil {
		fmt.Println("Error reading second file's records")
	}

	records1 = records1[1:]
	changes := make(map[string]*balanceChange, int(math.Round(1.5*float64(len(records1)))))
	for _, r := range records1 {
		start, err := strconv.ParseFloat(r[colQuantity], 64)
		if err != nil {
			panic("couldn't parse a balance")
		}
		acct := r[colAccount]
		changes[acct] = &balanceChange{account: acct, start: start, diff: -start}
	}

	records2 = records2[1:]
	for _, r := range records2 {
		acct := r[colAccount]
		end, err := strconv.ParseFloat(r[colQuantity], 64)
		if err != nil {
			panic("couldn't parse a balance")
		}
		change, ok := changes[acct]
		if ok {
			change.end = end
			change.diff = end - change.start
		} else {
			changes[acct] = &balanceChange{account: acct, start: 0, end: end, diff: end}
		}
	}

	for c := range maps.Values(changes) {
		fmt.Printf("%s\t%f\t%f\t%f\n", c.account, c.start, c.end, c.diff)
	}
}
