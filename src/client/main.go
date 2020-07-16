package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

var oauth = struct {
	authURL            string
	logout             string
}{
	authURL: "http://10.100.196.60:8080/auth/realms/learningApp/protocol/openid-connect/auth",
	logout: "http://10.100.196.60:8080/auth/realms/learningApp/protocol/openid-connect/logout",
}

var t = template.Must(template.ParseFiles("template/index.html"))
type AppVar struct {
	AuthCode string
	SessionState string
}
var appVar = AppVar{}

func main()  {
	fmt.Println("hello")
	http.HandleFunc("/", home)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/authCodeRedirect", authCodeRedirect)
	http.ListenAndServe(":8080", nil)
}

func logout(writer http.ResponseWriter, request *http.Request) {
	log.Printf("logout: Request queries: %v", request.URL.Query())

	qs:= url.Values{}
	qs.Add("redirect_uri", "http://localhost:8080")
	logoutURL, err := url.Parse(oauth.logout)
	if err != nil {
		log.Print(err)
		return
	}
	logoutURL.RawQuery = qs.Encode()
	http.Redirect(writer, request, logoutURL.String(), http.StatusFound)
	log.Printf("logout: done")
}

// AS will send us back here after login so we can collect the creds
//    and then redirect home
func authCodeRedirect(writer http.ResponseWriter, request *http.Request) {
	log.Printf("authCodeRedirect: Request queries: %v", request.URL.Query())
	appVar.AuthCode = request.URL.Query().Get("code")
	appVar.SessionState = request.URL.Query().Get("session_state")
	request.URL.RawQuery = ""
	log.Printf("Request queries: %+v", appVar)
	//t.Execute(writer, nil)
	http.Redirect(writer, request, "http://localhost:8080", http.StatusFound)
	log.Printf("authCodeRedirect: done")
}

func home(writer http.ResponseWriter, request *http.Request) {
	log.Printf("home: Request queries: %v", request.URL.Query())
	t.Execute(writer, nil)
	log.Printf("home: done")
}

func login(writer http.ResponseWriter, request *http.Request) {
	log.Printf("login: Request queries: %v", request.URL.Query())
	req, err := http.NewRequest("GET", oauth.authURL,nil)
	if err != nil {
		log.Print(err)
		return
	}
	qs:= url.Values{}
	qs.Add("state", "123")
	qs.Add("client_id", "billingApp")
	qs.Add("response_type", "code")
	qs.Add("redirect_uri", "http://localhost:8080/authCodeRedirect")
	req.URL.RawQuery = qs.Encode()
	http.Redirect(writer, request, req.URL.String(), http.StatusFound)
	log.Printf("login: done")
}

