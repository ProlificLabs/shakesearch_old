package main

import (
	"fmt"
	"net/http"
	"pulley.com/shakesearch/search"
	"pulley.com/shakesearch/utils"
)

const (
	fileName = "completeworks.txt"
)

var scanner *search.Scanner

func main() {
	scanner = search.NewScanner()
	utils.ErrorFatalCheck(scanner.Load(fileName))

	// serving static assets ...
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// search functionality handler ...
	http.HandleFunc("/search", handleSearch)

	port := utils.GetAppPort()
	fmt.Printf("Listening on port %s...", port)
	utils.ErrorFatalCheck(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	query, ok := r.URL.Query()["q"]
	if !ok || len(query[0]) < 1 {
		utils.WriteResponse(http.StatusBadRequest, []byte("missing search query in URL params"), w)
		return
	}
	results, correction := scanner.Search(query[0])
	buf, err := utils.EncodeResult(results, correction)
	if err != nil {
		utils.WriteResponse(http.StatusInternalServerError, []byte("encoding failure"), w)
		return
	}
	utils.WriteResponse(http.StatusOK, buf.Bytes(), w)
}
