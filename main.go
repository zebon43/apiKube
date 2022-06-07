package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//Display the home page to the user
func homePage(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, " | Home Page is requested.")
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)

	//routers for the application
	router.HandleFunc("/", homePage)
	//router.HandleFunc("/services", "services")
	//router.HandleFunc("/services/{applicationGroup}", "AppGrp")

	//Log in case there is an error while the service is running
	log.Fatal(http.ListenAndServe(":8000", router))
}

func main() {
	log.Println("apiKube Application Started.")
	handleRequests()
	log.Println("apiKube Application Stopped.")
}
