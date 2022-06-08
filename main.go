package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var Config *rest.Config

type podCount struct {
	ServiceName string `json:"name,omitempty"`
	AppGrp      string `json:"applicationGroup,omitempty"`
	No          int    `json:"runningPodsCount,omitempty"`
}

//Display the home page to the user
func apiHomePage(res http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, "| API Home Page is requested.")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(`This is the API for Kubernetes.`))
}

//GET request to display number of Pods
func services(res http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, "| Number of Pods is requested.")

	temp, err := getPodData("")
	if err != nil {
		log.Fatalln("Unable to fetch data from Kubernetes. Terminating the Application.", err)
	} else {
		json.NewEncoder(res).Encode(temp)
	}
}

//GET request to display number of Pods per Application Group
func appGrp(res http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, "| Number of Pods per Application Group is requested.")

	//GET Paramter from URL
	params := mux.Vars(req)
	grpName := params["applicationGroup"]

	temp, err := getPodData(grpName)
	if err != nil {
		log.Println(err)
		res.Write([]byte(`Either the application group is not valid or there are no pods in the application group.`))
	} else {
		json.NewEncoder(res).Encode(temp)
	}
}

//Default catch all rule
func catchAllHandler(res http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, "| REST Service not Available.")
	res.WriteHeader(http.StatusBadRequest)
	res.Write([]byte(`Ooopss....This REST Endpoint doesnot exists.`))
}

func apiRequests() {
	//routers for the application
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api/v1", apiHomePage)
	router.HandleFunc("/api/v1/services", services)
	router.HandleFunc("/api/v1/services/{applicationGroup}", appGrp)
	router.PathPrefix("/").HandlerFunc(catchAllHandler)

	//Log in case there is an error while the service is running
	log.Fatal("User has terminated the Application.", http.ListenAndServe(":8000", router))
}

func getPodData(grpName string) ([]podCount, error) {

	var listFilter v1.ListOptions
	//Set List Options
	if grpName != "" {
		listFilter = v1.ListOptions{LabelSelector: "applicationGroup=" + grpName}
	} else {
		listFilter = v1.ListOptions{}
	}

	// creates the clientset
	kD, err := kubernetes.NewForConfig(Config)
	if err == nil {
		log.Println("Connected to Kubernetes Cluster and data obtained.")
	} else {
		log.Fatalln("Unable to connect to Kubernetes Cluster. Terminating the Application.", err)
	}

	pods, err := kD.CoreV1().Pods("default").List(context.TODO(), listFilter)

	var pc []podCount
	addFlag := true
	//Check for invalid data in URL
	if err != nil {
		log.Fatalln("Unable to connect to Kubernetes. Terminating the Application.", err)
	} else if len(pods.Items) == 0 {
		return nil, errors.New("either the Application Group is not valid or there are no pods in the application group")
	} else {
		//Logic to get the data for the REST Endpoint
		for _, pod := range pods.Items {
			for j, temp := range pc {
				if pod.ObjectMeta.Labels["service"] == temp.ServiceName {
					pc[j].No++
					addFlag = false
					break
				}
				addFlag = true
			}
			if addFlag {
				pc = append(pc, podCount{pod.ObjectMeta.Labels["service"], pod.ObjectMeta.Labels["applicationGroup"], 1})
			}
		}
	}
	return pc, nil
}

func main() {
	log.Println("apiKube Application Started.")

	// path-to-kubeconfig -- for example, /root/.kube/config
	var err error
	Config, err = clientcmd.BuildConfigFromFlags("", "C:/Users/ankit/.kube/config")
	if err != nil {
		log.Println(err)
		log.Fatalln("Unable to find config file. Terminating the Application.")
	} else {
		apiRequests()
	}

	log.Println("apiKube Application Stopped.")
}
