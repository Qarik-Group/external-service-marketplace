package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/tweedproject/tweed/api"

	"github.com/gorilla/mux"
	"github.com/starkandwayne/external-service-marketplace/tweed"
	"github.com/starkandwayne/external-service-marketplace/util"

	realtweed "github.com/tweedproject/tweed"
)

type Catalog struct {
	realtweed.Catalog
}

func (c *Catalog) Merge(prefix string, other realtweed.Catalog) {
	for _, svc := range other.Services {
		svc.ID = fmt.Sprintf("%s--%s", prefix, svc.ID)
		svc.Name = fmt.Sprintf("[%s] %s", prefix, svc.Name)
		c.Services = append(c.Services, svc)
	}
}

type API struct {
	Config *util.Config
	Bind   string
}

func (a *API) Catalog() (realtweed.Catalog, error) {
	var c realtweed.Catalog

	/*for _, broker := range a.Config.ServiceBrokers {
		cat := tweed.Catalog(broker.URL)
		c.Merge(broker.Prefix, cat)
	}*/
	//loop over catalogs and condense
	cats := tweed.Connect(a.Config).Catalog()
	for _, catalog := range cats {
		c.Services = append(c.Services, catalog.Services[0])
	}

	return c, nil
}

func (a *API) Provision(cmd util.ProvisionCommand, prefix string, service string, plan string) (api.ProvisionResponse, error) {
	fmt.Printf("PROVISIONING [%s][%s][%s]\n", prefix, service, plan)

	broker, found := a.Config.Broker(prefix)
	var nothing api.ProvisionResponse
	if !found {
		return nothing, fmt.Errorf("no such broker '%s'", prefix)
	}

	fmt.Printf("PROVISIONING against tweed at %s (u:%s, p:%s)\n", broker.URL, broker.Username, broker.Password)
	// This is where we dispatch off to the actual broker.
	t := tweed.Connect(a.Config)
	provInst := t.Provision(broker.URL, cmd)

	// in James' ideal world, here's what we do
	/*
		instance, err := broker.Backend.Provision(service, plan)
		if err != nil {
			return "", err
		}

		return instance.ID, nil
	*/
	return provInst, nil
}
func (a *API) Unbind(prefix, instance string, binding string) (api.UnbindResponse, error) {
	fmt.Printf("Unbinding [%s][%s]\n", prefix, instance)

	broker, found := a.Config.Broker(prefix)
	var nothing api.UnbindResponse
	if !found {
		return nothing, fmt.Errorf("no such broker '%s'", prefix)
	}

	fmt.Printf("Unbinding against tweed at %s (u:%s, p:%s)\n", broker.URL, broker.Username, broker.Password)

	t := tweed.Connect(a.Config)
	unbindInst := t.UnBind(broker.URL, instance, binding)

	return unbindInst, nil
}

var config *util.Config
var url string

func findTweed(tweedname string) int {
	brokers := config.ServiceBrokers

	for i := 0; i < len(brokers); i++ {
		if tweedname == brokers[i].Prefix {
			return i
		}
	}
	return -1
}
func testResponse(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Endpoint Hit")
}

func (api API) Run() {
	config = api.Config
	//url = "http://10.128.32.138:31666"
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	//register broker
	r.HandleFunc("/register/{broker}", testResponse)
	//retrieve clouds
	r.HandleFunc("/clouds", testResponse)
	//retrieve a specific cloud
	r.HandleFunc("/clouds/{cloud}", testResponse)

	r.HandleFunc("/catalog", func(w http.ResponseWriter, r *http.Request) {
		cat, err := api.Catalog()
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "internal error: %s\n", err) // FIXME this is bad, don't do it.
			return
		}

		b, err := json.Marshal(cat)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "internal error: %s\n", err) // FIXME this is bad, don't do it.
			return
		}

		w.WriteHeader(200)
		fmt.Fprintf(w, "%s\n", string(b))
	})

	//provision a service
	r.HandleFunc("/provision/{service}/{plan}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		parts := strings.SplitN(vars["service"], "--", 2)
		if len(parts) != 2 {
			panic("ahhhhhh") // FIXME don't do this its bad
		}

		prefix := parts[0]
		service := parts[1]
		plan := vars["plan"]
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error Reading Request Body"))
			return
		}
		var provCmd util.ProvisionCommand
		s := make([]string, 2)
		s[0] = service
		s[1] = plan
		json.Unmarshal(body, &provCmd)

		inst, err := api.Provision(provCmd, prefix, service, plan)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "internal error: %s\n", err) // FIXME this is bad, don"t do it.
			return
		}
		w.WriteHeader(200)
		fmt.Fprintf(w, "OK %s\n", inst) // FIXME - use JSON, give some info back
	})

	//get an instance
	r.HandleFunc("/instances/{instance}", testResponse)
	//deprovision an instance
	r.HandleFunc("/deprovision/{prefix}/{instance}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		inst, err := api.Deprovision(vars["prefix"], vars["instance"])
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "internal error: %s\n", err) // FIXME this is bad, don"t do it.
			return
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(inst)
		//fmt.Fprintf(w, "OK %s\n", inst) // FIXME - use JSON, give some info back
	})
	//bind an instance
	r.HandleFunc("/bind/{prefix}/{instance}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		inst, err := api.BindSVC(vars["prefix"], vars["instance"])
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "internal error: %s\n", err) // FIXME this is bad, don"t do it.
			return
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(inst)
		//fmt.Fprintf(w, "OK %s\n", inst) // FIXME - use JSON, give some info back
	})
	//retrieve binding
	r.HandleFunc("/getbinding/{instance}", testResponse)
	//unbind an instance
	r.HandleFunc("/unbind/{prefix}/{instance}/{binding}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		inst, err := api.Unbind(vars["prefix"], vars["instance"], vars["binding"])
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "internal error: %s\n", err) // FIXME this is bad, don"t do it.
			return
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(inst)
		//fmt.Fprintf(w, "OK %s\n", inst) // FIXME - use JSON, give some info back
	})

	//start the server
	log.Fatal(http.ListenAndServe(api.Bind, r))
}
