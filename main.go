package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type StatusChecker interface {
	Check(ctx context.Context, name string) (status bool, err error)
}

type httpChecker struct {
}

/*
Function having context working as client to fetch status of websit
*/

func (h httpChecker) Check(ctx context.Context, sitename string) (status bool,
	err error) {

	req, err := http.NewRequest(http.MethodGet, sitename, nil)
	if err != nil {
		log.Println(err)
	}

	client := http.Client{
		Timeout: 1 * time.Minute,
	}
	req = req.WithContext(ctx)
	res, err := client.Do(req)
	if err != nil {
		return false, err
	}
	if res.StatusCode != http.StatusOK {
		return false, err
	}
	return true, nil
}

/*
go routines working as a server context with timeout of 1 minute
*/
func checkSites() {
	log.Println("STATUS CHECK STARTED")
	ctxAsClient := context.Background()
	httpChecker := httpChecker{}

	for {
		select {
		case <-time.After(1 * time.Minute):
			log.Println("STATUC CHECK :", websitesMap)

			for site := range websitesMap {
				status, err := httpChecker.Check(ctxAsClient, site)
				if err != nil || !status {
					websitesMap[site] = "DOWN"
				} else {
					websitesMap[site] = "UP"
				}
			}
		case <-ctxAsClient.Done():
			err := ctxAsClient.Err()
			log.Println(ctxAsClient, err.Error())
		}
	}
}

/*
api handler working as a server context to handle default api requests
*/
func defaultHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log.Println("Default Handler started")
	defer log.Println("Default Handler ended")

	select {
	case <-ctx.Done():
		err := ctx.Err()
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	default:
		fmt.Fprint(w, "Redirected to default page")
	}
}

type Websites struct {
	Websites []string `json:"websites"`
}

/*
api handler working as a server context to handle below POST api requests
URL : /websites
	Method : POST
	Request Body :
				{"websites":["http://www.google.com","http://www.facebook.com","http://www.fakewebsite1.com"]}
	Response :
				Websites updated successfully.
	Updates memory map of websites with website statuses
*/
func loadWebsitesHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log.Println("websitesHandler Handler started")
	defer log.Println("websitesHandler Handler ended")

	select {
	case <-ctx.Done():
		err := ctx.Err()
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	default:
		ctxAsClient := context.Background()
		httpChecker := httpChecker{}
		website := Websites{}

		err := json.NewDecoder(req.Body).Decode(&website)
		if err != nil {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Error while updating websites.")
		} else {
			for _, site := range website.Websites {
				_, err := httpChecker.Check(ctxAsClient, site)
				if err != nil {
					websitesMap[site] = "DOWN"
				} else {
					websitesMap[site] = "UP"
				}
			}
			fmt.Fprint(w, "Websites updated successfully.")
			log.Println(websitesMap)
		}
	}

}

/*
api handler working as a server context to handle below GET api requests
1. URL : /websites?name=""
	Method : GET
	Response :

				http://www.google.com - UP
				http://www.facebook.com - UP
				http://www.fakewebsite1.com - DOWN
*/
func checkSiteStatusHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log.Println("websitesHandler Handler started")
	defer log.Println("websitesHandler Handler ended")

	select {
	case <-ctx.Done():
		err := ctx.Err()
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	default:
		ctxAsClient := context.Background()
		httpChecker := httpChecker{}

		log.Println("IN GET HANDLER")

		var site string = req.URL.Query().Get("name")

		if site != "" {
			_, err := httpChecker.Check(ctxAsClient, site)
			var res = make(map[string]string)
			if err != nil {
				res[site] = "DOWN"
				websitesMap[site] = "DOWN"
			} else {
				res[site] = "UP"
				websitesMap[site] = "UP"
			}
			resp, _ := json.Marshal(res)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(resp)
		} else {
			resp, _ := json.Marshal(websitesMap)
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(resp)
		}
	}

}

type Website struct {
	Website string `json:"website"`
}

var websitesMap = make(map[string]string)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go checkSites()

	router := mux.NewRouter()

	router.HandleFunc("/", defaultHandler).Methods(http.MethodGet)
	router.HandleFunc("/websites", checkSiteStatusHandler).Methods(http.MethodGet)
	router.HandleFunc("/websites", loadWebsitesHandler).Methods(http.MethodPost)
	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
	}

	log.Fatal(srv.ListenAndServe())
	wg.Wait()
}
