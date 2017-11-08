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

// /auth. Initial page redirecting
func handleGitHubAuth(w http.ResponseWriter, r *http.Request) {
}

// /callback. Called by github after authorization is granted
func handleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(``))
}

func handleGitHubSuccess(w http.ResponseWriter, r *http.Request) {
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
	http.HandleFunc("/auth", handleGitHubAuth)
	http.HandleFunc("/success", handleGitHubSuccess)
	http.HandleFunc("/callback", handleGitHubCallback)
	fmt.Printf("Started running on %s\n", host)
	fmt.Println(http.ListenAndServe(host, nil))
}
