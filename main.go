package main

import (
	"fmt"
	"net/http"
	"os"

	githubclient "github.com/google/go-github/github"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"

	"./dotenv"
	rand "./randstr"
)

var (
	host = "localhost:3000"
	// random string for oauth2 API calls to protect against CSRF
	oauthStateString = rand.RandomString(64)
	oauthConf        = &oauth2.Config{
		Endpoint:     githuboauth.Endpoint,
		Scopes:       []string{"repo", "user"},
		ClientID:     os.Getenv("GITHUB_KEY"),
		ClientSecret: os.Getenv("GITHUB_SECRET"),
	}
)

const (
	script = `<!DOCTYPE html><html><head><script>
  if (!window.opener) {
    window.opener = {
      postMessage: function(action, origin) {
        console.log(action, origin);
      }
    }
  }
  (function(status, provider, result) {
    function recieveMessage(e) {
      console.log("Recieve message:", e);
      // send message to main window with da app
      window.opener.postMessage(
        "authorization:" + provider + ":" + status + ":" + result,
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
	res.Write([]byte(``))
}

// GET /auth  Initial page redirecting
func handleGitHubAuth(res http.ResponseWriter, req *http.Request) {
	url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	http.Redirect(res, req, url, http.StatusTemporaryRedirect)
}

// GET /callback  Called by github after authorization is granted
func handleGitHubCallback(res http.ResponseWriter, req *http.Request) {
	var (
		status string
		result string
	)
	state := req.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
	}
	provider := "github"
	code := req.FormValue("code")
	token, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		status = "error"
		result = fmt.Sprintf("%s", err)
	} else {
		oauthClient := oauthConf.Client(oauth2.NoContext, token)
		client := githubclient.NewClient(oauthClient)
		user, _, err := client.Users.Get(oauth2.NoContext, "")
		if err != nil {
			fmt.Printf("client.Users.Get() falled with '%s'\n", err)
			status = "error"
			result = fmt.Sprintf("%s", err)
		} else {
			fmt.Printf("Logged in as github user: %s (%s)\n", *user.Login, token.AccessToken)
			status = "success"
			result = fmt.Sprintf(`{"token":"%s", "provider":"%s"}`, token.AccessToken, provider)
		}
	}
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(fmt.Sprintf(script, status, provider, result)))
}

func handleRefresh(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("refresh with '%s'\n", req)
	res.Write([]byte(""))
}

// GET /success
func handleSuccess(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("success with '%s'\n", req)
	res.Write([]byte(""))
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
	//
	fmt.Printf("Started running on %s\n", host)
	fmt.Println(http.ListenAndServe(host, nil))
}
