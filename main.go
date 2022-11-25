package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type StatusChecker interface {
	Check(ctx context.Context, name string) (status bool, err error)
}
type httpChecker struct {
}

func (h httpChecker) Check(ctx context.Context, sitename string) (status bool,
	err error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	req, err := http.NewRequest(http.MethodGet, sitename, nil)
	if err != nil {
		log.Println(err)
	}
	req = req.WithContext(ctx)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	if res.StatusCode != http.StatusOK {
		return false, err
	}
	return true, nil
}

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

func defaultHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	fmt.Println("Default Handler started")
	defer fmt.Println("Default Handler ended")
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

func websitesHandler(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	fmt.Println("websitesHandler Handler started")
	defer fmt.Println("websitesHandler Handler ended")
	select {
	case <-ctx.Done():
		err := ctx.Err()
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	default:
		ctxAsClient := context.Background()
		httpChecker := httpChecker{}

		switch req.Method {
		case http.MethodGet:
			log.Println("IN GET HANDLER")

			var site string = req.URL.Query().Get("name")

			if site != "" {
				_, err := httpChecker.Check(ctxAsClient, site)
				if err != nil {
					websitesMap[site] = "DOWN"
				} else {
					websitesMap[site] = "UP"
				}
				fmt.Fprintf(w, "Site %s is %s", site, websitesMap[site])
			} else {
				for site, status := range websitesMap {
					fmt.Fprintf(w, "%s - %s\n", site, status)
				}
			}

		case http.MethodPost:

			website := Websites{}

			err := json.NewDecoder(req.Body).Decode(&website)
			if err != nil {
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
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}

}

type Website struct {
	Website string `json:"website"`
}

func checksitestatusHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	fmt.Println("checksitestatusHandler Handler started")
	defer fmt.Println("checksitestatusHandler Handler ended")
	select {
	case <-ctx.Done():
		err := ctx.Err()
		log.Print("Ctx done checksitestatusHandler : ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	default:

		ctxAsClient := context.Background()

		website := Website{}

		err := json.NewDecoder(req.Body).Decode(&website)
		if err != nil {
			fmt.Fprint(w, "Error while updating websites.")
		} else {
			httpChecker := httpChecker{}
			_, err := httpChecker.Check(ctxAsClient, website.Website)
			if err != nil {
				websitesMap[website.Website] = "DOWN"
			} else {
				websitesMap[website.Website] = "UP"
			}
			fmt.Fprintf(w, "Site %s is %s", website.Website, websitesMap[website.Website])
		}
	}
}

var websitesMap = make(map[string]string)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go checkSites()

	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/websites", websitesHandler)
	http.HandleFunc("/checksitestatus", checksitestatusHandler)
	http.ListenAndServe("127.0.0.1:8000", nil)

	log.Println(websitesMap)

}
