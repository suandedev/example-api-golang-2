package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "github.com/gorilla/mux"
)

type event struct {
    ID          string `json:"ID"`
    Title       string `json:"Title"`
    Description string `json:"Description"`
}

type allEvent []event

var events = allEvent {
    {
        ID :            "1",
        Title :         "Event one",
        Description :   "this is event one",
    },
}

func homeLink(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "welcome message")
}

func getAllEvents(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(events)
}

func createEvent(w http.ResponseWriter, r *http.Request) {
    var newEvent event
    reqBody, err := ioutil.ReadAll(r.Body)
    if err != nil {
        fmt.Fprintf(w, "something wrong")
    }

    json.Unmarshal(reqBody, &newEvent)
    events = append(events, newEvent)
    w.WriteHeader(http.StatusCreated)

    json.NewEncoder(w).Encode(newEvent)
}

func getEventById(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]

    for _, singgleEvent := range events {
        if singgleEvent.ID == id {
            json.NewEncoder(w).Encode(singgleEvent)
        }
    }
}

func deleteEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			events = append(events[:i], events[i+1:]...)
			fmt.Fprintf(w, "The event with ID %v has been deleted successfully", eventID)
		}
	}
}

func updateEvent(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    var updateEvent event

    reqBody, err := ioutil.ReadAll(r.Body)
    if err != nil {
        fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
    }
    json.Unmarshal(reqBody, &updateEvent)

    for i, singgleEvent := range events {
        if singgleEvent.ID == id {
            singgleEvent.Title = updateEvent.Title
            singgleEvent.Description = updateEvent.Description
            events = append(events[:i], singgleEvent)
            json.NewEncoder(w).Encode(singgleEvent)
        }
    }

}

func main() {
    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", homeLink)
    router.HandleFunc("/events", getAllEvents).Methods("GET")
    router.HandleFunc("/event", createEvent).Methods("POST")
    router.HandleFunc("/events/{id}", getEventById).Methods("GET")
    router.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")
    router.HandleFunc("/events/{id}", updateEvent).Methods("PUT")
    log.Fatal(http.ListenAndServe(":8080", router))
}