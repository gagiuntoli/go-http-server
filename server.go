package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"compress/gzip"
)

// Req: http://localhost:1234/upper?word=abc
// Res: ABC
func upperCaseHandler(w http.ResponseWriter, r *http.Request) {
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid request")
		return
	}
	word := query.Get("word")
	if len(word) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "missing word")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, strings.ToUpper(word))
}

// Req: http://localhost:1234/return-json
//
//	Res: {
//	  "description": "I am a JSON",
//	  "number": 12345,
//	}
//
// Bash:
// $ curl localhost:1234/return-json | json_pp
func returnJson(w http.ResponseWriter, r *http.Request) {

	obj := struct {
		Description string `json:"description"`
		Number      uint64 `json:"number"`
	}{
		Description: "I am a JSON",
		Number:      12345,
	}

	_, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid request")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(obj)
}

// Req: http://localhost:1234/return-struct-with-pointers
//
// Bash:
// $ curl localhost:1234/return-struct-with-pointers | json_pp
func returnStrucWithPointers(w http.ResponseWriter, r *http.Request) {

	type internal struct {
            SomeIntPointer    *int    `json:"SomeIntPointer,omitempty"`
            SomeInt           int     `json:"SomeInt,omitempty"`
            SomeStringPointer *string `json:"SomeStringPointer,omitempty"`
            SommeString       string  `json:"SomeString,omitempty"`
        }
        zeroPointer := 1

	_, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid request")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(internal{&zeroPointer, 0, nil, ""})
}

// Req: http://localhost:1234/return-gzip
//
// Bash:
// $ curl -i -H "Accepted-Enconding: gzip" http://localhost:1234/return-gzip | gunzip
// or
// $ curl -i --compressed http://localhost:1234/return-gzip
func returnGzip(w http.ResponseWriter, r *http.Request) {

	obj := []struct {
		Description string `json:"description"`
		Number      uint64 `json:"number"`
	}{
		{
			Description: "I am a JSON",
			Number:      12345,
		},
		{
			Description: "I am another",
			Number:      12346,
		},
	}

	_, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid request")
		return
	}

	encodings := r.Header.Get("Accept-encoding")
	fmt.Println("Accept-encoding", encodings, len(encodings))
	fmt.Println("body", r.Body)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(http.StatusOK)

	gw := gzip.NewWriter(w)
	defer gw.Close()

	json.NewEncoder(gw).Encode(obj)
}

func main() {
	http.HandleFunc("/upper", upperCaseHandler)
	http.HandleFunc("/return-json", returnJson)
	http.HandleFunc("/return-struct-with-pointers", returnStrucWithPointers)
	http.HandleFunc("/return-gzip", returnGzip)
	log.Fatal(http.ListenAndServe(":1234", nil))
}
