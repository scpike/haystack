package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"os"
	"math/rand"
	"sort"
	"strings"
)

func contains(arr []int, i int) bool {
   for _, a := range arr {
      if a == i {
         return true
      }
   }
   return false
}

func generateNeedles(needle []string, totalLength int) (map[int]string, error){
	var indices []int
	if (len(needle) > totalLength) {
		return nil, fmt.Errorf("Needle length %d can't be larger than haystack length %d", len(needle), totalLength)
	}
	for range needle {
		c := rand.Intn(totalLength)
		for (contains(indices, c)) {
			c = rand.Intn(totalLength)
		}
		indices = append(indices, c)
	}
	sort.Ints(indices)
	var retmap map[int]string
	retmap = make(map[int]string)
	for idx, i := range indices {
		retmap[i + 1] = needle[idx]
	}
	fmt.Println(indices)
	fmt.Println(retmap)
	return retmap, nil
}

func makeHandler(needle string, haystackSize int) (func(http.ResponseWriter, *http.Request), error) {
	nmap, err := generateNeedles(strings.Split(needle, ""), haystackSize)

	if (err == nil) {
		return func(w http.ResponseWriter, r *http.Request){
			path := r.URL.Path[1:]
			var secret string
			if (path != "") {
				i, err := strconv.Atoi(path)
				if (err != nil) {
					w.WriteHeader(500)
					fmt.Fprint(w, "Give number paths like /1")
					return
				}
				if (nmap[i] != "") {
					secret = nmap[i]
				} else {
					secret = "Nothing here, keep looking!"
					if (err == nil) && (i > haystackSize) {
						w.WriteHeader(404)
						fmt.Fprint(w, "You've reached the end")
						return
					}
				}
				fmt.Fprintf(w, "{\"secret\" => \"%s\"}", secret)
			} else {
				fmt.Fprintf(w, "{\"message\" => \"Go to /1 to get started!\"}")
			}
		}, nil
	} else {
		return nil, err
	}
}

func makeTotalHandler(haystackSize int) func(http.ResponseWriter, *http.Request) {
	return func (w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w, "{\"total_length\" => %d}", haystackSize)
	}
}

func main() {
	totalLength := 1000 // default
	needle := "SEVENFIFTY"
	if len(os.Args) >= 2 {
		length, err := strconv.Atoi(os.Args[1])
		if (err == nil) {
			totalLength = length
		}
	}
	if len(os.Args) >= 3 {
		needle = os.Args[2]
	}

	fmt.Printf("Setting up a server with %d entries, secret is %s\n", totalLength, needle)
	handler, err := makeHandler(needle, totalLength)
	if (err == nil) {
		http.HandleFunc("/", handler)
		http.HandleFunc("/total", makeTotalHandler(totalLength))
		log.Fatal(http.ListenAndServe(":8080", nil))
	} else {
		fmt.Println(err)
	}
}
