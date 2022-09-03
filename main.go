package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/bitbucket"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gitlab"
	log "github.com/sirupsen/logrus"
)

var (
	host = ":3000"
)

const (
	script = `<!DOCTYPE html><html><head><script>
  (function(status, provider, result) {
    function recieveMessage(e) {
      console.log("Recieve message:", e);
      // send message to main window with da app
      window.opener.postMessage(
        "authorization:" + provider + ":" + status + ":" + JSON.stringify(result),
        e.origin
      );
    }
    window.addEventListener("message", recieveMessage, false);
    // Start handshare with parent
    console.log("Sending message:", provider)
    window.opener.postMessage(
      "authorizing:" + provider,
      "*"
    );
  })("%s", "%s", %s)
  </script></head><body></body></html>`
)

// GET /
func handleMain(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	_, err := res.Write([]byte(``))
	if err != nil {
		log.Error(fmt.Printf("HandleMain error: %s", err))
	}
}

// GET /auth Page  redirecting after provider get param
func handleAuth(res http.ResponseWriter, req *http.Request) {
	url := fmt.Sprintf("%s/auth/%s", host, req.FormValue("provider"))
	log.Info(fmt.Printf("redirect to %s\n", url))
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}

// GET /auth/provider  Initial page redirecting by provider
func handleAuthProvider(res http.ResponseWriter, req *http.Request) {
	gothic.BeginAuthHandler(res, req)
}

// GET /callback/{provider}  Called by provider after authorization is granted
func handleCallbackProvider(res http.ResponseWriter, req *http.Request) {
	var (
		status string
		result string
	)
	provider, errProvider := gothic.GetProviderName(req)
	user, errAuth := gothic.CompleteUserAuth(res, req)
	status = "error"
	if errProvider != nil {
		fmt.Printf("provider failed with '%s'\n", errProvider)
		result = fmt.Sprintf("%s", errProvider)
	} else if errAuth != nil {
		fmt.Printf("auth failed with '%s'\n", errAuth)
		result = fmt.Sprintf("%s", errAuth)
	} else {
		fmt.Printf("Logged in as %s user: %s (%s)\n", user.Provider, user.Email, user.AccessToken)
		status = "success"
		result = fmt.Sprintf(`{"token":"%s", "provider":"%s"}`, user.AccessToken, user.Provider)
	}
	//res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	_, err := res.Write([]byte(fmt.Sprintf(script, status, provider, result)))
	if err != nil {
		log.Error(fmt.Printf("handleCallbackProvider error: %s", err))
	}
}

// GET /refresh
func handleRefresh(res http.ResponseWriter, req *http.Request) {
	log.Info(fmt.Printf("refresh"))
	res.WriteHeader(http.StatusNoContent)
}

// GET /success
func handleSuccess(res http.ResponseWriter, req *http.Request) {
	log.Info(fmt.Printf("success\n"))
	res.WriteHeader(http.StatusNoContent)
}

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if hostEnv, ok := os.LookupEnv("HOST"); ok {
		host = hostEnv
	}
	var (
		gitlabProvider goth.Provider
	)
	if gitlabServer, ok := os.LookupEnv("GITLAB_SERVER"); ok {
		gitlabProvider = gitlab.NewCustomisedURL(
			os.Getenv("GITLAB_KEY"), os.Getenv("GITLAB_SECRET"),
			fmt.Sprintf("%s/callback/gitlab", host),
			fmt.Sprintf("%s/oauth/authorize", gitlabServer),
			fmt.Sprintf("%s/oauth/token", gitlabServer),
			fmt.Sprintf("%s/api/v3/user", gitlabServer),
		)
	} else {
		gitlabProvider = gitlab.New(
			os.Getenv("GITLAB_KEY"), os.Getenv("GITLAB_SECRET"),
			fmt.Sprintf("%s/callback/gitlab", host),
		)
	}
	goth.UseProviders(
		github.New(
			os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"),
			fmt.Sprintf("%s/callback/github", host),
		),
		bitbucket.New(
			os.Getenv("BITBUCKET_KEY"), os.Getenv("BITBUCKET_SECRET"),
			fmt.Sprintf("%s/callback//bitbucket", host),
		),
		gitlabProvider,
	)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/callback/{provider}", handleCallbackProvider).Methods(http.MethodGet)
	router.HandleFunc("/auth/{provider}", handleAuthProvider).Methods(http.MethodGet)
	router.HandleFunc("/auth", handleAuth).Methods(http.MethodGet)
	router.HandleFunc("/refresh", handleRefresh).Methods(http.MethodGet)
	router.HandleFunc("/success", handleSuccess).Methods(http.MethodGet)
	router.HandleFunc("/", handleMain).Methods(http.MethodGet)
	//
	err := http.ListenAndServe(":3000", router)
	if err != nil {
		log.Fatal(fmt.Printf("Can't start webserver with error: %s", err.Error()))
	}
}
