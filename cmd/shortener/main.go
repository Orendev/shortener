package main

import (
	"github.com/Orendev/shortener/internal/app/random"
	"io"
	"net/http"
	"strings"
)

var mux = http.NewServeMux()
var db = make(map[string]string)

func main() {

	mux.HandleFunc("/", mainPage)

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	return http.ListenAndServe(`:8080`, mux)
}

func mainPage(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		w.Header().Set("Content-Type", "text/plain")

		body, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		code := random.Strn(8)
		db[code] = string(body)
		url := "http://localhost:8080/" + code

		w.WriteHeader(http.StatusCreated)

		_, _ = w.Write([]byte(url))

	} else {

		code := strings.TrimPrefix(r.URL.Path, "/")

		if url, ok := db[code]; ok {
			w.Header().Add("Location", url)
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}

		w.WriteHeader(http.StatusBadRequest)

	}
}
