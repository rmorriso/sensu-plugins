package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/golang/glog"
	"github.com/rmorriso/pager"
)

const (
	OK       = 0
	WARN     = 1
	CRITICAL = 2
)

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "f", "./pagerduty.yaml", "the PagerDuty config file")
}

// need to save incident id so it can be used to resolve in pager duty (Check.Action == "resolve")
// and to avoid creating duplicate incidents after the first (check pagerduty.rb for details)
// use http://sensuapp.org/docs/0.16/api_stashes to save state?
func main() {
	flag.Parse()

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		glog.Fatalf("CAS config file: %s\n", err)
	}

	args := flag.Args()
	glog.V(5).Infof("Args: %v\n", args)

	config, err := Init(configFile)
	if err != nil {
		glog.Fatalf("Error loading configuration %s\n", err)
	}

	// Global usage
	pager.ServiceKey = config.ServiceKey

	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Printf("Sensu system failure: %s", err)
		return
	}

	e := &event{}
	err = json.Unmarshal(input, &e)
	// TODO: should we create a PagerDuty incident for Sensu failures?
	if err != nil {
		fmt.Printf("Sensu system failure: %s", err)
		return
	}

	// have we already seen this event?
	// GET http://api:4567/stashes/events/id
	// POST http://api:4567/stashes/events/id

	id := e.ID

	switch e.Check.Status {
	case OK:
		if e.Action == "resolve" {
			// send Pagerduty a resolve request
			resolveIncident(id)
		}
		return
	case WARN:
		// if the history is all warnings or criticals then create an incident
		for _, x := range e.Check.History {
			if x != "1" || x != "2" {
				return
			}
		}
		// create a PagerDuty incident for CRITICAL alerts
		alert := fmt.Sprintf("%s/%s: recurring WARNINGS for: %s", e.Client.Name, e.Check.Name, e.Check.Output)
		createIncident(alert)
		return
	case CRITICAL:
		// create a PagerDuty incident for CRITICAL alerts
		alert := fmt.Sprintf("%s/%s: %s", e.Client.Name, e.Check.Name, e.Check.Output)
		createIncident(alert)
		return
	}

}

func createIncident(alert string) {
	if incidentKey, err := pager.Trigger(alert); err != nil {
		fmt.Printf("PagerDuty error %s\n", err)
	} else {
		fmt.Printf("Incident Key %s\n", incidentKey)
	}
}

func resolveIncident(id string) {
}
