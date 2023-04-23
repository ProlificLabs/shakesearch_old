package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/hexops/valast"
	_ "modernc.org/sqlite"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", home)

	http.HandleFunc("/search", handleSearch())

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	fmt.Printf("Listening on port %s...", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func writeErrResp(w http.ResponseWriter, message string) {
	errorMessage := map[string]string{"error": message}
	errorJson, _ := json.Marshal((errorMessage))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(errorJson)
}

func home(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("static/index.html"))
	works, err := getWorks()
	chars, err := getChars()

	if err != nil {
		//do something
	}

	pageData := HomePageData{Works: works, Characters: chars}

	tmpl.Execute(w, pageData)
}

func handleSearch() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)
		var qObject SearchQuery

		err := decoder.Decode((&qObject))

		if err != nil {
			writeErrResp(w, "query decoding error")
			return
		}

		if len(qObject.QueryText) == 0 && len(qObject.CharIds) == 0 && len(qObject.WorkIds) == 0 {
			res := []SearchResult{}
			r, _ := json.Marshal((res))
			w.Header().Set("Content-Type", "application/json")
			w.Write(r)
			return
		}

		results, err := executeFTS(qObject)

		if err != nil {
			fmt.Println(valast.String(err))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("db query err"))
			return
		}

		resJson, _ := json.Marshal(results)

		w.Header().Set("Content-Type", "application/json")
		w.Write(resJson)
	}
}
