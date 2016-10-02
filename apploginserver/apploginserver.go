// Copyright 2016 R. Stanley Hum
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
	App Login Server Microservice
*/
package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"

	"github.com/google/go-github/github"
)

var (
	githuboauthConf = &oauth2.Config{
		ClientID:     "7b018a98957cc0cc010f",
		ClientSecret: "fac68191a5cbb1bc999eb603cab2f4b228d51941",
		Endpoint:     githuboauth.Endpoint,
	}
	oauthStateString = "lajd39od2rnf48rjgk1341vsd"
)

const htmlIndex = `<html><body>Logged in <a href="/login">GitHub</a></body></html>`

type Profile struct {
	Name       string
	Activities []string
}

type Message struct {
	StatusCode int
	Status     string
	Url        string
}

type ReturnToken struct {
	Token      oauth2.Token
	StatusCode int
	Status     string
	User       string
}

// /
func handleMain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlIndex))
}

// /login
func handleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	provider := githuboauthConf
	url := provider.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// /github_oauth_cb
func handleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expended '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := githuboauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	oauthClient := githuboauthConf.Client(oauth2.NoContext, token)
	client := github.NewClient(oauthClient)
	user, _, err := client.Users.Get("")
	if err != nil {
		fmt.Printf("client.users.Get() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	w.Header().Set("Content-Type", "application.json; charset=utf-8")
	returnToken := ReturnToken{*token, http.StatusOK, "Logged In", *user.Login}
	js, err := json.Marshal(returnToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func main() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleGitHubLogin)
	http.HandleFunc("/github_oauth_cb", handleGitHubCallback)

	fmt.Printf("Server running on localhost:8080\n")
	fmt.Println(http.ListenAndServeTLS(":8080", "certs/server.pem", "certs/server.key", nil))
}
