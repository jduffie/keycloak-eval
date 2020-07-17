package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime"
)

// Billing response
type Billing struct {
	Services []string `json:"services"`
}

type AppVar struct {
}

var appVar = AppVar{}

func main() {
	fmt.Println("CLIENT")
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/billing/v1/list", enabledLog(services))
	http.ListenAndServe(":8081", nil)
}

func services(writer http.ResponseWriter, request *http.Request) {
	s := Billing{
		Services: []string {
			"electric",
			"phone",
			"internet",
			"water",
		},
	}

	encoder := json.NewEncoder(writer)
	writer.Header().Add("Content-Type", "application/json")
	encoder.Encode(s)

	// fmt.Fprintf(writer, "%v", res)
}

func enabledLog(h func(writer http.ResponseWriter, request *http.Request)) func(http.ResponseWriter, *http.Request) {

	return func(writer http.ResponseWriter, request *http.Request) {
		handlerName := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
		log.Printf("--> %s ", handlerName)
		log.SetPrefix(handlerName + " ")
		log.Printf("request: %+v\n", request)
		log.Printf("response: %+v\n", writer)
		h(writer, request)
		log.Printf("<-- %s ", handlerName)
	}

}
