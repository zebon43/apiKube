package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//Display the home page to the user
func apiHomePage(res http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, " | API Home Page is requested.")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(`This is the API for Kubernetes.`))
}

//GET request to display number of Pods
func services(res http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, " | Number of Pods is requested.")
}

//GET request to display number of Pods per Application Group
func appGrp(res http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, " | Number of Pods per Application Group is requested.")

	//GET Paramter from URL
	params := mux.Vars(req)
	grpName := params["applicationGroup"]

	log.Println(grpName)
}

func apiRequests() {
	//routers for the application
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api/v1", apiHomePage)
	router.HandleFunc("/api/v1/services", services)
	router.HandleFunc("/api/v1/services/{applicationGroup}", appGrp)

	//Log in case there is an error while the service is running
	log.Fatal(http.ListenAndServe(":8000", router))
}

func main() {
	log.Println("apiKube Application Started.")
	apiRequests()
	log.Println("apiKube Application Stopped.")
}
