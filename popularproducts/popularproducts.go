/*
Given a text file where each line has a JSON format like {"user_id": "KB/WRFTC", "product_id": "F5H", "quantity": 12},
  output the most popular products by two different metrics: the unique users purchasing the product and the quantity sold.
*/
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	fmt.Println("Hello World!")

	lines, err := parse("./pp_data.txt")
	if err != nil {
		panic(err) //meh
	}

	findPopularity(lines)
}

// The best way to trace this is with two hashmaps.
// The quantity hashmap is straightforward: increment quantity.
// The unique hashmap is slightly more complicated but can be done with nested hashmaps.
// The maxes can be tracked by maintaining a list of everything at that max then throwing it out when the max changes.
func findPopularity(lines []line) {
	qmap := make(map[string]int)             // products onto total quantity
	umap := make(map[string]map[string]bool) // product onto map of users onto nothing
	var qmax, umax int
	var qproduct, uproduct []string

	for _, l := range lines {
		// if it exists in qmap, increment. Else add.
		var newQ int
		if q, ok := qmap[l.ProductID]; ok {
			newQ = q + l.Quantity
			qmap[l.ProductID] = newQ
		} else {
			newQ = l.Quantity
			qmap[l.ProductID] = newQ
		}
		// replace or add to max list
		if newQ > qmax {
			// purge the prior list
			qmax = newQ
			qproduct = []string{l.ProductID}
		} else if newQ == qmax {
			// something cannot be added twice because changing the quantity would have purged the list
			qproduct = append(qproduct, l.ProductID)
		}

		// Same deal with unique users. If it exists then increment, else add.
		var newU int
		if m, ok := umap[l.ProductID]; ok {
			// if the item exists, makes sure the user does not
			if _, uexists := m[l.UserID]; !uexists {
				// m != umap[l.ProductID], so modify that in-place
				umap[l.ProductID][l.UserID] = true
				newU = len(umap[l.ProductID])
			} // do not set newU if the user already exists, otherwise the trace block might fire.
		} else {
			umap[l.ProductID] = make(map[string]bool)
			umap[l.ProductID][l.UserID] = true
			newU = 1
		}
		if newU > umax {
			umax = newU
			uproduct = []string{l.ProductID}
		} else if newU == umax {
			uproduct = append(uproduct, l.ProductID)
		}
	}
	fmt.Printf("Most popular product(s) based on the quantity of goods sold: %v at %d units\n", qproduct, qmax)
	fmt.Printf("Most popular product(s) based on the number of purchasers: %v at %d unique users\n", uproduct, umax)
}

type line struct {
	UserID    string `json:"user_id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

func parse(filename string) ([]line, error) {
	// The input file is "fairly small" so it can be read in rather than buffered
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	items := strings.Split(string(data), "\n")
	lines := make([]line, len(items))
	for i, s := range items {
		var l line
		err = json.Unmarshal([]byte(s), &l)
		if err != nil {
			return nil, err
		}
		lines[i] = l
	}
	return lines, nil
}
