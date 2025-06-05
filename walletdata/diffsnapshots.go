package main

import (
	"cmp"
	"encoding/csv"
	"fmt"
	"log"
	"maps"
	"math"
	"os"
	"slices"
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
	d1, d2 := "20250514", "20250605"

	f1, err := os.Open(fmt.Sprintf("top1000snapshots/%s.csv", d1))
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

	f2, err := os.Open(fmt.Sprintf("top1000snapshots/%s.csv", d2))
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
	changeMap := make(map[string]*balanceChange, int(math.Round(1.5*float64(len(records1)))))
	for _, r := range records1 {
		start, err := strconv.ParseFloat(r[colQuantity], 64)
		if err != nil {
			panic("couldn't parse a balance")
		}
		acct := r[colAccount]
		changeMap[acct] = &balanceChange{account: acct, start: start, diff: -start}
	}

	records2 = records2[1:]
	for _, r := range records2 {
		acct := r[colAccount]
		end, err := strconv.ParseFloat(r[colQuantity], 64)
		if err != nil {
			panic("couldn't parse a balance")
		}
		change, ok := changeMap[acct]
		if ok {
			change.end = end
			change.diff = end - change.start
		} else {
			changeMap[acct] = &balanceChange{account: acct, start: 0, end: end, diff: end}
		}
	}

	sortedChanges := slices.SortedFunc(maps.Values(changeMap), func(a, b *balanceChange) int {
		return cmp.Compare(a.account, b.account)
	})
	for _, c := range sortedChanges {
		fmt.Printf("%s\t%f\t%f\t%f\n", c.account, c.start, c.end, c.diff)
	}
}
