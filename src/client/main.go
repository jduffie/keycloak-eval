package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

var config = struct {
	authURL                  string
	afterAuthCodeRedirectURL string
	logout                   string
	afterLogoutRedirect      string
	clientId                 string
}{
	authURL:                  "http://10.100.196.60:8080/auth/realms/learningApp/protocol/openid-connect/auth",
	afterAuthCodeRedirectURL: "http://localhost:8080/afterAuthCodeRedirectURL",
	logout:                   "http://10.100.196.60:8080/auth/realms/learningApp/protocol/openid-connect/logout",
	afterLogoutRedirect:      "http://localhost:8080",
	clientId:                 "billingApp",
}

var t = template.Must(template.ParseFiles("template/index.html"))

type AppVar struct {
	AuthCode         string
	SessionState     string
}
var appVar = AppVar{}

func main()  {
	fmt.Println("hello")
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", home)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/afterAuthCodeRedirectURL", authCodeRedirect)
	http.ListenAndServe(":8080", nil)
}



func home(writer http.ResponseWriter, request *http.Request) {
	log.Printf("home: Request queries: %v", request.URL.Query())
	t.Execute(writer, appVar)
	log.Printf("home: done")
}

//(A)  The client initiates the flow by directing the resource owner's
//user-agent to the authorization endpoint.  The client includes
//its client identifier, requested scope, local state, and a
//redirection URI to which the authorization server will send the
//user-agent back once access is granted (or denied).
//(B)  The authorization server authenticates the resource owner (via
//the user-agent) and establishes whether the resource owner
//grants or denies the client's access request.
func login(writer http.ResponseWriter, request *http.Request) {
	log.Printf("login: Request queries: %v", request.URL.Query())
	req, err := http.NewRequest("GET", config.authURL,nil)
	if err != nil {
		log.Print(err)
		return
	}
	qs:= url.Values{}
	qs.Add("state", "123")
	qs.Add("client_id", config.clientId)
	qs.Add("response_type", "code")
	qs.Add("redirect_uri", config.afterAuthCodeRedirectURL)
	req.URL.RawQuery = qs.Encode()
	http.Redirect(writer, request, req.URL.String(), http.StatusFound)
	log.Printf("login: done")
}

//(C)  Assuming the resource owner grants access, the authorization
//server redirects the user-agent back to the client using the
//redirection URI provided earlier (in the request or during
//client registration).  The redirection URI includes an
//authorization code and any local state provided by the client
//earlier.
//(D)  The client requests an access token from the authorization
//server's token endpoint by including the authorization code
//received in the previous step.  When making the request, the
//client authenticates with the authorization server.  The client
//includes the redirection URI used to obtain the authorization
//code for verification.
func authCodeRedirect(writer http.ResponseWriter, request *http.Request) {
	log.Printf("authCodeRedirect: Request queries: %v", request.URL.Query())
	appVar.AuthCode = request.URL.Query().Get("code")
	appVar.SessionState = request.URL.Query().Get("session_state")
	request.URL.RawQuery = ""
	log.Printf("Request queries: %+v", appVar)
	http.Redirect(writer, request, "http://localhost:8080", http.StatusFound)
	log.Printf("authCodeRedirect: done")
}

func logout(writer http.ResponseWriter, request *http.Request) {
	log.Printf("logout: Request queries: %v", request.URL.Query())

	qs:= url.Values{}
	qs.Add("redirect_uri", config.afterLogoutRedirect)
	logoutURL, err := url.Parse(config.logout)
	if err != nil {
		log.Print(err)
		return
	}
	logoutURL.RawQuery = qs.Encode()
	appVar = AppVar{}
	http.Redirect(writer, request, logoutURL.String(), http.StatusFound)
	log.Printf("logout: done")
}
