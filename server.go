package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
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
	//fmt.Fprintf(w, strings.ToUpper(word))
}

func main() {
	http.HandleFunc("/upper", upperCaseHandler)
	http.HandleFunc("/return-json", returnJson)
	log.Fatal(http.ListenAndServe(":1234", nil))
}
