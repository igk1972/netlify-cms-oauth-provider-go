package main

import (
	"fmt"
	"net/http"
	"os"

	"./dotenv"
)

var (
	host = "localhost:3000"
)

// /
func handleMain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(``))
}

func init() {
	dotenv.File(".env")
	if hostEnv, ok := os.LookupEnv("HOST"); ok {
		host = hostEnv
	}
}

func main() {
	http.HandleFunc("/", handleMain)
	fmt.Printf("Started running on %s\n", host)
	fmt.Println(http.ListenAndServe(host, nil))
}
