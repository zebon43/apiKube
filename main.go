package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
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

//Default catch all rule
func catchAllHandler(res http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, " | REST Service not Available.")
	res.WriteHeader(http.StatusBadRequest)
	res.Write([]byte(`Ooopss....This REST Endpoint doesnot exists.`))
}

func apiRequests(kD *kubernetes.Clientset) {
	//routers for the application
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api/v1", apiHomePage)
	router.HandleFunc("/api/v1/services", services)
	router.HandleFunc("/api/v1/services/{applicationGroup}", appGrp)

	router.PathPrefix("/").HandlerFunc(catchAllHandler)

	//Log in case there is an error while the service is running
	log.Fatal(http.ListenAndServe(":8000", router))
}

func getKubeData() *kubernetes.Clientset {
	// path-to-kubeconfig -- for example, /root/.kube/config
	config, err := clientcmd.BuildConfigFromFlags("", "C:/Users/ankit/.kube/config")

	if err != nil {
		log.Println("Unable to find config file.")
	} else {
		// creates the clientset
		clientset, _ := kubernetes.NewForConfig(config)
		return clientset
	}
	return nil
}

func main() {
	log.Println("apiKube Application Started.")

	kD := getKubeData()
	if kD != nil {
		log.Println("Connected to Kubernetes Cluster and data obtained.")
		apiRequests(kD)
	} else {
		log.Println("Unable to connect to Kubernetes Cluster.")
	}

	log.Println("apiKube Application Stopped.")
}
